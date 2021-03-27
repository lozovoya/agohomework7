package app

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/jackc/pgx/v4/pgxpool"
	avia "github.com/lozovoya/agohomework7/pkg/avia/v1"
	"log"
	"sync"
	"time"
)

type Server struct {
	Databases []*DataBase
}

type DataBase struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

type Flight struct {
	Id             int64
	Departure_time string
	Duration       int64
	Price          int64
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) AddDb(ctx context.Context, pool *pgxpool.Pool) error {
	var db = DataBase{
		ctx:  ctx,
		pool: pool,
	}
	s.Databases = append(s.Databases, &db)
	return nil
}

func (s *Server) AviaTickets(request *avia.TicketRequest, server avia.AviaService_AviaTicketsServer) error {

	log.Println(request.From, request.To, request.Data)
	timestamp, _ := ptypes.Timestamp(request.Data)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var interval = 2 * time.Second

	for _, dbc := range s.Databases {

		wg.Add(1)
		go func(dbc *DataBase, interval time.Duration) {
			log.Println("goroutine started")
			defer wg.Done()
			time.Sleep(interval)
			rows, err := dbc.pool.Query(dbc.ctx,
				"SELECT id, departure, duration, price FROM flights WHERE iata_from=$1 and iata_to=$2 and departure_date=$3",
				request.From, request.To, timestamp)
			if err != nil {
				log.Println(err)
				return
			}
			defer rows.Close()

			var responce avia.TicketResponce
			var departureTime time.Time
			for rows.Next() {
				err = rows.Scan(&responce.Id, &departureTime, &responce.Duration, &responce.Price)
				if err != nil {
					log.Println(err)
					return
				}
				responce.DepTime, err = ptypes.TimestampProto(departureTime)
				if err != nil {
					log.Println(err)
					return
				}
				mu.Lock()
				if err := server.Send(&responce); err != nil {
					log.Println(err)
					return
				}
				mu.Unlock()
			}
			log.Println("goroutine finished")
		}(dbc, interval)
		interval = interval + 10*time.Second

	}
	wg.Wait()
	return nil
}
