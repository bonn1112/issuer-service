package main

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/lastrust/issuing-service/config"
	"github.com/lastrust/issuing-service/protocol"
	"github.com/lastrust/issuing-service/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := config.Env()
	if err != nil {
		logrus.WithError(err).Fatalln("invalid configuration")
	}

	lis, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to listen")
	}

	logOpts := configureLogger(conf.LogLevel)

	srv := grpc.NewServer(logOpts...)
	protocol.RegisterIssuingServiceServer(srv, service.New())

	if conf.ProcessEnv == "dev" {
		logrus.Info("reflection GRPC is registered")
		reflection.Register(srv)
	}

	logrus.Printf("Listening and serving GRPC on %s", conf.Addr)
	if err = srv.Serve(lis); err != nil {
		logrus.WithError(err).Fatalln("failed to serve")
	}
}

func configureLogger(logLevelString string) []grpc.ServerOption {
	logLevel, err := logrus.ParseLevel(logLevelString)
	if err != nil {
		logrus.WithError(err).Fatalln("failed parse log level")
	}
	logrus.Printf("Log level: %d %s", logLevel, logLevelString)

	logger := logrus.New()
	logger.Level = logLevel

	return []grpc.ServerOption{
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logger))),
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger))),
	}
}
