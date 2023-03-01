package users

import (
	"fmt"
	"log"
	"tfg/graph/model"
	"tfg/internal/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       string `bson:"_id,omitempty"`
    Username string `bson:"username,omitempty"`
    Password string `bson:"password,omitempty"`
    Rol 	 model.Rol `bson:"rol,omitempty"`
}

func (user *User) Authenticate() (bool, error) {
    var mongoUser User

    log.Printf("Authenticate user")

    // Query to get user by username
    query := bson.D{
        {Key: "username", Value: user.Username},
    }

    // Filter query to only obtain password
    filter := bson.D{
        {Key: "password", Value: 1},
        {Key: "_id", Value: 0},
    }

    // Execute query
    cursor, err := mongo.Query("users", query, filter)
    if err != nil {
        log.Printf("Error: Get user in db")
        return false, err
    }

    // Decode query
    cursor.Next(mongo.GetCtx())

    if err = cursor.Decode(&mongoUser); err != nil {
        log.Printf("Error: Decoding user")
        return false, err
    }

    // Check if passwords are equal
    return CheckPasswordHash(user.Password, mongoUser.Password), nil
}

func CheckPasswordHash(password, hash string) (bool) {
    log.Printf("Check password")
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (user *User) Create() (error) {
    log.Printf("Create user")

    // Hash password
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
        log.Printf("Error: Hash password")
        return err
    }

    // Execute insert
    res, err := mongo.InsertOne("users", bson.D{
        {Key: "username", Value: user.Username},
        {Key: "password", Value: hashedPassword},
        {Key: "rol", Value: user.Rol},
    })
    if err != nil {
        log.Printf("Error: Create user in db")
        return err
    }

    // Save ID
    user.ID = res.InsertedID.(primitive.ObjectID).Hex()
    return nil
}

func HashPassword(password string) (string, error) {
    log.Printf("Hash password")
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func GetAllUsers() ([]User, error) {
    var users []User

    log.Printf("Get all users")

    // Query to get all users
    query := bson.D{}

    // Filter to exclude operations
    filter := bson.D{
        {Key: "operations", Value: 0},
    }

    // Execute query
    cursor, err := mongo.Query("users", query, filter)
    if err != nil {
        log.Printf("Error: Get users in db")
        return users, err
    }

    // Decode query
    if err = cursor.All(mongo.GetCtx(), &users); err != nil {
        log.Printf("Error: Decoding users")
        return users, err
    }

    return users, nil
}

func (user *User) GetUserById() (error) {
    log.Printf("Get user by id")

	// String id to ObjectId
	id, err :=  primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

    // Query to get user by id
    query := bson.D{
        {Key: "_id", Value: id},
    }

    // Filter to exclude operations
    filter := bson.D{
        {Key: "operations", Value: 0},
    }

    // Execute query
    cursor, err := mongo.Query("users", query, filter)
    if err != nil {
        log.Printf("Error: Get user in db")
        return err
    }

    // Decode query
    cursor.Next(mongo.GetCtx())

    if err = cursor.Decode(&user); err != nil {
        log.Printf("Error: Decoding user")
        return err
    }

    return nil
}

func (user *User) GetUserByUsername() (error) {
    log.Printf("Get user by username")

    // Query to get user by username
    query := bson.D{
        {Key: "username", Value: user.Username},
    }

    // Filter to exclude operations
    filter := bson.D{
        {Key: "operations", Value: 0},
    }

    // Execute query
    cursor, err := mongo.Query("users", query, filter)
    if err != nil {
        log.Printf("Error: Get user in db")
        return err
    }

    // Decode query
    cursor.Next(mongo.GetCtx())

    if err = cursor.Decode(&user); err != nil {
        log.Printf("Error: Decoding user")
        return err
    }

    return nil
}

func (user *User) Update() (error) {
    log.Printf("Update user")

	// String id to ObjectId
	id, err :=  primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return err
	}

    // Filter to get user by id
    filter := bson.D{
        {Key: "_id", Value: id},
    }

	// Hash password
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
        log.Printf("Error: Hash password")
        return err
    }

    // Update field
    update := bson.D{
		{Key: "$set", Value: bson.D{
				{Key: "password", Value: hashedPassword},
			},
		},
	}

    // Execute update
    result, err := mongo.UpdateOne("users", filter, update)
    if err != nil {
        log.Printf("Error: Update user in db")
        return err
    }
	if result.ModifiedCount == 0 {
		log.Printf("Error: Update user in db")
        return fmt.Errorf("no changes")
	}

	user.GetUserById()

    return nil
}

func Delete(userId string) (bool, error) {
    log.Printf("Delete user by id")

	// String id to ObjectId
	id, err :=  primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Printf("Error: Convert string to id")
        return false, err
	}

    // Query to delete user by id
    query := bson.D{
        {Key: "_id", Value: id},
    }

    // Execute query
    result, err := mongo.DeleteOne("users", query)
    if err != nil {
        log.Printf("Error: Delete user in db")
        return false, err
    }

    if result.DeletedCount == 0 {
		log.Printf("Error: Delete user in db")
        return false, fmt.Errorf("no changes")
	}

    return true, nil
}