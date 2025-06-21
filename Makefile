.DEFAULT_GOAL := all

GOOS=windows
export GOOS

.PHONY: all
all:
	go build -v

