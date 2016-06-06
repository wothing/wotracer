/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/06/06 09:55
 */

package main

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"

	gctx "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/urfave/negroni"

	"github.com/wothing/wotracer"
	"github.com/wothing/wotracer/example/gateway/middleware"
	"github.com/wothing/wotracer/example/pb"
)

func main() {
	wotracer.InitTracer(":1708")

	router := mux.NewRouter()
	router.HandleFunc("/", Home)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(&middleware.LogMiddleware{})

	n.UseHandler(router)
	n.Run(":8699")
}

func Home(rw http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(":10010", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// mock delay 5ms
	time.Sleep(time.Millisecond * 5)

	client := pb.NewHelloServiceClient(conn)

	client.SayHello(gctx.Get(r, "ctx").(context.Context), &pb.HelloRequest{Greeting: "elvizlai"})

	// mock delay 10ms
	time.Sleep(time.Millisecond * 50)

	client.SayGoodbye(gctx.Get(r, "ctx").(context.Context), &pb.HelloRequest{Greeting: "elvizlai"})

	fmt.Fprintf(rw, `<p>GRPC requests have been made!<p>`)
	fmt.Fprintf(rw, `<p><a href="http://localhost:8700/traces" target="_">View the trace</a></p>`)
}
