package http_server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type server struct {
	srv *http.Server
}

func New(address string, handler http.Handler) *server {
	return &server{
		srv: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

func (s *server) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for {
			err := s.srv.ListenAndServe()
			if err == http.ErrServerClosed {
				return
			}
			if err != nil {
				log.Printf("http server listen: %v", err)
				// wait
				select {
				case <-ctx.Done(): // cancellation
					return
				case <-ticker.C:
					continue
				}
			}
		}
	}()
}

func (s *server) Stop() {
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		log.Printf("http server shutdown: %v", err)
	}
}
