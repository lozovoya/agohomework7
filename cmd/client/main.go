package main

import (
	avia "agohomework7/pkg/avia/v1"
	"context"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(2)
	}
}

func execute(addr string) (err error) {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	var t = time.Date(2021, time.May, 05, 6, 0, 0, 0, time.Local)
	departure, err := ptypes.TimestampProto(t)
	var request = avia.TicketRequest{
		From: "SVO",
		To:   "BER",
		Data: departure,
	}
	client := avia.NewAviaServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*600)
	stream, err := client.AviaTickets(ctx, &request)

	for {
		responce, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		t, _ := ptypes.Timestamp(responce.DepTime)
		log.Printf("flightid %d, Departure time %d:%d, Duration %dh, Price %d rubles",
			responce.Id, t.Hour(), t.Minute(), responce.Duration/60, responce.Price/100)
	}
	return nil
}
