UNAME := $(shell uname)
NAME := ft_onion-app
KEY_PATH ?= $$HOME/.ssh/id_ed25519.pub

BFLAGS = --secret id=ssh_key,src=$(KEY_PATH) -t $(NAME)
ifeq ($(UNAME), Darwin)
	BFLAGS +=
else
	BFLAGS +=
endif

all: build run

build:
	docker build $(BFLAGS) .

run:
	docker run -d -p 9050:9050 --name $(NAME)-running $(NAME)

itrain:
	@docker exec -it $(NAME)-running sh

logs:
	@docker logs $(NAME)-running

address:
	@docker logs ft_onion-app-running 2>&1 | grep -oE '[a-z2-7]{56}\.onion'

clean:
	docker rm -f $(NAME)-running

fclean: clean
	docker image prune -af

.PHONY: build run itrain logs clean fclean address