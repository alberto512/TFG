package transactions

import (
	"fmt"
	"log"
	"tfg/internal/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          string  	`bson:"_id"`
	Description string  	`bson:"description"`
	Date        time.Time	`bson:"date"`
	Amount      float64 	`bson:"amount"`
	Category    string  	`bson:"category"`
	AccountID   string		`bson:"accountId"`
}

func (transaction *Transaction) Create() (error) {
	log.Printf("Create transaction")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(transaction.AccountID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Execute insert
	res, err := mongo.InsertOne("transactions", bson.D{
		{Key: "description", Value: transaction.Description},
		{Key: "date", Value: transaction.Date},
		{Key: "amount", Value: transaction.Amount},
		{Key: "category", Value: transaction.Category},
		{Key: "accountId", Value: id},
	})
	if err != nil {
		log.Printf("Error: Create transaction in db")
		return err
	}

	// Save ID
	transaction.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func GetAllTransactions(accountId string) ([]Transaction, error) {
	var transactions []Transaction

	log.Printf("Get all transactions")

	// Query to get all transactions
	query := bson.D{}

	if accountId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(accountId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return transactions, err
		}

		// Change query to only search the account transactions
		query = bson.D{
			{Key: "accountId", Value: id},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("transactions", query, filter)
	if err != nil {
		log.Printf("Error: Get all transactions in db")
		return transactions, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &transactions); err != nil {
		log.Printf("Error: Decoding transactions")
		return transactions, err
	}

	return transactions, nil
}

func GetTransactionsByDate(accountId string, initDate time.Time, endDate time.Time) ([]Transaction, error) {
	var transactions []Transaction

	log.Printf("Get all transactions by date")

	// Query to get all transactions between dates
	query := bson.D{
		{Key: "date", Value: bson.M{"$gte": initDate}},
		{Key: "date", Value: bson.M{"$lte": endDate}},
	}

	if accountId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(accountId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return transactions, err
		}

		// Change query to only search the account transactions
		query = bson.D{
			{Key: "accountId", Value: id},
			{Key: "date", Value: bson.M{"$gte": initDate}},
			{Key: "date", Value: bson.M{"$lte": endDate}},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("transactions", query, filter)
	if err != nil {
		log.Printf("Error: Get all transactions by date in db")
		return transactions, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &transactions); err != nil {
		log.Printf("Error: Decoding transactions")
		return transactions, err
	}

	return transactions, nil
}

func GetTransactionsByCategory(accountId string, category string) ([]Transaction, error) {
	var transactions []Transaction

	log.Printf("Get all transactions by category")

	// Query to get all transactions by category
	query := bson.D{
		{Key: "category", Value: category},
	}

	if accountId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(accountId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return transactions, err
		}

		// Change query to only search the account transactions
		query = bson.D{
			{Key: "accountId", Value: id},
			{Key: "category", Value: category},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("transactions", query, filter)
	if err != nil {
		log.Printf("Error: Get all transactions by category in db")
		return transactions, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &transactions); err != nil {
		log.Printf("Error: Decoding transactions")
		return transactions, err
	}

	return transactions, nil
}

func GetTransactionsByCategoryAndDate(accountId string, category string, initDate time.Time, endDate time.Time) ([]Transaction, error) {
	var transactions []Transaction

	log.Printf("Get all transactions by category and date")

	// Query to get all transactions by category and date
	query := bson.D{
		{Key: "category", Value: category},
		{Key: "date", Value: bson.M{"$gte": initDate}},
		{Key: "date", Value: bson.M{"$lte": endDate}},
	}

	if accountId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(accountId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return transactions, err
		}

		// Change query to only search the account transactions
		query = bson.D{
			{Key: "accountId", Value: id},
			{Key: "category", Value: category},
			{Key: "date", Value: bson.M{"$gte": initDate}},
			{Key: "date", Value: bson.M{"$lte": endDate}},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("transactions", query, filter)
	if err != nil {
		log.Printf("Error: Get all transactions by category and date in db")
		return transactions, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &transactions); err != nil {
		log.Printf("Error: Decoding transactions")
		return transactions, err
	}

	return transactions, nil
}


func (transaction *Transaction) GetTransactionById(accountId string) error {
	log.Printf("Get transaction by id")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(transaction.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Query to get transaction by id
	query := bson.D{
		{Key: "_id", Value: id},
	}

	if accountId != "" {
		// String id to ObjectId
		accountIdObject, err := primitive.ObjectIDFromHex(accountId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return err
		}

		// Change query to only search the account transactions
		query = bson.D{
			{Key: "accountId", Value: accountIdObject},
			{Key: "_id", Value: id},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("transactions", query, filter)
	if err != nil {
		log.Printf("Error: Get transaction in db")
		return err
	}

	// Decode query
	cursor.Next(mongo.GetCtx())

	if err = cursor.Decode(&transaction); err != nil {
		log.Printf("Error: Decoding transaction")
		return err
	}

	return nil
}

func (transaction *Transaction) Update(accountId string, description *string, date *time.Time, amount *float64, category *string) error {
	log.Printf("Update transaction")

	// Get current state of transaction
	transaction.GetTransactionById(accountId)

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(transaction.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	accountID, err := primitive.ObjectIDFromHex(transaction.AccountID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Set fields if provided
	if description != nil {
		transaction.Description = *description
	}
	if date != nil {
		transaction.Date = *date
	}
	if amount != nil {
		transaction.Amount = *amount
	}
	if category != nil {
		transaction.Category = *category
	}

	// Filter to get user by id
	filter := bson.D{
		{Key: "_id", Value: id},
	}

	// Update fields
	update := bson.M{
		"description": transaction.Description,
		"date": transaction.Date,
		"amount": transaction.Amount,
		"category": transaction.Category,
		"accountId": accountID,
	}

	// Execute update
	result, err := mongo.ReplaceOne("transactions", filter, update)
	if err != nil {
		log.Printf("Error: Update transaction in db")
		return err
	}
	if result.ModifiedCount == 0 {
		log.Printf("Error: Update transaction in db")
		return fmt.Errorf("no changes")
	}

	return nil
}

func Delete(transactionId string, accountId string) (bool, error) {
	log.Printf("Delete transaction by id")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(transactionId)
	if err != nil {
		log.Printf("Error: Convert string to id")
		return false, err
	}

	// Query to delete transaction by id
	query := bson.D{
		{Key: "_id", Value: id},
	}

	if accountId != "" {
		// String id to ObjectId
		accountIdObject, err := primitive.ObjectIDFromHex(accountId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return false, err
		}

		// Change query to only search the account transactions
		query = bson.D{
			{Key: "accountId", Value: accountIdObject},
			{Key: "_id", Value: id},
		}
	}

	// Execute query
	result, err := mongo.DeleteOne("transactions", query)
	if err != nil {
		log.Printf("Error: Delete transaction in db")
		return false, err
	}

	if result.DeletedCount == 0 {
		log.Printf("Error: Delete transaction in db")
		return false, fmt.Errorf("no changes")
	}

	return true, nil
}