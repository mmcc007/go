package main

import (
	"context"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/mmcc007/go/examples/route_guide/routeguide"
)

var client pb.RouteGuideClient
var conn *grpc.ClientConn
var serverCmd *exec.Cmd

func TestMain(m *testing.M) {
	grpclog.Printf("TestMain()")
	startServer()
	//startServer2()
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

	if !serverUp() {
		grpclog.Fatal("Server failed to open port")
	}

}

func startServer2() {
	go main()
	if !serverUp() {
		grpclog.Fatal("Server failed to open port")
	}
}

func serverUp() bool {
	// wait for port to open
	for i := 0; i < 100; i++ {
		if checkServerUp() {
			grpclog.Println("server up!")
			return true
		}
		grpclog.Println("server not up yet...trying again!")
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

func checkServerUp() bool {
	// Check if server port is in use

	// Try to create a server with the port
	server, err := net.Listen("tcp", ":10000")

	// if it fails then the port is likely taken
	if err != nil {
		return true
	}

	err = server.Close()
	if err != nil {
		return true
	}

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
	//conn, err := grpc.Dial("127.0.0.1:"+string(*port), opts...)
	conn, err := grpc.Dial("127.0.0.1:10000", opts...)
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
	grpclog.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	feature, err := client.GetFeature(context.Background(), point)
	if err != nil {
		grpclog.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	grpclog.Println(feature)
}

func TestPrintFeatures(t *testing.T) {
	// Looking for features between 40, -75 and 42, -73.
	rect := &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	}
	grpclog.Printf("Looking for features within %v", rect)
	stream, err := client.ListFeatures(context.Background(), rect)
	if err != nil {
		grpclog.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpclog.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		grpclog.Println(feature)
	}
}

func TestRunRecordRoutes(t *testing.T) {
	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	grpclog.Printf("Traversing %d points.", len(points))
	stream, err := client.RecordRoute(context.Background())
	if err != nil {
		grpclog.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			grpclog.Fatalf("%v.Send(%v) = %v", stream, point, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		grpclog.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	grpclog.Printf("Route summary: %v", reply)
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}

// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func TestRunRouteChat(t *testing.T) {
	notes := []*pb.RouteNote{
		{&pb.Point{Latitude: 0, Longitude: 1}, "First message"},
		{&pb.Point{Latitude: 0, Longitude: 2}, "Second message"},
		{&pb.Point{Latitude: 0, Longitude: 3}, "Third message"},
		{&pb.Point{Latitude: 0, Longitude: 1}, "Fourth message"},
		{&pb.Point{Latitude: 0, Longitude: 2}, "Fifth message"},
		{&pb.Point{Latitude: 0, Longitude: 3}, "Sixth message"},
	}
	stream, err := client.RouteChat(context.Background())
	if err != nil {
		grpclog.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				grpclog.Fatalf("Failed to receive a note : %v", err)
			}
			grpclog.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			grpclog.Fatalf("Failed to send a note: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}
