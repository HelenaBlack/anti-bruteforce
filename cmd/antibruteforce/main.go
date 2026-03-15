package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/HelenaBlack/anti-bruteforce/api/gen"          //nolint:depguard
	"github.com/HelenaBlack/anti-bruteforce/internal/app"        //nolint:depguard
	"github.com/HelenaBlack/anti-bruteforce/internal/config"     //nolint:depguard
	"github.com/HelenaBlack/anti-bruteforce/internal/limiter"    //nolint:depguard
	"github.com/HelenaBlack/anti-bruteforce/internal/repository" //nolint:depguard
	"github.com/HelenaBlack/anti-bruteforce/internal/server"     //nolint:depguard
	"github.com/redis/go-redis/v9"                               //nolint:depguard
	"google.golang.org/grpc"                                     //nolint:depguard
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func() { _ = db.Close() }()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer func() { _ = rdb.Close() }()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("failed to connect to redis: %v", err)
		return
	}

	ipRepo := repository.NewPostgresIPRepository(db)
	limitSvc := limiter.NewRedisLimiter(rdb)
	svc := app.NewAntiBruteforceService(limitSvc, ipRepo, cfg)
	grpcSrv := server.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterAntibruteforceServer(s, grpcSrv)

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	s.GracefulStop()
}
