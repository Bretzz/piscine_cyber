package main

import (
	"fmt"
	"os"

	"wasted-kiwi.com/ft_otp/internal/autoenv"
	"wasted-kiwi.com/ft_otp/internal/totp"
)

func main() {

	// Parse input
	mode, err := parse()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	// fmt.Println(mode)

	// -e
	if mode.Autoenv == true {
		// Load .env
		// fmt.Print("opening file ", mode.EnvFILE, "\n")
		err := autoenv.Load(mode.EnvFILE)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
	}

	// load env variables
	{
		err := totp.LoadEnv()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
	}

	// -g FILE
	if mode.NewKey == true {
		// fmt.Printf("reading all file '%s'\n", mode.KeyFILE)
		secret, err := os.ReadFile(mode.KeyFILE)
		if err != nil {
			fmt.Printf("error2: %s\n", err)
			os.Exit(1)
		}
		err = totp.GenerateKey(secret, "ft_otp.key")
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Key was successfully saved in ft_otp.key.\n")
		return
	}

	// -s SECRET
	/* if mode.SecOTP == true {
		err := GenerateKey([]byte(mode.SecARG))
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		return
	} */

	// -k
	if mode.GetOTP == true {
		// generate value
		values, err := totp.ProduceOTP(mode.KeyFILE)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			os.Exit(1)
		}
		// print otp [-1, 0, 1]
		fmt.Println(totp.StringOTP(values[1]))
		return
	}
}
