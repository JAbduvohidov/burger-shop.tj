package main

import (
	"burger-shop.tj/cmd/crud/app"
	"burger-shop.tj/pkg/crud/services/burgers"
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"path/filepath"
)

const (
	ENV_HOST = "HOST"
	ENV_PORT = "PORT"
	ENV_DSN  = "DATABASE_URL"
)

var (
	hostF = flag.String("host", "0.0.0.0", "Server host")
	portF = flag.String("port", "9999", "Server port")
	dsnF  = flag.String("dsn", "", "Postgres DSN")
)

func main() {
	flag.Parse()
	addr := net.JoinHostPort(*hostF, *portF)
	log.Printf("address to connect: %s", addr)

	log.Printf("try start server on: %s, dbUrl: %s", addr, *dsnF)
	start(addr, *dsnF)
	log.Printf("server success on: %s, dbUrl: %s", addr, *dsnF)
}

func start(addr string, dsn string) {
	router := app.NewExactMux()
	log.Println("trying to create pool to connect")
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("can't create pool: %v", err)
		panic(err)
	}
	burgersSvc := burgers.NewBurgersSvc(pool)
	log.Println("server starting")
	server := app.NewServer(
		router,
		pool,
		burgersSvc,
		filepath.Join("web", "templates"),
		filepath.Join("web", "assets"),
	)
	server.InitRoutes()

	panic(http.ListenAndServe(addr, server))
}
