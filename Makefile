.PHONY: all clean

COMMANDS=$(wildcard cmd/*)
BINARIES=$(foreach cmd,${COMMANDS},bin/$(notdir ${cmd}))

bin/%: cmd/%/main.go
	cd ./cmd/$* ; packr2
	go build -o $@ ./cmd/$*

all: $(BINARIES)

clean:
	rm -rf cmd/*/packrd
	rm -f cmd/*/*-packr.go
	rm -f $(BINARIES)
