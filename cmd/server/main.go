package main

import (
	"agohomework7/cmd/server/app"
	avia "agohomework7/pkg/avia/v1"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"
const Db1 = "postgres://app:pass@localhost:5432/db1"
const Db2 = "postgres://app:pass@localhost:6432/db2"
const Db3 = "postgres://app:pass@localhost:7432/db3"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	db1, ok := os.LookupEnv("DB1")
	if !ok {
		db1 = Db1
	}

	db2, ok := os.LookupEnv("DB2")
	if !ok {
		db2 = Db2
	}

	db3, ok := os.LookupEnv("DB3")
	if !ok {
		db3 = Db3
	}

	var dbs = []string{db1, db2, db3}
	//log.Println(db2, db3)
	//var dbs = []string{db1}

	if err := execute(net.JoinHostPort(host, port), dbs); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string, dbs []string) error {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	server := app.NewServer()

	for _, db := range dbs {
		ctx := context.Background()
		pool, err := pgxpool.Connect(ctx, db)
		if err != nil {
			return err
		}
		defer pool.Close()

		server.AddDb(ctx, pool)
	}

	avia.RegisterAviaServiceServer(grpcServer, server)

	return grpcServer.Serve(listener)
}
