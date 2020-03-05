package main

import (
	"burger-shop.tj/cmd/crud/app"
	"burger-shop.tj/pkg/crud/services"
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
	log.Println("host setting to connect")
	host, ok := services.FlagOrEnv(hostF, ENV_HOST)
	if !ok {
		log.Panic("can't host setting")
	}
	log.Println("get port to connect")
	port, ok := services.FlagOrEnv(portF, ENV_PORT)
	if !ok {
		log.Panic("can't port setting")
	}
	log.Println("set address to connect")
	addr := net.JoinHostPort(host, port)
	log.Printf("address to connect: %s", addr)

	log.Println("set database to connect")
	dsn, ok := services.FlagOrEnv(dsnF, ENV_DSN)
	if !ok {
		log.Panic("How get DB url?")
	}
	log.Printf("try start server on: %s, dbUrl: %s", addr, dsn)
	start(addr, dsn)
	log.Printf("server success on: %s, dbUrl: %s", addr, dsn)
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
