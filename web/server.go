package web

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Server struct {
	goSrv *http.Server
	ctx   context.Context
}

type Path struct {
	path    string
	handler http.HandlerFunc
}

func NewPath(urlpath string, handler http.HandlerFunc) Path {
	return Path{
		path:    urlpath,
		handler: handler,
	}
}

func NewServer(path string, port string, urls []Path) *Server {
	ctx := context.Background()
	router := RouterWithUrls("", urls)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(path, port),
		Handler: router.mux,
	}
	srv := &Server{
		goSrv: httpServer,
		ctx:   ctx,
	}
	return srv
}

func (srv *Server) Run() error {
	ctx, cancel := signal.NotifyContext(srv.ctx, os.Interrupt)
	defer cancel()
	go func() {
		log.Printf("listening on %s\n", srv.goSrv.Addr)
		if err := srv.goSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error in server: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := srv.goSrv.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		} else {
			log.Printf("stop listening on %s\n", srv.goSrv.Addr)
		}
	}()
	wg.Wait()
	return nil
}
