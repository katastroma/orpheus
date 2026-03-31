//revive:disable:package-comments
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	foundationclient "git.sonicoriginal.software/grpc-foundation/client"
	foundationotel "git.sonicoriginal.software/grpc-foundation/otel"
	"git.sonicoriginal.software/grpc-foundation/server"

	pb "github.com/katastroma/keleustes"

	"github.com/katastroma/orpheus/internal/config"
	grpc_health "github.com/katastroma/orpheus/internal/health/grpc"
	http_health "github.com/katastroma/orpheus/internal/health/http"
	"github.com/katastroma/orpheus/internal/orderer"
	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/render/helm"
	"github.com/katastroma/orpheus/internal/render/kustomize"
	"github.com/katastroma/orpheus/internal/render/plain"
	"github.com/katastroma/orpheus/internal/serve"
)

const shutdownTimeout = 10 * time.Second

func main() {
	mainCtx := context.Background()

	log, shutdown, err := foundationotel.Init(mainCtx, "orpheus")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize otel: %v\n", err)
		os.Exit(1)
	}
	defer shutdown(mainCtx)

	ordererAddr, err := config.RequireEnv("ORDERER_ADDR")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	ordererConn, err := foundationclient.New(ordererAddr, nil, nil)
	if err != nil {
		log.Error("orderer connection failed", "error", err)
		os.Exit(1)
	}
	defer ordererConn.Close()

	var router render.Router
	router.Register(helm.Type, helm.Match, helm.Render)
	router.Register(kustomize.Type, kustomize.Match, kustomize.Render)
	router.Register(plain.Type, plain.Match, plain.Render)

	streamFn := orderer.NewStreamFunc(log, ordererConn)
	service := serve.New(log, router.Render, streamFn)

	mux := http.NewServeMux()
	mux.Handle("GET /healthz", http_health.New(log))

	httpPort := config.StringEnv("PORT", "8080")
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: mux,
	}

	grpcServer := server.New(log)
	healthpb.RegisterHealthServer(grpcServer, grpc_health.New())
	pb.RegisterRendererServiceServer(grpcServer, service)

	sigNotifyContext, stop := context.WithCancel(mainCtx)
	defer stop()

	var wg sync.WaitGroup

	wg.Go(func() {
		log.Info("starting http server", "port", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server error", "error", err)
			stop()
		}
	})

	wg.Go(func() {
		addr := server.Address()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("grpc listen error", "error", err)
			stop()
			return
		}
		log.Info("starting grpc server", "address", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server error", "error", err)
			stop()
		}
	})

	server.HandleGracefulShutdown(sigNotifyContext, stop, log, grpcServer, shutdownTimeout)

	shutdownCtx, shutdownCancel := context.WithTimeout(mainCtx, shutdownTimeout)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown error", "error", err)
	}

	wg.Wait()
}
