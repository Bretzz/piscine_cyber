NAME := ft_otp
VALIDATOR := $(NAME)_validator
ROOT := $(shell pwd)
SRCDIR := $(ROOT)/src

all: $(NAME)

$(NAME):
	@cd $(SRCDIR) && go build ./cmd/basic
	@mv $(SRCDIR)/basic $(ROOT)/$@

validator: $(VALIDATOR)

$(VALIDATOR):
	@cd $(SRCDIR) && go build ./cmd/server
	@mv $(SRCDIR)/server $(ROOT)/$(VALIDATOR)

key.txt:
	@openssl rand -hex 64 | tr -d '\n' > $@

.env:
	@echo "AES_KEY=\"$$(openssl rand -base64 32)\"" > $@

fclean:
	@rm -f $(NAME) .env key.txt ft_otp.key
	@rm -rf otpkeys

re: fclean all

.PHONY: fclean re