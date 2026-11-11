package totp

import (
	"crypto/subtle"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"wasted-kiwi.com/ft_otp/internal/aes"
	"wasted-kiwi.com/ft_otp/internal/algo"
)

// LoadKey reads and decrypts a key file, returning the RAW key bytes.
// The stored plaintext is the 64-char hex string, so decode it here.
// (If you changed GenerateKey to store raw bytes already, drop the decode.)
func LoadKey(path string) ([]byte, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	plaintext, err := aes.Decrypt(buf, []byte(AES_KEY))
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Verifier holds the state for ONE account.
type Verifier struct {
	mu        sync.Mutex
	key       []byte
	lastUsed  uint64 // highest step already consumed
	statePath string // where lastUsed is persisted
}

func (v *Verifier) Verify(code uint32) (bool, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	cur := uint64(time.Now().UTC().Unix()) / TIME_STEP

	for _, step := range []uint64{cur - 1, cur, cur + 1} {
		if step <= v.lastUsed {
			continue // replayed, or older than a used step
		}
		expected := algo.HOTP(v.key, step)
		if subtle.ConstantTimeEq(int32(expected), int32(code)) != 1 {
			continue
		}

		// Persist first; fail closed if we can't record the use.
		prev := v.lastUsed
		v.lastUsed = step
		if err := v.save(); err != nil {
			v.lastUsed = prev
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// save writes lastUsed atomically (temp file + rename).
func (v *Verifier) save() error {
	tmp := v.statePath + ".tmp"
	data := []byte(strconv.FormatUint(v.lastUsed, 10))
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, v.statePath)
}

func loadVerifier(keyPath string) (*Verifier, error) {
	key, err := LoadKey(keyPath)
	if err != nil {
		return nil, err
	}

	v := &Verifier{
		key:       key,
		statePath: strings.TrimSuffix(keyPath, ".key") + ".step",
	}

	data, err := os.ReadFile(v.statePath)
	switch {
	case err == nil:
		last, perr := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
		if perr != nil {
			return nil, perr
		}
		v.lastUsed = last
	case os.IsNotExist(err):
		// first use: lastUsed stays 0
	default:
		return nil, err
	}
	return v, nil
}

// Registry maps account -> Verifier. The zero value is ready to use.
type Registry struct {
	mu        sync.RWMutex
	verifiers map[string]*Verifier
}

func (r *Registry) get(account, keyPath string) (*Verifier, error) {
	r.mu.RLock()
	v := r.verifiers[account] // reading a nil map is fine
	r.mu.RUnlock()
	if v != nil {
		return v, nil
	}

	// Load outside the lock so disk/crypto work doesn't block other accounts.
	loaded, err := loadVerifier(keyPath)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.verifiers == nil {
		r.verifiers = make(map[string]*Verifier)
	}
	// Another request may have loaded it meanwhile; keep the first one.
	if existing := r.verifiers[account]; existing != nil {
		return existing, nil
	}
	r.verifiers[account] = loaded
	return loaded, nil
}

func (r *Registry) Verify(account, keyPath string, code uint32) (bool, error) {
	v, err := r.get(account, keyPath)
	if err != nil {
		return false, err
	}
	return v.Verify(code)
}

// Reset drops the cached verifier and its saved step. Call it whenever
// the account's key file is (re)generated.
func (r *Registry) Reset(account, keyPath string) {
	r.mu.Lock()
	delete(r.verifiers, account)
	r.mu.Unlock()
	_ = os.Remove(strings.TrimSuffix(keyPath, ".key") + ".step")
}
