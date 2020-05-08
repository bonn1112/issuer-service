package main

import (
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/kelseyhightower/envconfig"
	"github.com/lastrust/issuing-service/env"
	"github.com/lastrust/issuing-service/infra/database"
	"github.com/lastrust/issuing-service/infra/repos/certrepo"
	"github.com/lastrust/issuing-service/infra/repos/issuerrepo"
	"github.com/lastrust/issuing-service/protocol"
	"github.com/lastrust/issuing-service/service"
	"github.com/lastrust/utils-go/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var conf env.Config
	if err := envconfig.Process("", &conf); err != nil {
		logging.Err().Fatal(err)
	}

	lis, err := net.Listen("tcp", conf.ServerAddr)
	if err != nil {
		logging.Err().WithError(err).Fatalln("failed to listen")
	}

	logOpts := configureLogger(conf.LogLevel)
	srv := grpc.NewServer(logOpts...)

	db, err := database.Open(conf.DB)
	if err != nil {
		logging.Err().WithError(err).Fatalln("error database initialization")
	}
	svc := service.New(conf.Service, issuerrepo.New(db), certrepo.New(db))

	protocol.RegisterIssuingServiceServer(srv, svc)

	if conf.Service.ProcessEnv == "dev" {
		logging.Out().Info("reflection GRPC is registered")
		reflection.Register(srv)
	}

	logging.Out().Printf("Listening GRPC on %s", conf.ServerAddr)
	if err = srv.Serve(lis); err != nil {
		logging.Err().WithError(err).Fatalln("failed to serve")
	}
}

func configureLogger(lvl string) []grpc.ServerOption {
	err := logging.Init(lvl, os.Stdout, os.Stderr)
	if err != nil {
		logging.Err().WithError(err).Fatalln("failed logging initialization")
	}
	logging.Out().Printf("logging level: %s\n", lvl)

	return []grpc.ServerOption{
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logging.Out()))),
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logging.Out()))),
	}
}
