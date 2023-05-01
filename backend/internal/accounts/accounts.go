package accounts

import (
	"fmt"
	"log"
	"tfg/internal/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
    ID			string		`bson:"_id"`
    Iban		string		`bson:"iban"`
    Name		string		`bson:"name"`
    Currency	string		`bson:"currency"`
	Amount		float64		`bson:"amount"`
	Bank		string		`bson:"bank"`
	UpdateDate  time.Time	`bson:"updateDate"`
	UserID		string		`bson:"userId"`
}

func (account *Account) Create() (error) {
	log.Printf("Create account")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(account.UserID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Execute insert
	res, err := mongo.InsertOne("accounts", bson.D{
		{Key: "iban", Value: account.Iban},
		{Key: "name", Value: account.Name},
		{Key: "currency", Value: account.Currency},
		{Key: "amount", Value: account.Amount},
		{Key: "bank", Value: account.Bank},
		{Key: "updateDate", Value: time.Now()},
		{Key: "userId", Value: id},
	})
	if err != nil {
		log.Printf("Error: Create account in db")
		return err
	}

	// Save ID
	account.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func GetAllAccounts(userId string) ([]Account, error) {
	var accounts []Account

	log.Printf("Get all accounts")

	// Query to get all accounts
	query := bson.D{}

	if userId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return accounts, err
		}

		// Change query to only search the user accounts
		query = bson.D{
			{Key: "userId", Value: id},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("accounts", query, filter)
	if err != nil {
		log.Printf("Error: Get all accounts in db")
		return accounts, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &accounts); err != nil {
		log.Printf("Error: Decoding accounts")
		return accounts, err
	}

	return accounts, nil
}

func (account *Account) GetAccountById(userId string) error {
	log.Printf("Get account by id")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(account.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Query to get account by id
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

		// Change query to only search the user accounts
		query = bson.D{
			{Key: "userId", Value: userIdObject},
			{Key: "_id", Value: id},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("accounts", query, filter)
	if err != nil {
		log.Printf("Error: Get account in db")
		return err
	}

	// Decode query
	cursor.Next(mongo.GetCtx())

	if err = cursor.Decode(&account); err != nil {
		log.Printf("Error: Decoding account")
		return err
	}

	return nil
}

func (account *Account) GetAccountByIban(userId string) error {
	log.Printf("Get account by iban")

	// Query to get account by iban
	query := bson.D{
		{Key: "iban", Value: account.Iban},
	}

	if userId != "" {
		// String id to ObjectId
		userIdObject, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return err
		}

		// Change query to only search the user accounts
		query = bson.D{
			{Key: "userId", Value: userIdObject},
			{Key: "iban", Value: account.Iban},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("accounts", query, filter)
	if err != nil {
		log.Printf("Error: Get account in db")
		return err
	}

	// Decode query
	cursor.Next(mongo.GetCtx())

	println(account)

	if err = cursor.Decode(&account); err != nil {
		log.Printf("Error: Decoding account")
		return err
	}

	println(account)

	return nil
}

func (account *Account) Update(userId string, iban *string, name *string, currency *string, amount *float64, bank *string) error {
	log.Printf("Update account")

	// Get current state of account
	account.GetAccountById(userId)

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(account.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	userID, err := primitive.ObjectIDFromHex(account.UserID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

	// Set fields if provided
	if iban != nil {
		account.Iban = *iban
	}
	if name != nil {
		account.Name = *name
	}
	if currency != nil {
		account.Currency = *currency
	}
	if amount != nil {
		account.Amount = *amount
	}
	if bank != nil {
		account.Bank = *bank
	}

	// Filter to get user by id
	filter := bson.D{
		{Key: "_id", Value: id},
	}

	// Update fields
	update := bson.M{
		"iban": account.Iban,
		"name": account.Name,
		"currency": account.Currency,
		"amount": account.Amount,
		"bank": account.Bank,
		"updateDate": time.Now(),
		"userId": userID,
	}

	// Execute update
	result, err := mongo.ReplaceOne("accounts", filter, update)
	if err != nil {
		log.Printf("Error: Update account in db")
		return err
	}
	if result.ModifiedCount == 0 {
		log.Printf("Error: Update account in db")
		return fmt.Errorf("no changes")
	}

	return nil
}

func Delete(accountId string, userId string) (bool, error) {
	log.Printf("Delete account by id")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(accountId)
	if err != nil {
		log.Printf("Error: Convert string to id")
		return false, err
	}

	// Query to delete account by id
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

		// Change query to only search the user accounts
		query = bson.D{
			{Key: "userId", Value: userIdObject},
			{Key: "_id", Value: id},
		}
	}

	// Execute query
	result, err := mongo.DeleteOne("accounts", query)
	if err != nil {
		log.Printf("Error: Delete account in db")
		return false, err
	}

	if result.DeletedCount == 0 {
		log.Printf("Error: Delete account in db")
		return false, fmt.Errorf("no changes")
	}

	return true, nil
}