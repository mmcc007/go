pipeline {
  agent any
  /*
  environment {
        CC = 'clang'
	// Install the desired Go version
    	def root = tool name: 'Go 1.8.1', type: 'go'

    	// Export environment variables pointing to the directory where Go was installed
    	GOROOT=${root}
	PATH='${root}/bin'
	GO='${root}/bin'
    }
    */
  stages {
    stage('build') {
      steps {
        sh "/usr/local/bin/go get github.com/tebeka/go2xunit"
	sh "/usr/local/bin/go get google.golang.org/grpc"
	sh "/usr/local/bin/go get -u github.com/golang/protobuf/proto"
	sh "/usr/local/bin/go get -u github.com/golang/protobuf/protoc-gen-go"
	sh "/usr/local/bin/go get github.com/mmcc007/go/examples/route_guide/routeguide"
	sh "/usr/local/bin/go build -o examples/route_guide/server/server examples/route_guide/server/server.go"
	sh "/usr/local/bin/go get github.com/golang/mock/gomock"
	sh "/usr/local/bin/go get github.com/golang/mock/mockgen"
	sh "/usr/local/bin/go test -v ./... > test.output"
	sh "cat test.output | ~/go/bin/go2xunit -output tests.xml"
      }
    }
  }
}
