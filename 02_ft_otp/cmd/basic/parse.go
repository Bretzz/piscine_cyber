package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
)

type Mode struct {
	NewKey bool // Generate a new .key FILE
	GetOTP bool // Get new OTP value from previousy stored key
	// SecOTP  bool
	Autoenv bool   // Load the specified file as environment
	EnvFILE string // the env file
	SecFILE string // the file with the plaintext secret
	KeyFILE string // the .key file
	// SecARG  string
}

func visitedFlags() []string {
	visited := []string{}
	flag.Visit(func(f *flag.Flag) {
		visited = append(visited, f.Name)
	})
	return visited
}

func parse() (Mode, error) {
	// fmt.Printf("Hello World\n")

	// f1 := flag.NewFlagSet("f1", flag.ContinueOnError)
	// help := f1.Bool("h", false, "Print an help menu")
	key := flag.String("k", "ft_otp.key", "Produce an OTP value from the encrypted key stored in FILE.")
	gen := flag.String("g", "key.txt", "Generate a new key based on a plaintext file.")
	// sec := flag.String("s", "64-bytes-otp-secret", "Generate a new key based on SECRET argument")
	env := flag.String("e", ".env", "Source the .env file for enviroment variables.")
	// out := flag.String("o", "ft_otp.key", "Specifies the output file for the key generation.")

	// flag.
	flag.Parse()

	visited := visitedFlags()
	// fmt.Printf("visited %v\n", visited)

	// Errors
	// No flags
	if len(os.Args) == 1 {
		return Mode{}, fmt.Errorf("No flag passed.")
		// os.Exit(1)

	}
	// -k, -g (just one of them)
	if slices.Contains(visited, "k") && slices.Contains(visited, "g") {
		return Mode{}, fmt.Errorf("Incompatible flags: -k, -g.")
	}

	// -s
	/* if slices.Contains(visited, "s") == true {
		// : Invalid arg number
		if len(os.Args) != 3+envAdd {
			return Mode{}, fmt.Errorf("-s flag needs precisely 1 argument.\n")
			// os.Exit(1)
		}
		// : Invalid filename
		if *sec == "" {
			return Mode{}, fmt.Errorf("Invalid -s ARG.\n")
		}
	} */

	mode := Mode{}

	// -h
	// if *help == true {
	// 	fmt.Println("ft_otp client")
	// 	fmt.Println("-------------")
	// 	fmt.Println("Usage: ft_otp -[kgeh] [FILE]")
	// 	fmt.Println("flags:")
	// 	fmt.Println("\t-k [FILE]: Generate an OTP value from the encrypted key stored in FILE (default ft_otp.key).")
	// 	fmt.Println("\t-g FILE: Generate a new key based on a plaintext file.")
	// 	fmt.Println("\t-s SECRET: Generate a new key based on SECRET argument.")
	// 	fmt.Println("\t-e FILE: Source the .env file for enviroment variables.")
	// 	fmt.Println("\t-h: Print this menu.")

	// 	mode.Help = true
	// 	return mode, nil
	// }

	// -e
	if slices.Contains(visited, "e") == true {
		mode.Autoenv = true
		mode.EnvFILE = *env
	}

	// -k
	if slices.Contains(visited, "k") == true {
		mode.GetOTP = true
		mode.KeyFILE = *key
		return mode, nil
	}

	// -s
	/* if slices.Contains(visited, "s") == true {
		mode.SecOTP = true
		mode.SecARG = *sec
		return mode, nil
	} */

	// -g
	mode.NewKey = true
	mode.KeyFILE = *gen
	return mode, nil
}
