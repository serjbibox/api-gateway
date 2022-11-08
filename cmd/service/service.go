package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
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
		[]string{"192.168.0.109:29092"},
		"news-stream",
		0,
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
	id := 1
	go readTopic(id, db)
	//id++
	//go readTopic(id, db)
	elog.Fatal(srv.Run(HTTP_PORT, handlers.InitRoutes()))
}

func readTopic(id int, db *_kafka.Client) {
	conn, err := kafka.Dial("tcp", "192.168.0.109:29092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}
	/*db.Reader.SetOffset(30)
	for {
		msg, err := db.FetchProcessCommit()
		if err != nil {
			log.Println("reading kafka err", id, err)
		}
		log.Printf("%+v\n%s\n%s\nid: %d", msg, string(msg.Key), string(msg.Value), id)
		time.Sleep(100 * time.Millisecond)
	}*/
}
