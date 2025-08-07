package initializers

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v3/log"
	_ "github.com/lib/pq"
	"os"
	"time"
)

var DB *sql.DB

func LoadDatabase() {
	var err error
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("postgres", connectionString)
		if err == nil && DB.Ping() == nil {
			log.Info("✅ Successfully connected to Postgres")
			return
		}
		log.Info("⏳ Waiting for database to be ready...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("❌ Failed to connect to database after several attempts:", err)
}
