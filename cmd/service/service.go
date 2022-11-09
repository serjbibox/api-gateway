package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/serjbibox/api-gateway/pkg/handler"
	"github.com/serjbibox/api-gateway/pkg/storage"
	"github.com/serjbibox/api-gateway/pkg/storage/_kafka"
)

var elog = log.New(os.Stderr, "service error\t", log.Ldate|log.Ltime|log.Lshortfile)
var ilog = log.New(os.Stdout, "service info\t", log.Ldate|log.Ltime)
var ctx = context.Background()

const (
	HTTP_PORT = "8080"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func main() {
	var err error
	db, err := _kafka.New(
		ctx,
		[]string{"192.168.0.109:29092"},
		"news-request-stream",
		"api-gateway",
	)
	if err != nil {
		elog.Fatal(err)
	}

	s, err := storage.NewStorageKafka(ctx, db)
	if err != nil {
		elog.Fatal(err)
	}
	handlers, err := handler.New(s)
	if err != nil {
		elog.Fatal(err)
	}
	srv := new(Server)
	elog.Fatal(srv.Run(HTTP_PORT, handlers.InitRoutes()))
}
