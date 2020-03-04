package main

import (
	"net"
	"os"

	"github.com/lastrust/issuing-service/config"
	"github.com/lastrust/issuing-service/protocol"
	"github.com/lastrust/issuing-service/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logFile *os.File

func main() {
	conf, err := config.Env()
	if err != nil {
		logrus.WithError(err).Fatalln("invalid configuration")
	}
	configureLogger(conf.LogFilename)

	lis, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to listen")
	}

	s := grpc.NewServer()
	protocol.RegisterIssuingServiceServer(s, service.New())

	if conf.ProcessEnv == "dev" {
		logrus.Info("reflection GRPC is registered")
		reflection.Register(s)
	}

	logrus.Printf("Listening and serving GRPC on %s\n", conf.Addr)
	if err = s.Serve(lis); err != nil {
		logrus.WithError(err).Fatalln("failed to serve")
	}
}

func configureLogger(logFilename string) {
	if logFilename == "" {
		logrus.Warn("logging filename is empty")
		return
	}

	var err error
	logFile, err = os.OpenFile(logFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		logrus.WithError(err).Fatalln("error opening log file")
		return
	}

	logrus.SetOutput(logFile)
}
