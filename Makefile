.PHONY: dist/ambihass

GC = go build
GO_FLAGS = GOGC=off
ENTRYPOINT = ./ambihass.go

all: dist

dist/ambihass:
	$(GO_FLAGS) $(GC) -o dist/ambihass $(ENTRYPOINT)
