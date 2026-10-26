# About ft_otp
Simple Hash based One Time Password implementation in Golang.

## Usage
```
ft_otp -[egkh] [FILE]
flags:
	-e string
        Source the .env file for enviroment variables. (default ".env")
	-g string
        Generate a new key based on a plaintext file. (default "key.txt")
	-k string
        Produce an OTP value from the encrypted key stored in FILE. (default "ft_otp.key")
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
- [x] Clean `gen` back up
- [x] Make `gen` a library and import it
- [x] Move OTPKey class in separate file
- [x] Auth server (timestep, tollerance, one time validation)
