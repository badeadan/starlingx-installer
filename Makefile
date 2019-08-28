.PHONY: all clean

COMMANDS=$(wildcard cmd/*)
BINARIES=$(foreach cmd,${COMMANDS},bin/$(notdir ${cmd}))
LDFLAGS=-ldflags='-s -w'
UPX := $(shell command -v upx 2> /dev/null)

bin/%: cmd/%/main.go
	cd ./cmd/$* ; packr2
	go build ${LDFLAGS} -o $@ ./cmd/$*
ifdef UPX
	upx --best $@
endif

all: $(BINARIES)

clean:
	rm -rf cmd/*/packrd
	rm -f cmd/*/*-packr.go
	rm -f $(BINARIES)
