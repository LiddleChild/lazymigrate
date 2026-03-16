.PHONY: debug
debug:
	tail -f ~/.local/share/lazymigrate/debug.log

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: install
install:
	go install
