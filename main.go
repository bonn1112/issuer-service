package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/lastrust/issuing-service/protocol"
	"github.com/lastrust/issuing-service/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const addr = ":8082"

func main() {
	privateKey := os.Getenv("BLOCKCHAIN_PRIVATE_KEY")
	if privateKey == "" {
		logrus.Fatal(errors.New("private key couldn't be empty string"))
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	protocol.RegisterCertIssuerServer(s, service.New())

	if os.Getenv("PROCESS_ENV") == "dev" {
		logrus.Info("reflection GRPC is registered")
		reflection.Register(s)
	}

	if err = s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
