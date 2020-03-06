package main

import (
	"context"
	"flag"
	"github.com/JAbduvohidov/burger-shop.tj/cmd/crud/app"
	"github.com/JAbduvohidov/burger-shop.tj/pkg/crud/services/burgers"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	hostF = flag.String("host", "0.0.0.0", "Server host")
	portF = flag.String("port", "8080", "Server port")
	dsnF  = flag.String("dsn", "postgres://nurzyxgxduryxt:a91147d43b56869a99a0815d324323f5f22071d6dfa17cdd789c93388a392072@ec2-52-86-73-86.compute-1.amazonaws.com:5432/dc1ns5rpr9g4e3", "Postgres DSN")
)

func main() {
	flag.Parse()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = *portF
	}
	log.Println("set address to connect")
	addr := net.JoinHostPort(*hostF, port)
	log.Printf("address to connect: %s", addr)

	log.Printf("try start server on: %s, dbUrl: %s", addr, *dsnF)
	start(addr, *dsnF)
	log.Printf("server success on: %s, dbUrl: %s", addr, *dsnF)
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
