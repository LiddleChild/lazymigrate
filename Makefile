.PHONY: debug
debug:
	tail -f debug.log

.PHONY: test
test:
	go test -v ./...
