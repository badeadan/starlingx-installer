.PHONY: all clean

COMMANDS=$(wildcard cmd/*)
BINARIES=$(foreach cmd,${COMMANDS},bin/$(notdir ${cmd}))

bin/%: cmd/%/main.go
	go build -o $@ ./cmd/$*

all: $(BINARIES)

clean:
	rm -f $(BINARIES)
