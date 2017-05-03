pipeline {
  agent any
  environment {
        CC = 'clang'
	// Install the desired Go version
    	def root = tool name: 'Go 1.8.1', type: 'go'

    	// Export environment variables pointing to the directory where Go was installed
    	GOROOT=${root}
	PATH+GO=${root}/bin
        sh 'go version'
    }
  stages {
    stage('build') {
      steps {
        sh "go get github.com/tebeka/go2xunit"
	sh "go get google.golang.org/grpc"
	sh "go get -u github.com/golang/protobuf/proto"
	sh "go get -u github.com/golang/protobuf/protoc-gen-go"
	sh "go get github.com/mmcc007/go/examples/route_guide/routeguide"
	sh "go build -o examples/route_guide/server/server examples/route_guide/server/server.go"
	sh "go get github.com/golang/mock/gomock"
	sh "go get github.com/golang/mock/mockgen"
	sh "go test -v ./... > test.output"
	sh "cat test.output | ~/go/bin/go2xunit -output tests.xml"
      }
    }
  }
}
