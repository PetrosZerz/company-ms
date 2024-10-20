package adapters

import (
	"company-ms/internal"
	"company-ms/internal/application"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database, collectionName string) application.CompanyRepository {
	collection := db.Collection(collectionName)
	return &MongoRepository{collection}
}

func (r *MongoRepository) Create(company *internal.Company) error {
	_, err := r.collection.InsertOne(context.TODO(), company)
	return err
}

func (r *MongoRepository) GetAll() ([]*internal.Company, error) {
	var companies []*internal.Company

	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var company internal.Company
		if err := cursor.Decode(&company); err != nil {
			return nil, err
		}
		companies = append(companies, &company)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *MongoRepository) GetByID(id string) (*internal.Company, error) {
	var company internal.Company

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("company not found")
		}
		return nil, err
	}
	return &company, nil
}

func (r *MongoRepository) GetByName(name string) error {
	var company internal.Company
	filter := bson.M{"name": name}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			return err
		}
	}
	return fmt.Errorf("company with name %v already defined", name)
}

func (r *MongoRepository) GetByNameAndId(id string, name string) error {
	var company internal.Company
	filter := bson.M{"_id": bson.M{"$ne": id}, "name": name}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			return err
		}
	}
	return fmt.Errorf("company with name %v already defined", name)
}

func (r *MongoRepository) Update(company *internal.Company) error {
	filter := bson.M{"_id": company.ID}
	update := bson.M{
		"$set": bson.M{
			"name":                company.Name,
			"description":         company.Description,
			"amount_of_employees": company.AmountOfEmployees,
			"registered":          company.Registered,
			"type":                company.Type,
		},
	}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoRepository) Delete(id string) error {
	filter := bson.M{"_id": id}
	results, err := r.collection.DeleteOne(context.TODO(), filter)
	if results.DeletedCount == 0 {
		return errors.New("no document deleted")
	}
	return err
}
