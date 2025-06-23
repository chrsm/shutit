.DEFAULT_GOAL := all

GOOS=windows
export GOOS

COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date)

.PHONY: all
all:
	go build -trimpath -ldflags="-X 'bits.chrsm.org/shutit.BuildDate=$(DATE)' -X bits.chrsm.org/shutit.BuildRev=$(COMMIT)" -v ./cmd/shutit

