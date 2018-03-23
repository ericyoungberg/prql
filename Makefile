#-- Define the world
PROJDIR := $(shell pwd)

PRQL_BIN  = prql
PRQLD_BIN = prqld


#-- Build the world
all: prql prqld staticcheck install


.PHONY: prql
prql: test-prql build-prql

.PHONY: prqld
prqld: test-prqld build-prqld


build-prql: cli/*.go
				@echo "+ $@"
				@go build -o $(PRQL_BIN) ./cli


build-prqld: prqld/*.go
				@echo "+ $@"
				@go build -o $(PRQLD_BIN) ./prqld


.PHONY: test-prql
test-prql:
				@echo "+ $@"
				@echo "No tests exist for prql"


.PHONY: test-prqld
test-prqld:
				@echo "+ $@"
				@echo "No tests exist for prqld"


.PHONY: staticcheck
staticcheck:
				@echo "+ $@"
				@staticcheck $(shell go list ./... | grep -v vendor) | grep -v '.pb.go:' | tee /dev/stderr


.PHONY: install
install:
				@echo "+ $@"
				@go install -a


.PHONY: clean
clean:
				@echo "+ $@"
				@go clean
				$(RM) $(PRQL_BIN) $(PRQLD_BIN)
