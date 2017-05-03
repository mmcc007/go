.DEFAULT_GOAL := test
setup:
	# get the grpc stuff
	go get google.golang.org/grpc

	# get the protobuffers
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go

	# generate the stubs for helloworld
	$(GOPATH)/bin/protoc -I examples/helloworld/helloworld examples/helloworld/helloworld/helloworld.proto --go_out=plugins=grpc:examples/helloworld/helloworld		
	# generate the stubs for route_guide
	$(GOPATH)/bin/protoc -I examples/route_guide/routeguide examples/route_guide/routeguide/route_guide.proto --go_out=plugins=grpc:examples/route_guide/routeguide
	# build the server for testing 
	go build -o examples/route_guide/server/server examples/route_guide/server/server.go

	# download the mock tools
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen

	# generate the mock server for helloworld client test
	$(GOPATH)/bin/mockgen -destination examples/helloworld/mock_helloworld/hw_mock.go github.com/mmcc007/go/examples/helloworld/helloworld GreeterClient

test: setup
	# run all tests
	go test -v ./... > test.output

	# convert the test output to junit format for reporting on Jenkins
	go get github.com/tebeka/go2xunit
	cat test.output | $(GOPATH)/bin/go2xunit -output tests.xml
