package main

import (
	"context"
	"flag"
	"github.com/JAbduvohidov/burger-shop.tj/cmd/crud/app"
	"github.com/JAbduvohidov/burger-shop.tj/pkg/crud/services"
	"github.com/JAbduvohidov/burger-shop.tj/pkg/crud/services/burgers"
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
	portF = flag.String("port", "8080", "Server port")
	dsnF  = flag.String("dsn", "postgres://nurzyxgxduryxt:a91147d43b56869a99a0815d324323f5f22071d6dfa17cdd789c93388a392072@ec2-52-86-73-86.compute-1.amazonaws.com:5432/dc1ns5rpr9g4e3", "Postgres DSN")
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
		log.Panic("unable to get db url")
	}

	log.Printf("try start server on: %s, dbUrl: %s", addr, dsn)
	start(addr, dsn)
	log.Printf("server success on: %s, dbUrl: %s", addr, dsn)
}

func start(addr string, dsn string) {
	router := app.NewExactMux()

	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	burgersSvc := burgers.NewBurgersSvc(pool)
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
