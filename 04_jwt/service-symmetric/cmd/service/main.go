package main

import (
	"context"
	"encoding/base64"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"net/http"
	"service/cmd/service/app"
	"time"

	//"context"
	//"github.com/jackc/pgx/v4/pgxpool"
	//"log"
	//"net"
	//"net/http"
	//"os"
	//"service/cmd/service/app"
	//"service/pkg/business"
	//"service/pkg/security"
	"net"
	"os"
	"service/pkg/business"
	"service/pkg/security"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@localhost:5432/db"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = defaultDSN
	}

	tokenLifeTime := time.Hour

	keyBytes, err := ioutil.ReadFile("keys/symmetric.key")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	secretKey, err := base64.StdEncoding.DecodeString(string(keyBytes))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	if err := execute(net.JoinHostPort(host, port), dsn, secretKey, tokenLifeTime); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string, dsn string, secretKey []byte, tokenLifeTime time.Duration) error {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Print(err)
		return err
	}
	defer pool.Close()

	securitySvc := security.NewService(pool, secretKey, tokenLifeTime)
	businessSvc := business.NewService(pool)
	router := chi.NewRouter()
	application := app.NewServer(securitySvc, businessSvc, router)
	err = application.Init()
	if err != nil {
		log.Print(err)
		return err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
