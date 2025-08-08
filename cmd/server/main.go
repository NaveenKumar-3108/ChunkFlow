package main

import (
	"CHUNKFLOW/internal/db"
	router "CHUNKFLOW/routes"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("Server started...")
	f, err := os.OpenFile("././log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	defer f.Close()
	db.Init()
	defer db.GDB.Close()
	Router := router.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", Router))

}
