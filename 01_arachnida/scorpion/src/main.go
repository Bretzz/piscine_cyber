package main

// https://pkg.go.dev/github.com/evanoberholster/imagemeta#Decode

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/djherbis/times"
	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta/exif"
)

func print(tag string, value any) {
	fmt.Printf("%40s: %v\n", tag, value)
}

func print_file_header(path string) {
	fmt.Println("========================================================================================")
	fmt.Printf("> %s\n", filepath.Base(path))
	fmt.Println("========================================================================================")
}

// WalkExif mimics the behavior of a walk function by using reflection
// to iterate over every field in the imagemeta.Exif struct.
func WalkExif(ex *exif.Exif, walkFn func(name string, value interface{})) {
	val := reflect.ValueOf(ex).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if fieldType.PkgPath != "" {
			continue
		}

		// If the field is a nested struct (like GPS or a sub-block),
		// you can choose to recurse, or just print its value.
		if fieldVal.Kind() == reflect.Struct {
			// Optional: handle nested structs if needed
		}

		// Skip zero/empty values to keep output clean (optional)
		if fieldVal.IsZero() {
			continue
		}

		walkFn(fieldType.Name, fieldVal.Interface())
	}
}

func print_full_exif_metadata(path string) {

	f := func(path string) error {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Failed to open file: %v", err)
		}
		defer file.Close()

		ex, err := imagemeta.Decode(file)
		if err != nil {
			return fmt.Errorf("Failed to decode metadata: %v", err)
		}

		// Convert to JSON bytes
		data, err := json.MarshalIndent(ex, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(data))

		return nil
	}

	err := f(path)

	if err != nil {
		fmt.Printf("[Warn] %s\n", err.Error())
	}

	fmt.Println("----------------------------------------------------------------------------------------")
}

func get_birth_time(path string) (string, error) {
	t, err := times.Stat(path)
	if err != nil {
		return "<3", err
	}

	if t.HasBirthTime() {
		return t.BirthTime().Local().String(), nil
	} else {
		return "Not supported", nil
	}
}

func print_standard_metadata(path string) {

	f := func(path string) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return err
		}

		// Get directory
		dir, err := filepath.Abs(path)
		if err != nil {
			dir = "Error"
		} else {
			dir = filepath.Dir(dir)
		}

		// Get btime
		btime, err := get_birth_time(path)
		if err != nil {
			btime = "Error"
		}

		print("Filemame", stat.Name())
		print("Directory", dir)
		print("Size", fmt.Sprintf("%v bytes", stat.Size()))
		print("Permissions", stat.Mode())
		print("Creation Date", btime)
		print("Modification Date", stat.ModTime())

		return nil
	}

	err := f(path)

	if err != nil {
		fmt.Printf("[Warn] %s\n", err.Error())
	}

	fmt.Println("----------------------------------------------------------------------------------------")
}

func is_dir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// do directory stuff
		return true, nil
	}
	return false, nil
}

func ask_for_enter() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("[Enter]")

	_, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Missing Argument(s)")
	}

	for i, file := range os.Args[1:] {
		print_file_header(file)

		// check if dir
		ret, err := is_dir(file)
		if err != nil {
			fmt.Printf("[Error] opening file: %s\n", err)
			continue
		}
		if ret == true {
			fmt.Printf("[Error] '%s' is a directory.\n", file)
			continue
		}
		print_standard_metadata(file)
		// print_exif_medatada(file)
		print_full_exif_metadata(file)

		if i != len(os.Args)-2 {
			ask_for_enter()
		}
	}
}
