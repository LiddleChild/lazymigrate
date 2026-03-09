.PHONY: debug
debug:
	tail -f ~/.local/share/lazymigrate/debug.log

.PHONY: test
test:
	go test -v ./...
