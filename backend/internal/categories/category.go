package categories

import (
	"fmt"
	"log"
	"tfg/internal/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID          string  	`bson:"_id"`
	Name		string		`bson:"name"`
	UserID		string		`bson:"userId"`
}

func (category *Category) Create() (error) {
	log.Printf("Create category")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(category.UserID)
	if err != nil {
		log.Printf("Error: Convert string to id")
		return err
	}

	// Execute insert
	res, err := mongo.InsertOne("categories", bson.D{
		{Key: "name", Value: category.Name},
		{Key: "userId", Value: id},
	})
	if err != nil {
		log.Printf("Error: Create category in db")
		return err
	}

	// Save ID
	category.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func GetAllCategories(userId string) ([]Category, error) {
	var categories []Category

	log.Printf("Get all categories")

	// Query to get all categories
	query := bson.D{}

	if userId != "" {
		// String id to ObjectId
		id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Printf("Error: Convert string to id")
			return categories, err
		}

		// Change query to only search the user categories
		query = bson.D{
			{Key: "userId", Value: id},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("categories", query, filter)
	if err != nil {
		log.Printf("Error: Get all categories in db")
		return categories, err
	}

	// Decode query
	if err = cursor.All(mongo.GetCtx(), &categories); err != nil {
		log.Printf("Error: Decoding categories")
		return categories, err
	}

	return categories, nil
}

func (category *Category) GetCategoryById(userId string) error {
	log.Printf("Get category by id")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(category.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
		return err
	}

	// Query to get category by id
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

		// Change query to only search the user categories
		query = bson.D{
			{Key: "userId", Value: userIdObject},
			{Key: "_id", Value: id},
		}
	}

	// Empty filter
	filter := bson.D{}

	// Execute query
	cursor, err := mongo.Query("categories", query, filter)
	if err != nil {
		log.Printf("Error: Get category in db")
		return err
	}

	// Decode query
	cursor.Next(mongo.GetCtx())

	if err = cursor.Decode(&category); err != nil {
		log.Printf("Error: Decoding category")
		return err
	}

	return nil
}

func Delete(categoryId string, userId string) (bool, error) {
	log.Printf("Delete category by id")

	// String id to ObjectId
	id, err := primitive.ObjectIDFromHex(categoryId)
	if err != nil {
		log.Printf("Error: Convert string to id")
		return false, err
	}

	// Query to delete category by id
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

		// Change query to only search the user categories
		query = bson.D{
			{Key: "userId", Value: userIdObject},
			{Key: "_id", Value: id},
		}
	}

	// Execute query
	result, err := mongo.DeleteOne("categories", query)
	if err != nil {
		log.Printf("Error: Delete category in db")
		return false, err
	}

	if result.DeletedCount == 0 {
		log.Printf("Error: Delete category in db")
		return false, fmt.Errorf("no changes")
	}

	return true, nil
}