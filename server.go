package gss

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Server struct {
	Port uint
	Mux  *http.ServeMux
}

func NewServer(port uint) *Server {
	return &Server{
		Port: port,
		Mux:  http.NewServeMux(),
	}
}

func (s *Server) RegisterHandle(path string, h http.Handler) {
	s.Mux.Handle(path, h)
}

func (s *Server) Start(ctx context.Context, wg *sync.WaitGroup) {
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(int(s.Port)),
		Handler: s.Mux,
	}

	log.Println("Starting subscription server on " + srv.Addr)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				if err := srv.Shutdown(ctx); err != nil {
					log.Fatal(err)
				}
				return
			}
		}
	}()
}
