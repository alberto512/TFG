package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.25

import (
	"context"
	"fmt"
	"log"
	"tfg/graph/model"
	"tfg/internal/jwt"
	"tfg/internal/middleware"
	"tfg/internal/operations"
	"tfg/internal/users"
	"time"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (string, error) {
	// Set user
	var user users.User

	log.Printf("Route: Login")

	user.Username = username
	user.Password = password

	// Authenticate user
	correct, err := user.Authenticate()
	if err != nil || !correct {
		log.Printf("Error: Authenticating user")
		return "", err
	}

	// Generate jwt token to user
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		log.Printf("Error: Generating token")
		return "", err
	}

	return token, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (string, error) {
	log.Printf("Route: Refresh")

	// Parse actual token
	username, err := jwt.ParseToken(token)
	if err != nil {
		log.Printf("Error: Parsing token")
		return "", fmt.Errorf("access denied")
	}

	// Generate new token
	newToken, err := jwt.GenerateToken(username)
	if err != nil {
		log.Printf("Error: Generating token")
		return "", err
	}

	return newToken, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, username string, password string, rol model.Rol) (*model.User, error) {
	// Set user
	var user users.User

	log.Printf("Route: CreateUser")

	// Create user struct
	user.Username = username
	user.Password = password
	user.Rol = rol

	// Create user in db
	err := user.Create()
	if err != nil {
		log.Printf("Error: Create user")
		return &model.User{}, err
	}

	// Parse to model.User
	graphqlUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Rol:      user.Rol,
	}

	return graphqlUser, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	// Set variables
	var graphqlUsers []*model.User
	var userAuth *users.User

	log.Printf("Route: Users")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlUsers, fmt.Errorf("access denied")
	}

	// Only allow operation if is an admin
	if userAuth.Rol != model.RolAdmin {
		log.Printf("Error: Access denied")
		return graphqlUsers, fmt.Errorf("access denied")
	}

	// Get all users
	dbUsers, err := users.GetAllUsers()
	if err != nil {
		log.Printf("Error: Get all users")
		return graphqlUsers, err
	}

	// Parse al users to model.User
	for _, user := range dbUsers {
		graphqlUsers = append(graphqlUsers, &model.User{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
			Rol:      user.Rol,
		})
	}

	return graphqlUsers, nil
}

// UserByID is the resolver for the userById field.
func (r *queryResolver) UserByID(ctx context.Context, id string) (*model.User, error) {
	// Set variables
	var user users.User
	var userAuth *users.User

	log.Printf("Route: UserByID")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.User{}, fmt.Errorf("access denied")
	}

	// Only allow operation if is an admin
	if userAuth.Rol != model.RolAdmin {
		log.Printf("Error: Access denied")
		return &model.User{}, fmt.Errorf("access denied")
	}

	user.ID = id

	// Get user
	err := user.GetUserById()
	if err != nil {
		log.Printf("Error: Get user")
		return &model.User{}, err
	}

	// Parse to model.User
	graphqlUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Rol:      user.Rol,
	}

	return graphqlUser, nil
}

// UserByToken is the resolver for the userByToken field.
func (r *queryResolver) UserByToken(ctx context.Context) (*model.User, error) {
	// Set variables
	var user users.User
	var userAuth *users.User

	log.Printf("Route: UserByToken")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.User{}, fmt.Errorf("access denied")
	}

	user.ID = userAuth.ID

	// Get user
	err := user.GetUserById()
	if err != nil {
		log.Printf("Error: Get user")
		return &model.User{}, err
	}

	// Parse to model.User
	graphqlUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Rol:      user.Rol,
	}

	return graphqlUser, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, password string) (*model.User, error) {
	// Set variables
	var user users.User
	var userAuth *users.User

	log.Printf("Route: UpdateUser")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.User{}, fmt.Errorf("access denied")
	}

	// Set fields. Admin can edit any user and User can only edit himself
	if userAuth.Rol == model.RolAdmin {
		user.ID = id
	} else {
		user.ID = userAuth.ID
	}
	user.Password = password

	// Update user
	err := user.Update()
	if err != nil {
		log.Printf("Error: Create user")
		return &model.User{}, err
	}

	// Parse to model.User
	graphqlUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Rol:      user.Rol,
	}

	return graphqlUser, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (string, error) {
	// Set variables
	var userAuth *users.User

	log.Printf("Route: DeleteUser")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return "Error", fmt.Errorf("access denied")
	}

	// Only allow deletion if is an admin
	if userAuth.Rol != model.RolAdmin {
		log.Printf("Error: Access denied")
		return "Error", fmt.Errorf("access denied")
	}

	// Delete user
	correct, err := users.Delete(id)
	if err != nil || !correct {
		log.Printf("Error: Delete user")
		return "Error", err
	}

	return "Deletion success", nil
}

// CreateOperation is the resolver for the createOperation field.
func (r *mutationResolver) CreateOperation(ctx context.Context, input model.NewOperation) (*model.Operation, error) {
	// Set variables
	var userAuth *users.User
	var operation operations.Operation

	log.Printf("Route: CreateOperation")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.Operation{}, fmt.Errorf("access denied")
	}

	// Parse date
	date := time.UnixMilli(int64(input.Date))

	// Create operation struct
	operation.Description = input.Description
	operation.Date = date
	operation.Amount = input.Amount
	operation.Category = input.Category
	operation.UserID = userAuth.ID

	// Create operation in db
	err := operation.Create()
	if err != nil {
		log.Printf("Error: Create operation")
		return &model.Operation{}, err
	}

	// Parse to model.Operation
	graphqlOperation := &model.Operation{
		ID:          operation.ID,
		Description: operation.Description,
		Date:        int(operation.Date.UnixMilli()),
		Amount:      operation.Amount,
		Category:    operation.Category,
		UserID:      operation.UserID,
	}

	return graphqlOperation, nil
}

// Operations is the resolver for the operations field.
func (r *queryResolver) Operations(ctx context.Context) ([]*model.Operation, error) {
	// Set variable
	var userAuth *users.User
	var graphqlOperations []*model.Operation
	var userId string

	log.Printf("Route: Operations")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlOperations, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any operation and user can only get his own operations
	if userAuth.Rol == model.RolAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Get all operations of user
	operations, err := operations.GetAllOperations(userId)
	if err != nil {
		log.Printf("Error: Get all operations")
		return graphqlOperations, err
	}

	// Parse operations of user
	for _, operation := range operations {
		graphqlOperation := &model.Operation{
			ID:          operation.ID,
			Description: operation.Description,
			Date:        int(operation.Date.UnixMilli()),
			Amount:      operation.Amount,
			Category:    operation.Category,
			UserID:      operation.UserID,
		}
		graphqlOperations = append(graphqlOperations, graphqlOperation)
	}

	return graphqlOperations, nil
}

// OperationsByDate is the resolver for the operationsByDate field.
func (r *queryResolver) OperationsByDate(ctx context.Context, initDate int, endDate int) ([]*model.Operation, error) {
	// Set variable
	var userAuth *users.User
	var graphqlOperations []*model.Operation
	var userId string

	log.Printf("Route: OperationsByDate")

	// Parse dates
	date1 := time.UnixMilli(int64(initDate))
	date2 := time.UnixMilli(int64(endDate))

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlOperations, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any operation and user can only get his own operations
	if userAuth.Rol == model.RolAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Get all operations of user by date
	operations, err := operations.GetOperationsByDate(userId, date1, date2)
	if err != nil {
		log.Printf("Error: Get all operations by date")
		return graphqlOperations, err
	}

	// Parse operations of user
	for _, operation := range operations {
		graphqlOperation := &model.Operation{
			ID:          operation.ID,
			Description: operation.Description,
			Date:        int(operation.Date.UnixMilli()),
			Amount:      operation.Amount,
			Category:    operation.Category,
			UserID:      operation.UserID,
		}
		graphqlOperations = append(graphqlOperations, graphqlOperation)
	}

	return graphqlOperations, nil
}

// OperationsByCategory is the resolver for the operationsByCategory field.
func (r *queryResolver) OperationsByCategory(ctx context.Context, category string) ([]*model.Operation, error) {
	// Set variable
	var userAuth *users.User
	var graphqlOperations []*model.Operation
	var userId string

	log.Printf("Route: OperationsByCategory")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlOperations, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any operation and user can only get his own operations
	if userAuth.Rol == model.RolAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Get all operations of user by category
	operations, err := operations.GetOperationsByCategory(userId, category)
	if err != nil {
		log.Printf("Error: Get all operations by category")
		return graphqlOperations, err
	}

	// Parse operations of user
	for _, operation := range operations {
		graphqlOperation := &model.Operation{
			ID:          operation.ID,
			Description: operation.Description,
			Date:        int(operation.Date.UnixMilli()),
			Amount:      operation.Amount,
			Category:    operation.Category,
			UserID:      operation.UserID,
		}
		graphqlOperations = append(graphqlOperations, graphqlOperation)
	}

	return graphqlOperations, nil
}

// OperationByID is the resolver for the operationById field.
func (r *queryResolver) OperationByID(ctx context.Context, id string) (*model.Operation, error) {
	// Set variables
	var userAuth *users.User
	var operation operations.Operation
	var userId string

	log.Printf("Route: OperationByID")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.Operation{}, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any operation and user can only get his own operations
	if userAuth.Rol == model.RolAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}
	operation.ID = id

	// Get operation in db
	err := operation.GetOperationById(userId)
	if err != nil {
		log.Printf("Error: Get operation")
		return &model.Operation{}, err
	}

	// Parse to model.Operation
	graphqlOperation := &model.Operation{
		ID:          operation.ID,
		Description: operation.Description,
		Date:        int(operation.Date.UnixMilli()),
		Amount:      operation.Amount,
		Category:    operation.Category,
		UserID:      operation.UserID,
	}

	return graphqlOperation, nil
}

// UpdateOperation is the resolver for the updateOperation field.
func (r *mutationResolver) UpdateOperation(ctx context.Context, input *model.UpdateOperation) (*model.Operation, error) {
	// Set variables
	var userAuth *users.User
	var operation operations.Operation
	var userId string

	log.Printf("Route: UpdateOperation")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.Operation{}, fmt.Errorf("access denied")
	}

	// Set fields. Admin can update any operation and user can only update his own operations
	if userAuth.Rol == model.RolAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}
	operation.ID = input.ID

	// Parse dates
	var date *time.Time
	if input.Date != nil {
		dateLocal := time.UnixMilli(int64(*input.Date))
		date = &dateLocal
	} else {
		date = nil
	}

	// Update operation in db
	err := operation.Update(userId, input.Description, date, input.Amount, input.Category)
	if err != nil {
		log.Printf("Error: Update operation")
		return &model.Operation{}, err
	}

	// Parse to model.Operation
	graphqlOperation := &model.Operation{
		ID:          operation.ID,
		Description: operation.Description,
		Date:        int(operation.Date.UnixMilli()),
		Amount:      operation.Amount,
		Category:    operation.Category,
		UserID:      operation.UserID,
	}

	return graphqlOperation, nil
}

// DeleteOperation is the resolver for the deleteOperation field.
func (r *mutationResolver) DeleteOperation(ctx context.Context, id string) (string, error) {
	// Set variables
	var userAuth *users.User
	var userId string

	log.Printf("Route: DeleteOperation")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return "Error", fmt.Errorf("access denied")
	}

	// Set fields. Admin can delete any operation and user can only delete his own operations
	if userAuth.Rol == model.RolAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Delete operation in db
	correct, err := operations.Delete(id, userId)
	if err != nil || !correct {
		log.Printf("Error: Delete operation")
		return "Error", err
	}

	return "Deletion success", nil
}

// Operations is the resolver for the operations field.
func (r *userResolver) Operations(ctx context.Context, obj *model.User) ([]*model.Operation, error) {
	// Set variable
	var graphqlOperations []*model.Operation

	log.Printf("Resolver: Operations")

	// Get operations of user
	operations, err := operations.GetAllOperations(obj.ID)
	if err != nil {
		log.Printf("Error: Get operations by id")
		return graphqlOperations, err
	}

	// Parse operations of user
	for _, operation := range operations {
		graphqlOperation := &model.Operation{
			ID:          operation.ID,
			Description: operation.Description,
			Date:        int(operation.Date.UnixMilli()),
			Amount:      operation.Amount,
			Category:    operation.Category,
			UserID:      operation.UserID,
		}
		graphqlOperations = append(graphqlOperations, graphqlOperation)
	}

	return graphqlOperations, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
