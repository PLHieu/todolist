package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todolist/db"
	"todolist/handlers"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	port, ok := os.LookupEnv("TODOLIST_PORT")
	if !ok {
		log.Fatal("$TODOLIST_PORT not found")
	}
	postgresURL, ok := os.LookupEnv("TODOLIST_DB_URL")
	if !ok {
		log.Fatal("$TODOLIST_DB_URL not found")
	}
	m, err := migrate.New("file://db/migrations", postgresURL)
	if err != nil {
		log.Fatalf("migrate:%v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration up:%v", err)
	}
	// defer func() {
	// 	m.Down()
	// }()
	dbConn, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("db open:%v", err)
	}

	ns := db.NewNotificationService()
	ts := db.NewTodoService(dbConn, ns)
	us := db.NewUserService(dbConn, ts)
	h := handlers.NewHandler(ts, us)

	router := handlers.ConfigureServer(h)
	log.Printf("Listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), router))
}
