package grpcs

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/meowalien/go-meowalien-lib/errs"
)

func StartGRPCServer() (grpcServer *grpc.Server) {
	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
	grpcServer = grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	return grpcServer
}

func ListenAndServe(grpcServer *grpc.Server, grpcServerListen string) (err error) {
	listener, err := net.Listen("tcp", grpcServerListen)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcServerListen, err)
	}
	fmt.Println("Start listening on port", grpcServerListen)
	err = grpcServer.Serve(listener)
	if err != nil {
		err = errs.New(err)
		return
	}
	return nil
}
