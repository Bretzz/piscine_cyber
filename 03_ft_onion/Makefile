UNAME := $(shell uname)
NAME := ft_onion-app

BFLAGS = -t $(NAME)
ifeq ($(UNAME), Darwin)
	BFLAGS +=
else
	BFLAGS +=
endif

all: build run

build:
	docker build $(BFLAGS) .

run:
	#docker run -d -p 4242:4242 -p 8080:80 -p 9050:9050 --name $(NAME)-running $(NAME)
	docker run -d -p 8080:80 -p 9050:9050 --name $(NAME)-running $(NAME)

itrain:
	docker exec -it $(NAME)-running sh

logs:
	docker logs $(NAME)-running

clean:
	docker rm -f $(NAME)-running

fclean: clean
	docker image prune -af

.PHONY: build run itrain logs clean fclean