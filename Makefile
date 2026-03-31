.PHONY: check
check: lint test

PROJECT = github.com/lgarron/sd-card-backup/...

.PHONY: lint
lint:
	go vet ${PROJECT}

.PHONY: test
test: test-go test-build-lib test-build-bin

.PHONY: test-go
test-go:
	go test -cover ${PROJECT}

.PHONY: test-build-lib
test-build-lib:
	go build .

.PHONY: test-build-bin
test-build-bin:
	go build -o /dev/null ./cmd/sd-card-backup

.PHONY: test-verbose
test-verbose: lint
	go test -v -cover ${PROJECT}

.PHONY: run
run:
	go run cmd/sd-card-backup/*.go

.PHONY: install
install:
	go install ${PROJECT}
