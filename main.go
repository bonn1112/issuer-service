package main

import (
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/lastrust/issuing-service/protocol"
	"github.com/lastrust/issuing-service/service"
	"github.com/lastrust/issuing-service/utils/env"
	"github.com/lastrust/utils-go/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	processEnv     = env.GetDefault("PROCESS_ENV", "dev")
	logLevelString = env.GetDefault("LOG_LEVEL", "info")
	cloudService   = env.GetDefault("CLOUD_SERVICE", "gcp")
	port           = ":8080"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logging.Err().WithError(err).Fatalln("failed to listen")
	}

	logOpts := configureLogger()

	srv := grpc.NewServer(logOpts...)
	protocol.RegisterIssuingServiceServer(srv, service.New(cloudService, processEnv))

	if processEnv == "dev" {
		logging.Out().Info("reflection GRPC is registered")
		reflection.Register(srv)
	}

	logging.Out().Printf("Listening GRPC on %s", port)
	if err = srv.Serve(lis); err != nil {
		logging.Err().WithError(err).Fatalln("failed to serve")
	}
}

func configureLogger() []grpc.ServerOption {
	err := logging.Init(logLevelString, os.Stdout, os.Stderr)
	if err != nil {
		logging.Err().WithError(err).Fatalln("failed logging initialization")
	}
	logging.Out().Printf("logging level: %s\n", logLevelString)

	return []grpc.ServerOption{
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logging.Out()))),
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logging.Out()))),
	}
}
