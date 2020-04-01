CHS_ENV_HOME  ?= $(HOME)/.chs_env
TESTS         ?= ./...

bin           := company-search-consumer
chs_envs      := $(CHS_ENV_HOME)/global_env $(CHS_ENV_HOME)/company-search-consumer/env
source_env    := for chs_env in $(chs_envs); do test -f $$chs_env && . $$chs_env; done
xunit_output  := test.xml
lint_output   := lint.txt

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

.PHONY: all
all: clean build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build: fmt
	go build

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	go test $(TESTS) -run 'Unit' -coverprofile=coverage.out


.PHONY: test-integration
test-integration:
	$(source_env); go test $(TESTS) -run 'Integration'

.PHONY: clean
clean:
	go mod tidy
	rm -f $(bin)
	rm -f *.zip
	rm -rf build-*

.PHONY: package
package:
	$(eval tmpdir:=$(shell mktemp -d build-XXXXXXXXXX))
	cp ./$(bin) $(tmpdir)/$(bin)
	cp ./start.sh $(tmpdir)/start.sh
	cd $(tmpdir) && zip -r ../$(bin)-$(version).zip $(bin) start.sh
	rm -rf $(tmpdir)

.PHONY: dist
dist: clean build package

.PHONY: xunit-tests
xunit-tests: GO111MODULE=off
xunit-tests:
	go get github.com/tebeka/go2xunit
	go test -v $(TESTS) -run 'Unit' | go2xunit -output $(xunit_output)

.PHONY: lint
lint: GO111MODULE=off
lint:
	go get github.com/golang/lint/golint
	golint ./... > $(lint_output)
