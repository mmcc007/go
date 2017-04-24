setup:
	go get github.com/tebeka/go2xunit
	go get google.golang.org/grpc
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/mmcc007/go/examples/route_guide/routeguide
	go build -o examples/route_guide/server/server examples/route_guide/server/server.go
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen
	go test -v ./... > test.output
	cat test.output | ~/go/bin/go2xunit -output tests.xml
