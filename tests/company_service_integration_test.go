package tests

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"company-ms/internal"
	"company-ms/internal/adapters"
	"company-ms/internal/application"
)

const (
	mongoURI    = "mongodb://localhost:27017"
	kafkaBroker = "localhost:9092"
	kafkaTopic  = "company_topic"
)

func TestCreateCompanyIntegration(t *testing.T) {
	// Set up MongoDB client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Set up Kafka writer
	producer, err := adapters.NewKafkaProducer(kafkaBroker)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Set up Zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	repo := adapters.NewMongoRepository(client.Database("company_db"), "companies")

	// Create a company service instance
	companyService := application.NewCompanyService(repo, logger, producer)

	// Test company data
	newCompany := internal.Company{
		ID:                "uuid",
		Name:              "Test Company",
		Description:       "A company for testing",
		AmountOfEmployees: 50,
		Registered:        true,
		Type:              "Corporations",
	}

	// Perform the operation
	err = companyService.Create(&newCompany)
	assert.NoError(t, err, "Expected no error when creating company")

	// Verify the company was saved to MongoDB
	collection := client.Database("company_db").Collection("companies")
	var result internal.Company
	err = collection.FindOne(context.TODO(), bson.M{"_id": newCompany.ID}).Decode(&result)
	assert.NoError(t, err, "Expected no error when finding company in MongoDB")
	assert.Equal(t, newCompany, result, "The company returned should match the created company")

	message, err := json.Marshal(newCompany)
	assert.NoError(t, err, "Expected no error when marshaling company to JSON")

	// Produce the message
	err = producer.Produce("company_topic", message)
	assert.NoError(t, err, "Expected no error when producing message to Kafka")
}
