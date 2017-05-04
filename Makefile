.DEFAULT_GOAL := test

CHDIR_SHELL := $(SHELL)
define chdir
   $(eval _D=$(firstword $(1) $(@D)))
   $(info $(MAKE): cd $(_D)) $(eval SHELL = cd $(_D); $(CHDIR_SHELL))
endef

# check to see if protobuf compiler needs to be built
PROTOC_VERSION = "$(shell protoc --version)"
PROTOC_VERSION_REQ = "libprotoc 3.3.0"
ifneq ($(PROTOC_VERSION),$(PROTOC_VERSION_REQ))
	BUILD_PROTOC = TRUE
endif

# check to see if go needs to be installed
GO_VERSION = "$(shell go version)"
GO_VERSION_REQ = "go version go1.8.1 linux/amd64"
ifneq ($(GO_VERSION),$(GO_VERSION_REQ))
	INSTALL_GO = TRUE
endif

protobuf/install: protocdownload
ifdef BUILD_PROTOC
	$(call chdir)
	./autogen.sh
	./configure
	make
	#make check
	sudo make install
	sudo ldconfig
endif

protocdownload:
ifdef BUILD_PROTOC
	# no protoc at version 3.3.0 download and build and install
	sudo apt-get install autoconf automake libtool curl make g++ unzip
	git clone https://github.com/google/protobuf.git
endif

goinstall:
ifdef INSTALL_GO
	# no golang at version 1.8.1 download and install
	wget https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz
	sudo tar -C /usr/local -xzf go1.8.1.linux-amd64.tar.gz
	export PATH=/usr/local/go/bin:$PATH
endif

setup: goinstall protobuf/install
	rm -rf protobuf

build:
	# get the grpc stuff
	go get google.golang.org/grpc

	# get the protobuffers
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go

	# generate the stubs for helloworld
	protoc -I examples/helloworld/helloworld examples/helloworld/helloworld/helloworld.proto --go_out=plugins=grpc:examples/helloworld/helloworld		
	# generate the stubs for route_guide
	protoc -I examples/route_guide/routeguide examples/route_guide/routeguide/route_guide.proto --go_out=plugins=grpc:examples/route_guide/routeguide
	# build the server for testing 
	go build -o examples/route_guide/server/server examples/route_guide/server/server.go

	# download the mock tools
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen

	# generate the mock server for helloworld client test
	$(GOPATH)/bin/mockgen -destination examples/helloworld/mock_helloworld/hw_mock.go github.com/mmcc007/go/examples/helloworld/helloworld GreeterClient

test: 
	# run all tests
	go test -v ./... > test.output

	# convert the test output to junit format for reporting on Jenkins
	go get github.com/tebeka/go2xunit
	cat test.output | $(GOPATH)/bin/go2xunit -output tests.xml
