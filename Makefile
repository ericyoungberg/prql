#-- Define the world
BUILDDIR = build
PKG := github.com/prql/prql
ARCH ?= darwin/amd64

PRQL_BIN  = prql
PRQLD_BIN = prqld
PRQL_DIR  = cli
PRQLD_DIR = prqld

BUILD_CONTAINER = prql-builder



#-- Generate flags
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif

CTIMEVAR = -X $(PKG)/version.GITCOMMIT=$(GITCOMMIT)
GO_LDFLAGS = -ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC = -ldflags "-w $(CTIMEVAR) -extldflags -static"

GOOSARCHES = darwin/amd64 darwin/386 freebsd/amd64 freebsd/386 linux/arm linux/arm64 linux/amd64 linux/386 solaris/amd64 windows/amd64 windows/386


#-- Build the world
all: clean build-prql build-prqld staticcheck install


.PHONY: prql
prql: test-prql build-prql install

.PHONY: prqld
prqld: test-prqld build-prqld install


build-prql: $(PRQL_DIR)/*.go
		@echo "+ $@"
		@go build ${GO_LDFLAGS} -o $(BUILDDIR)/$(PRQL_BIN) ./$(PRQL_DIR)


build-prqld: $(PRQLD_DIR)/*.go
		@echo "+ $@"
		@go build ${GO_LDFLAGS} -o $(BUILDDIR)/$(PRQLD_BIN) ./$(PRQLD_DIR)


.PHONY: with-docker
with-docker:
		@echo "+ $@"
		docker build -t $(BUILD_CONTAINER) .
		docker run --rm -t -v $(PWD):/go/src/$(PKG) -e "ARCH=$(ARCH)" $(BUILD_CONTAINER)
		docker rmi $(BUILD_CONTAINER)
		@chown -R $(whomai):$(whoami) $(BUILDDIR)


define buildrelease
GOOS=$(3) GOARCH=$(4) CGO_ENABLED=0 go build \
	 -o $(BUILDDIR)/$(1)-$(3)-$(4) \
	 -a -tags "static_build netgo" \
	 -installsuffix netgo ${GO_LDFLAGS_STATIC} ./$(2);
md5sum $(BUILDDIR)/$(1)-$(3)-$(4) > $(BUILDDIR)/$(1)-$(3)-$(4).md5;
sha256sum $(BUILDDIR)/$(1)-$(3)-$(4) > $(BUILDDIR)/$(1)-$(3)-$(4).sha256;
endef

.PHONY: static
static:
		@echo "+ $@"
		$(call buildrelease,$(PRQL_BIN),$(PRQL_DIR),$(subst /,,$(dir $(ARCH))),$(notdir $(ARCH)))
		$(call buildrelease,$(PRQLD_BIN),$(PRQLD_DIR),$(subst /,,$(dir $(ARCH))),$(notdir $(ARCH)))


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
		@cp $(BUILDDIR)/* $(GOPATH)/bin


.PHONY: clean
clean:
		@echo "+ $@"
		rm -rf $(BUILDDIR)
		rm -rf $(GOPATH)/bin/$(PRQL_BIN) $(GOPATH)/bin/$(PRQLD_BIN)
