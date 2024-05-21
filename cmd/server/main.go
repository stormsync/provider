package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	slogenv "github.com/cbrewster/slog-env"
	"github.com/jackc/pgx/v5"
	"github.com/stormsync/database"

	api "github.com/jason-costello/weather/accesssvc/api/go"
	"github.com/jason-costello/weather/accesssvc/consumer"
)

func main() {
	logger := slog.New(slogenv.NewHandler(slog.NewTextHandler(os.Stderr, nil)))
	var groupID = "transform-consumer"

	dbAddress := os.Getenv("DB_ADDRESS")
	if dbAddress == "" {
		log.Fatal("dbAddress is required.  Use env var DB_ADDRESS")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("dbUser is required.  Use env var DB_USER")
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		log.Fatal("dbPass is required.  Use env var DB_PASS")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("dbName is required.  Use env var DB_NAME")
	}

	webServerAddress := os.Getenv("WEB_SERVER_ADDRESS")
	if webServerAddress == "" {
		log.Fatal("webServerAddress is required.  Use env var WEB_SERVER_ADDRESS")
	}

	apikey := os.Getenv("API_KEY")
	if apikey == "" {
		log.Fatal("apikey is required.  Use env var API_KEY")
	}

	address := os.Getenv("KAFKA_ADDRESS")
	if address == "" {
		log.Fatal("address is required.  Use env var ADDRESS")
	}

	user := os.Getenv("KAFKA_USER")
	if user == "" {
		log.Fatal("kafka user is required.  Use env var KAFKA_USER")
	}

	pw := os.Getenv("KAFKA_PASSWORD")
	if pw == "" {
		log.Fatal("kafka password is required.  Use env var KAFKA_PASSWORD")
	}

	consumerTopic := os.Getenv("CONSUMER_TOPIC")
	if consumerTopic == "" {
		log.Fatal("consumer topic is required.  Use env var CONSUMER_TOPIC")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := pgx.Connect(ctx, "user=postgres dbname=weather sslmode=disable")
	if err != nil {
		log.Fatal("no db", err)
	}
	defer conn.Close(ctx)
	db := database.New(conn)

	consumer, err := consumer.NewConsumer(address, consumerTopic, user, pw, groupID, logger, db)
	if err != nil {
		log.Fatal("unable to create consumer: ", err)
	}

	if err != nil {
		log.Fatal("failed to create the collect: %w", err)
	}

	logger.Info("Starting transform service")

	rc := api.RouterConfig{
		ROKey:  "rokey",
		RWKey:  "rwkey",
		DB:     db,
		Logger: logger,
	}
	sdb := api.NewRouter(rc)
	go func() {
		sdb.Web.Logger.Fatal(sdb.Web.Start("0.0.0.0:8080"))
	}()

	for {
		if err := consumer.GetMessage(ctx); err != nil {
			logger.Error("failed to collect message: ", "error", err)
			cancel()
			time.Sleep(10 * time.Second)
			break
		}
	}
}
