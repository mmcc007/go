package main

import (
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/mmccabe/go/examples/route_guide/routeguide"
)

var client pb.RouteGuideClient
var conn *grpc.ClientConn
var serverCmd *exec.Cmd

func TestMain(m *testing.M) {
	grpclog.Printf("TestMain()")
	startServer()
	client, conn = startClient()
	returnCode := m.Run()
	stopClient()
	stopServer()
	os.Exit(returnCode)
}
func startServer() {
	grpclog.Printf("startServer()")
	cmdStr := "server/server"
	serverCmd = exec.Command(cmdStr)
	serverCmd.Dir = ".."
	err := serverCmd.Start()
	if err != nil {
		grpclog.Fatal("Server failed to start: ", err)
	}

	// wait for port to open
	for i := 0; i < 50; i++ {
		if checkServerUp() {
			grpclog.Println("server up!")
			break
		}
		grpclog.Println("server not up yet...trying again!")
		time.Sleep(10 * time.Millisecond)
	}

}

func checkServerUp() bool {
	// Check if server port is in use

	// Try to create a server with the port
	//server, err := net.Listen("tcp", *serverAddr)
	server, err := net.Listen("tcp", ":10000")

	// if it fails then the port is likely taken
	if err != nil {
		return true
	}

	server.Close()

	// we successfully used and closed the port
	// so it's now available to be used again
	return false

}

func stopServer() {
	grpclog.Printf("stopServer()")
	if err := serverCmd.Process.Kill(); err != nil {
		grpclog.Fatal("failed to kill: ", err)
	}
}

func startClient() (pb.RouteGuideClient, *grpc.ClientConn) {
	grpclog.Printf("startClient()")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewRouteGuideClient(conn)

	return client, conn
}

func stopClient() {
	conn.Close()
}

func TestPrintFeature(t *testing.T) {
	grpclog.Printf("TestPrintFeature()")
	point := &pb.Point{Latitude: 409146138, Longitude: -746188906}
	printFeature(client, point)
}

func TestPrintFeatures(t *testing.T) {
	printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})
}

func TestRunRecordRoutes(t *testing.T) {
	grpclog.Printf("TestRunRecordRoutes()")
	runRecordRoute(client)
}

func TestRunRouteChat(t *testing.T) {
	runRouteChat(client)
}
