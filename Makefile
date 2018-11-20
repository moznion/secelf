.PHONY: check test lint vet fmt-check fmt
PKGS := $(shell go list ./...)

check: test lint vet fmt-check

build:
	go build cmd/secelf/secelf.go

test:
	go test -v -cover ./...

lint:
	golint $(PKGS) 2>&1 \
		| grep -v 'should have comment or be unexported' \
		| grep . ; \
		EXIT_CODE=$$? ; \
		if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi

vet:
	go vet ./...

fmt-check:
	gofmt -l -s **/*.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi; \
	goimports -l **/*.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi \

fmt:
	gofmt -w -s **/*.go
	goimports -w **/*.go

