package operations

import (
	"fmt"
	"log"
	"tfg/internal/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Operation struct {
	ID          string  	`bson:"_id,omitempty"`
	Description string  	`bson:"description,omitempty"`
	Date        time.Time	`bson:"date,omitempty"`
	Amount      float64 	`bson:"amount,omitempty"`
	Category    string  	`bson:"category,omitempty"`
	UserID      string		`bson:"userId,omitempty"`
}

func (operation *Operation) Create() (error) {
	log.Printf("Create operation")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(operation.UserID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Execute insert
	res, err := mongo.InsertOne("operations", bson.D{
		{Key: "description", Value: operation.Description},
		{Key: "date", Value: operation.Date},
		{Key: "amount", Value: operation.Amount},
		{Key: "category", Value: operation.Category},
		{Key: "userId", Value: id},
	})
	if err != nil {
		log.Printf("Error: Create operation in db")
		return err
	}

	// Save ID
	operation.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func GetAllOperations(userId string) ([]Operation, error) {
	var operations []Operation

	log.Printf("Get all operations")

	// Query to get all operations
	query := bson.D{}

	if userId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return operations, err
		}

		// Change query to only search the user operations
		query = bson.D{
			{Key: "userId", Value: id},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("operations", query, filter)
	if err != nil {
		log.Printf("Error: Get all operations in db")
		return operations, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &operations); err != nil {
		log.Printf("Error: Decoding operations")
		return operations, err
	}

	return operations, nil
}

func GetOperationsByDate(userId string, initDate time.Time, endDate time.Time) ([]Operation, error) {
	var operations []Operation

	log.Printf("Get all operations by date")

	// Query to get all operations between dates
	query := bson.D{
		{Key: "date", Value: bson.M{"$gte": initDate}},
		{Key: "date", Value: bson.M{"$lte": endDate}},
	}

	if userId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return operations, err
		}

		// Change query to only search the user operations
		query = bson.D{
			{Key: "userId", Value: id},
			{Key: "date", Value: bson.M{"$gte": initDate}},
			{Key: "date", Value: bson.M{"$lte": endDate}},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("operations", query, filter)
	if err != nil {
		log.Printf("Error: Get all operations by date in db")
		return operations, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &operations); err != nil {
		log.Printf("Error: Decoding operations")
		return operations, err
	}

	return operations, nil
}

func GetOperationsByCategory(userId string, category string) ([]Operation, error) {
	var operations []Operation

	log.Printf("Get all operations by category")

	// Query to get all operations by category
	query := bson.D{
		{Key: "category", Value: category},
	}

	if userId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return operations, err
		}

		// Change query to only search the user operations
		query = bson.D{
			{Key: "userId", Value: id},
			{Key: "category", Value: category},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("operations", query, filter)
	if err != nil {
		log.Printf("Error: Get all operations by category in db")
		return operations, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &operations); err != nil {
		log.Printf("Error: Decoding operations")
		return operations, err
	}

	return operations, nil
}


func (operation *Operation) GetOperationById(userId string) error {
	log.Printf("Get operation by id")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(operation.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Query to get operation by id
	query := bson.D{
		{Key: "_id", Value: id},
	}

	if userId != "" {
		// String id to ObjectId
		userIdObject, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return err
		}

		// Change query to only search the user operations
		query = bson.D{
			{Key: "userId", Value: userIdObject},
			{Key: "_id", Value: id},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("operations", query, filter)
	if err != nil {
		log.Printf("Error: Get operation in db")
		return err
	}

	// Decode query
	cursor.Next(mongo.GetCtx())

	if err = cursor.Decode(&operation); err != nil {
		log.Printf("Error: Decoding operation")
		return err
	}

	return nil
}

func (operation *Operation) Update(userId string, description *string, date *time.Time, amount *float64, category *string) error {
	log.Printf("Update operation")

	// Get current state of operation
	operation.GetOperationById(userId)

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(operation.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	userID, err := primitive.ObjectIDFromHex(operation.UserID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Set fields if provided
	if description != nil {
		operation.Description = *description
	}
	if date != nil {
		operation.Date = *date
	}
	if amount != nil {
		operation.Amount = *amount
	}
	if category != nil {
		operation.Category = *category
	}

	// Filter to get user by id
	filter := bson.D{
		{Key: "_id", Value: id},
	}

	// Update fields
	update := bson.M{
		"description": operation.Description,
		"date": operation.Date,
		"amount": operation.Amount,
		"category": operation.Category,
		"userId": userID,
	}

	// Execute update
	result, err := mongo.ReplaceOne("operations", filter, update)
	if err != nil {
		log.Printf("Error: Update user in db")
		return err
	}
	if result.ModifiedCount == 0 {
		log.Printf("Error: Update user in db")
		return fmt.Errorf("no changes")
	}

	return nil
}

func Delete(operationId string, userId string) (bool, error) {
	log.Printf("Delete operation by id")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(operationId)
	if err != nil {
		log.Printf("Error: Convert string to id")
		return false, err
	}

	// Query to delete operation by id
	query := bson.D{
		{Key: "_id", Value: id},
	}

	if userId != "" {
		// String id to ObjectId
		userIdObject, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return false, err
		}

		// Change query to only search the user operations
		query = bson.D{
			{Key: "userId", Value: userIdObject},
			{Key: "_id", Value: id},
		}
	}

	// Execute query
	result, err := mongo.DeleteOne("operations", query)
	if err != nil {
		log.Printf("Error: Delete operation in db")
		return false, err
	}

	if result.DeletedCount == 0 {
		log.Printf("Error: Delete operation in db")
		return false, fmt.Errorf("no changes")
	}

	return true, nil
}