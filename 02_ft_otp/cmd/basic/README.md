# About ft_otp
Simple Hash based One Time Password implementation in Golang.

## Usage
```
ft_otp -[kgeh] [FILE]
flags:
	-k: Generate an OTP value from the encrypted key stored in ft_otp.key.
	-g FILE: Generate a new key based on a plaintext file.
	-e FILE: Source the .env file for enviroment variables.
	-h: Print this menu.
```
```bash
#.env
TIME_STEP=30                                # Default
AES_KEY="AES256Key-32Characters1234567890"  # Mandatory
```

## Resources
- [RFC. 4226](https://www.rfc-editor.org/info/rfc4226/)

## TODO
- [x] Read file main.go
- [ ] Clean `gen` back up
- [ ] Make `gen` a library and import it
- [ ] Move OTPKey class in separate file
- [ ] Auth server (timestep, tollerance, one time validation)
