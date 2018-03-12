GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

test:
	@go test -ldflags -s -v $(GOPACKAGES)


.PHONY: test