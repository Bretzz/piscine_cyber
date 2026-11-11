# About ft_otp
This project let's you build your **OTP generator** and **validator**. The algorithm used is the **TOTP** (RFC 6238), which is a time-based version of the **HOTP** algorithm.

## How to build ft_otp
To build **ft_otp** a makefile is provided.
- `make` to build **ft_otp**
- `make validator` to build **ft_otp_validator**

## How to use ft_otp
The source code provided builds 2 different executables:
- `ft_otp`: core algorithm implementation. Its behaviour is described in the **42 - ft_otp** subject (`./ft_otp -h` for further explaination)
- `ft_otp_validator`: simple TOTP validator with very basic GUI hosted at `https://localhost:8080` (Note that the 'public/index.html' folder must live besides the executable)

## .env example
```bash
AES_KEY="AES-GCM-32-bytes-encription-key-with-base64-encoding"
```
to generate 32 bytes with base64 encoding: `openssl rand -base64 32`

## key.txt example
To generate a valid key.txt file
`openssl rand -hex N | tr -d '\n' > key.txt`
- **N** is the length of the data (min 64)
- **openssl** is the program that generates keys
- **tr -d '\n'** removes trailing newlines