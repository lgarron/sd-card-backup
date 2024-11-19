PROJECT = github.com/lgarron/sd-card-backup/...

.PHONY: all
all: test

.PHONY: lint
lint:
	go vet ${PROJECT}
	golint -set_exit_status ${PROJECT}

.PHONY: test
test: lint
	go test -cover ${PROJECT}

.PHONY: test-verbose
test-verbose: lint
	go test -v -cover ${PROJECT}

.PHONY: run
run:
	go run cmd/sd-card-backup/*.go

.PHONY: install
install:
	go install ${PROJECT}
