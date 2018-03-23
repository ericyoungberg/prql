#-- Define the world
PROJDIR := $(shell pwd)
BUILDDIR := build

PRQL_BIN  = prql
PRQLD_BIN = prqld
PRQL_DIR  = cli
PRQLD_DIR = prqld

RM = rm -rf


#-- Build the world
all: clean prqld prql staticcheck install


.PHONY: prql
prql: test-prql build-prql

.PHONY: prqld
prqld: test-prqld build-prqld


build-prql: $(PRQL_DIR)/*.go
		@echo "+ $@"
		@go build -o $(BUILDDIR)/$(PRQL_BIN) -v ./$(PRQL_DIR)


build-prqld: $(PRQLD_DIR)/*.go
		@echo "+ $@"
		@go build -o $(BUILDDIR)/$(PRQLD_BIN) -v ./$(PRQLD_DIR)


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
		cp $(BUILDDIR)/* $(GOPATH)/bin


.PHONY: clean
clean:
		@echo "+ $@"
		$(RM) $(BUILDDIR)
		$(RM) $(GOPATH)/bin/$(PRQL_BIN) $(GOPATH)/bin/$(PRQLD_BIN)
