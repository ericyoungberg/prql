#-- Define the world
BUILD_DIR = build
PKG 		  = github.com/prql/prql
ARCH 		 ?= darwin/amd64

PRQL_BIN  = prql
PRQLD_BIN = prqld
PRQL_DIR  = cli
PRQLD_DIR = prqld

BUILD_CONTAINER = prql-builder

# Set default compiler
GO := go

#-- Generate flags
go_ld_flags = -ldflags "-w -X ${PKG}/$(1)/version.VERSION=$(shell cat $1/VERSION)"
go_ld_flags_static = -ldflags "-w -X ${PKG}/$(1)/version.VERSION=$(shell cat $1/VERSION) -extldflags -static"

GOOSARCHES = darwin/amd64 darwin/386 freebsd/amd64 freebsd/386 linux/arm linux/arm64 linux/amd64 linux/386 solaris/amd64 windows/amd64 windows/386


#-- Build the world
all: clean build-prql build-prqld staticcheck install

.PHONY: prql
prql: test-prql build-prql install

.PHONY: prqld
prqld: test-prqld build-prqld install

define build
@${GO} build $(call go_ld_flags,${2}) -o ${BUILD_DIR}/${1} ./${2}
endef

build-prql: $(PRQL_DIR)/*.go
		@echo "+ $@"
		$(call build,${PRQL_BIN},${PRQL_DIR})

build-prqld: $(PRQLD_DIR)/*.go
		@echo "+ $@"
		$(call build,${PRQLD_BIN},${PRQLD_DIR})

.PHONY: with-docker
with-docker:
		@echo "+ $@"
		docker build -t $(BUILD_CONTAINER) .
		docker run --rm -t -v $(PWD):/go/src/$(PKG) -e "ARCH=$(ARCH)" $(BUILD_CONTAINER)
		docker rmi $(BUILD_CONTAINER)
		@chown -R $(whomai):$(whoami) $(BUILD_DIR)

define buildrelease
GOOS=$(3) GOARCH=$(4) CGO_ENABLED=0 go build \
	 -o $(BUILD_DIR)/$(1)-$(3)-$(4) \
	 -a -tags "static_build netgo" \
	 -installsuffix netgo ${GO_LDFLAGS_STATIC} ./$(2);
md5sum $(BUILD_DIR)/$(1)-$(3)-$(4) > $(BUILD_DIR)/$(1)-$(3)-$(4).md5;
sha256sum $(BUILD_DIR)/$(1)-$(3)-$(4) > $(BUILD_DIR)/$(1)-$(3)-$(4).sha256;
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
		@cp $(BUILD_DIR)/* $(GOPATH)/bin

.PHONY: clean
clean:
		@echo "+ $@"
		rm -rf $(BUILD_DIR)
		rm -rf $(GOPATH)/bin/$(PRQL_BIN) $(GOPATH)/bin/$(PRQLD_BIN)


.PHONY: bump-version
BUMP := patch
bump-version: ## Bump the version in the version file. Set BUMP to [ patch | major | minor ].
	@$(GO) get -u github.com/jessfraz/junk/sembump # update sembump tool
	$(eval NEW_VERSION = $(shell sembump --kind $(BUMP) $(VERSION)))
	@echo "Bumping VERSION from $(VERSION) to $(NEW_VERSION)"
	echo $(NEW_VERSION) > VERSION
	git add VERSION.txt
	git commit -vsam "Bump version to $(NEW_VERSION)"
	@echo "Run make tag to create and push the tag for new version $(NEW_VERSION)"
