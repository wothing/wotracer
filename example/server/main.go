/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/06/06 09:48
 */

package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/wothing/wotracer"
	"github.com/wothing/wotracer/example/pb"
)

const (
	grpcPort  = ":10010"
	debugPort = ":10011"
)

func main() {
	wotracer.InitTracer(":1708")

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}

	log.Printf("starting grpc at %s", grpcPort)

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &helloServer{})
	go s.Serve(lis)

	log.Printf("starting debug at %s", debugPort)
	log.Fatal(http.ListenAndServe(debugPort, nil))
}

type helloServer struct {
}

func (helloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	span, _ := wotracer.JoinRPC(ctx, "SayHello")
	defer span.Finish()

	//mock delay 100ms
	time.Sleep(time.Millisecond * 100)

	return &pb.HelloResponse{Reply: "Hello, " + req.Greeting}, nil
}

func (helloServer) SayGoodbye(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	span, ctx := wotracer.JoinRPC(ctx, "SayGoodbye")
	defer span.Finish()

	//mock delay 150ms
	time.Sleep(time.Millisecond * 150)

	{
		conn, err := grpc.Dial(":10010", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		client := pb.NewHelloServiceClient(conn)
		client.SayGoodnight(wotracer.PackCtx(ctx), &pb.HelloRequest{Greeting: "elvizlai"})
	}

	return &pb.HelloResponse{Reply: "Goodbye, " + req.Greeting}, nil
}

func (helloServer) SayGoodnight(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	span, _ := wotracer.JoinRPC(ctx, "SayGoodnight")
	defer span.Finish()

	//mock delay 200ms
	time.Sleep(time.Millisecond * 200)

	return &pb.HelloResponse{Reply: "Goodnight, " + req.Greeting}, nil
}
