package main

import (
	"company-ms/internal/adapters"
	"company-ms/internal/application"
	"company-ms/internal/config"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.StacktraceKey = ""
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	return config.Build()
}

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	logger, err := setupLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// Mongo
	clientOptions := options.Client().ApplyURI(cfg.Database.Host + "://" + cfg.Database.Environment)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logger.Sugar().Fatal("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Ping MongoDB to check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	// Get MongoDB database and initialize repository
	db := client.Database(cfg.Database.Name)
	repo := adapters.NewMongoRepository(db, "companies")

	kafkaProd, err := adapters.NewKafkaProducer(cfg.Kafka.Host + ":" + cfg.Kafka.Port) // Set your Kafka broker address
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}
	defer kafkaProd.Close()

	companyService := application.NewCompanyService(repo, logger, kafkaProd)

	// Set up HTTP handlers
	companyHandler := adapters.NewCompanyHandler(companyService, logger)

	// Set up routes
	r := mux.NewRouter()
	r.Handle("/v1/companies", http.HandlerFunc(companyHandler.CreateCompany)).Methods("POST")
	r.Handle("/v1/companies", http.HandlerFunc(companyHandler.GetAllCompanies)).Methods("GET")
	r.Handle("/v1/companies/{id}", http.HandlerFunc(companyHandler.GetByID)).Methods("GET")
	r.Handle("/v1/companies/{id}", http.HandlerFunc(companyHandler.Update)).Methods("PATCH")
	r.Handle("/v1/companies/{id}", http.HandlerFunc(companyHandler.Delete)).Methods("DELETE")

	// Start the server
	logger.Sugar().Fatal(http.ListenAndServe(":8080", r))
}
