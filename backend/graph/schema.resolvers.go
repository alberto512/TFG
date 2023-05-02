package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.25

import (
	"context"
	"fmt"
	"log"
	"tfg/graph/model"
	"tfg/internal/accounts"
	"tfg/internal/categories"
	"tfg/internal/jwt"
	"tfg/internal/middleware"
	"tfg/internal/santander"
	"tfg/internal/transactions"
	"tfg/internal/users"
	"time"
)

// User is the resolver for the user field.
func (r *accountResolver) User(ctx context.Context, obj *model.Account) (*model.User, error) {
	// Set variables
	var userAuth *users.User
	var account accounts.Account
	var user users.User

	log.Printf("Resolver: User-Account")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.User{}, fmt.Errorf("access denied")
	}

	account.ID = obj.ID

	// Get account in db
	err := account.GetAccountById("")
	if err != nil {
		log.Printf("Error: Get account")
		return &model.User{}, err
	}

	user.ID = account.UserID

	// Get user
	err = user.GetUserById()
	if err != nil {
		log.Printf("Error: Get user")
		return &model.User{}, err
	}

	// Parse to model.User
	graphqlUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}

	return graphqlUser, nil
}

// Transactions is the resolver for the transactions field.
func (r *accountResolver) Transactions(ctx context.Context, obj *model.Account) ([]*model.Transaction, error) {
	// Set variables
	var userAuth *users.User
	var graphqlTransactions []*model.Transaction

	log.Printf("Resolver: Transactions-Account")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlTransactions, fmt.Errorf("access denied")
	}

	// Get all transactions of account
	transactions, err := transactions.GetAllTransactions(obj.ID)
	if err != nil {
		log.Printf("Error: Get all transactions")
		return graphqlTransactions, err
	}

	// Parse transactions of account
	for _, transaction := range transactions {
		graphqlTransaction := &model.Transaction{
			ID:          transaction.ID,
			Description: transaction.Description,
			Date:        int(transaction.Date.UnixMilli()),
			Amount:      transaction.Amount,
		}
		graphqlTransactions = append(graphqlTransactions, graphqlTransaction)
	}

	return graphqlTransactions, nil
}

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
func (r *mutationResolver) CreateUser(ctx context.Context, username string, password string, role model.Role) (*model.User, error) {
	// Set user
	var user users.User
	var category categories.Category

	log.Printf("Route: CreateUser")

	// Create user struct
	user.Username = username
	user.Password = password
	user.Role = role

	// Create user in db
	err := user.Create()
	if err != nil {
		log.Printf("Error: Create user")
		return &model.User{}, err
	}

	// Create default categories for user
	category.UserID = user.ID

	// Food category
	category.Name = "Food"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Home category
	category.Name = "Home"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Lifestyle category
	category.Name = "Lifestyle"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Health category
	category.Name = "Health"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Shopping category
	category.Name = "Shopping"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Children category
	category.Name = "Children"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Vacation category
	category.Name = "Vacation"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Education category
	category.Name = "Education"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Salary category
	category.Name = "Salary"
	err = category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.User{}, err
	}

	// Parse to model.User
	graphqlUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}

	return graphqlUser, nil
}

// UpdatePassword is the resolver for the updatePassword field.
func (r *mutationResolver) UpdatePassword(ctx context.Context, id string, password string) (*model.User, error) {
	// Set variables
	var user users.User
	var userAuth *users.User

	log.Printf("Route: UpdatePassword")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.User{}, fmt.Errorf("access denied")
	}

	// Set fields. Admin can edit any user and User can only edit himself
	if userAuth.Role == model.RoleAdmin {
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
		Role:     user.Role,
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
	if userAuth.Role != model.RoleAdmin {
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

// CreateCategory is the resolver for the createCategory field.
func (r *mutationResolver) CreateCategory(ctx context.Context, name string) (*model.Category, error) {
	// Set variables
	var category categories.Category
	var userAuth *users.User

	log.Printf("Route: CreateCategory")

	// Create user struct
	category.Name = name

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.Category{}, fmt.Errorf("access denied")
	}

	// Set category values
	category.Name = name
	category.UserID = userAuth.ID

	// Create category in db
	err := category.Create()
	if err != nil {
		log.Printf("Error: Create category")
		return &model.Category{}, err
	}

	// Parse to model.Category
	graphqlCategory := &model.Category{
		ID:   category.ID,
		Name: category.Name,
	}

	return graphqlCategory, nil
}

// DeleteCategory is the resolver for the deleteCategory field.
func (r *mutationResolver) DeleteCategory(ctx context.Context, id string) (string, error) {
	// Set variables
	var userAuth *users.User
	var userId string

	log.Printf("Route: DeleteCategory")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return "Error", fmt.Errorf("access denied")
	}

	// Set fields. Admin can delete any category and user can only delete his own categories
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Delete user
	correct, err := categories.Delete(id, userId)
	if err != nil || !correct {
		log.Printf("Error: Delete category")
		return "Error", err
	}

	return "Deletion success", nil
}

// RefreshBankData is the resolver for the refreshBankData field.
func (r *mutationResolver) RefreshBankData(ctx context.Context) (string, error) {
	// Set variables
	var userAuth *users.User

	log.Printf("Route: RefreshBankData")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return "Error", fmt.Errorf("access denied")
	}

	err := santander.RefreshData(userAuth.ID)
	if err != nil {
		log.Printf("Error: Refresh data")
		return "Error", err
	}

	return "Refresh success", nil
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
	if userAuth.Role != model.RoleAdmin {
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
			Role:     user.Role,
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
	if userAuth.Role != model.RoleAdmin {
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
		Role:     user.Role,
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
		Role:     user.Role,
	}

	return graphqlUser, nil
}

// Accounts is the resolver for the accounts field.
func (r *queryResolver) Accounts(ctx context.Context) ([]*model.Account, error) {
	// Set variables
	var userAuth *users.User
	var graphqlAccounts []*model.Account
	var userId string

	log.Printf("Route: Accounts")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlAccounts, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any account and user can only get his own accounts
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Get all accounts of user
	accounts, err := accounts.GetAllAccounts(userId)
	if err != nil {
		log.Printf("Error: Get all accounts")
		return graphqlAccounts, err
	}

	// Parse accounts of user
	for _, account := range accounts {
		graphqlAccount := &model.Account{
			ID:       account.ID,
			Iban:     account.Iban,
			Name:     account.Name,
			Currency: account.Currency,
			Amount:   account.Amount,
			Bank:     account.Bank,
		}
		graphqlAccounts = append(graphqlAccounts, graphqlAccount)
	}

	return graphqlAccounts, nil
}

// AccountByID is the resolver for the accountById field.
func (r *queryResolver) AccountByID(ctx context.Context, id string) (*model.Account, error) {
	// Set variables
	var userAuth *users.User
	var account accounts.Account
	var userId string

	log.Printf("Route: AccountByID")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.Account{}, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any account and user can only get his own account
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}
	account.ID = id

	// Get account in db
	err := account.GetAccountById(userId)
	if err != nil {
		log.Printf("Error: Get account")
		return &model.Account{}, err
	}

	// Parse to model.Account
	graphqlAccount := &model.Account{
		ID:       account.ID,
		Iban:     account.Iban,
		Name:     account.Name,
		Currency: account.Currency,
		Amount:   account.Amount,
		Bank:     account.Bank,
	}

	return graphqlAccount, nil
}

// Transactions is the resolver for the transactions field.
func (r *queryResolver) Transactions(ctx context.Context) ([]*model.Transaction, error) {
	// Set variables
	var userAuth *users.User
	var graphqlTransactions []*model.Transaction
	var userId string

	log.Printf("Route: Transactions")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlTransactions, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any transaction and user can only get his own transactions
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Get all accounts of user
	accounts, err := accounts.GetAllAccounts(userId)
	if err != nil {
		log.Printf("Error: Get all accounts")
		return graphqlTransactions, err
	}

	// Iterate for every account
	for _, account := range accounts {
		// Get all transactions of account
		transactions, err := transactions.GetAllTransactions(account.ID)
		if err != nil {
			log.Printf("Error: Get all transactions")
			return graphqlTransactions, err
		}

		// Parse transactions of account
		for _, transaction := range transactions {
			graphqlTransaction := &model.Transaction{
				ID:          transaction.ID,
				Description: transaction.Description,
				Date:        int(transaction.Date.UnixMilli()),
				Amount:      transaction.Amount,
			}
			graphqlTransactions = append(graphqlTransactions, graphqlTransaction)
		}
	}

	return graphqlTransactions, nil
}

// TransactionByID is the resolver for the transactionById field.
func (r *queryResolver) TransactionByID(ctx context.Context, id string) (*model.Transaction, error) {
	// Set variables
	var userAuth *users.User
	var transaction transactions.Transaction
	var userId string

	log.Printf("Route: TransactionByID")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.Transaction{}, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any transaction and user can only get his own transactions
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	transaction.ID = id

	// Get all accounts of user
	accounts, err := accounts.GetAllAccounts(userId)
	if err != nil {
		log.Printf("Error: Get all accounts")
		return &model.Transaction{}, err
	}

	// Iterate for every account
	for _, account := range accounts {
		// Get transaction by id
		err := transaction.GetTransactionById(account.ID)
		if err != nil || transaction.Description == "" {
			log.Printf("Error: Get account")
			continue
		}

		break
	}

	// Check if transaction exists in one of the accounts
	if transaction.Description == "" {
		log.Printf("Error: Get account")
		return &model.Transaction{}, fmt.Errorf("not found")
	}

	// Parse to model.Transaction
	graphqlTransaction := &model.Transaction{
		ID:          transaction.ID,
		Description: transaction.Description,
		Date:        int(transaction.Date.UnixMilli()),
		Amount:      transaction.Amount,
	}

	return graphqlTransaction, nil
}

// TransactionsByDate is the resolver for the transactionsByDate field.
func (r *queryResolver) TransactionsByDate(ctx context.Context, initDate int, endDate int) ([]*model.Transaction, error) {
	// Set variables
	var userAuth *users.User
	var graphqlTransactions []*model.Transaction
	var userId string

	log.Printf("Route: TransactionsByDate")

	// Parse dates
	date1 := time.UnixMilli(int64(initDate))
	date2 := time.UnixMilli(int64(endDate))

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlTransactions, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any transaction and user can only get his own transactions
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Get all accounts of user
	accounts, err := accounts.GetAllAccounts(userId)
	if err != nil {
		log.Printf("Error: Get all accounts")
		return graphqlTransactions, err
	}

	// Iterate for every account
	for _, account := range accounts {
		// Get all transactions of account by date
		transactions, err := transactions.GetTransactionsByDate(account.ID, date1, date2)
		if err != nil {
			log.Printf("Error: Get all transactions by date")
			return graphqlTransactions, err
		}

		// Parse transactions of account
		for _, transaction := range transactions {
			graphqlTransaction := &model.Transaction{
				ID:          transaction.ID,
				Description: transaction.Description,
				Date:        int(transaction.Date.UnixMilli()),
				Amount:      transaction.Amount,
			}
			graphqlTransactions = append(graphqlTransactions, graphqlTransaction)
		}
	}

	return graphqlTransactions, nil
}

// TransactionsByCategory is the resolver for the transactionsByCategory field.
func (r *queryResolver) TransactionsByCategory(ctx context.Context, category string) ([]*model.Transaction, error) {
	// Set variables
	var userAuth *users.User
	var graphqlTransactions []*model.Transaction
	var userId string

	log.Printf("Route: TransactionsByCategory")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlTransactions, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any transaction and user can only get his own transactions
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Get all accounts of user
	accounts, err := accounts.GetAllAccounts(userId)
	if err != nil {
		log.Printf("Error: Get all accounts")
		return graphqlTransactions, err
	}

	// Iterate for every account
	for _, account := range accounts {
		// Get all transactions of account by category
		transactions, err := transactions.GetTransactionsByCategory(account.ID, category)
		if err != nil {
			log.Printf("Error: Get all transactions by category")
			return graphqlTransactions, err
		}

		// Parse transactions of account
		for _, transaction := range transactions {
			graphqlTransaction := &model.Transaction{
				ID:          transaction.ID,
				Description: transaction.Description,
				Date:        int(transaction.Date.UnixMilli()),
				Amount:      transaction.Amount,
			}
			graphqlTransactions = append(graphqlTransactions, graphqlTransaction)
		}
	}

	return graphqlTransactions, nil
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context) ([]*model.Category, error) {
	// Set variables
	var userAuth *users.User
	var graphqlCategories []*model.Category
	var userId string

	log.Printf("Route: Categories")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlCategories, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any category and user can only get his own categories
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Get all categories of user
	categories, err := categories.GetAllCategories(userId)
	if err != nil {
		log.Printf("Error: Get all accounts")
		return graphqlCategories, err
	}

	// Parse categories of user
	for _, category := range categories {
		graphqlCategory := &model.Category{
			ID:   category.ID,
			Name: category.Name,
		}
		graphqlCategories = append(graphqlCategories, graphqlCategory)
	}

	return graphqlCategories, nil
}

// Balances is the resolver for the balances field.
func (r *queryResolver) Balances(ctx context.Context, accountIds []string, categoryIds []string, initDate int, endDate int) ([]*model.Balance, error) {
	// Set variables
	var userAuth *users.User
	var account accounts.Account
	var category categories.Category
	var graphqlBalances []*model.Balance
	var userId string

	log.Printf("Route: Balances")

	// Parse dates
	date1 := time.UnixMilli(int64(initDate))
	date2 := time.UnixMilli(int64(endDate))

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlBalances, fmt.Errorf("access denied")
	}

	// Set fields. Admin can get any transaction and user can only get his own transactions
	if userAuth.Role == model.RoleAdmin {
		userId = ""
	} else {
		userId = userAuth.ID
	}

	// Start array to get total value of a category
	balanceAmounts := make([]float64, len(categoryIds))

	// Iterate every account
	for _, accountId := range accountIds {
		// Check if user can access account
		account.ID = accountId
		err := account.GetAccountById(userId)
		if err != nil {
			log.Printf("Error: Get account")
			return graphqlBalances, err
		}

		// Iterate every category
		for index, categoryId := range categoryIds {
			// Get all transactions of account with the category and with a valid date
			transactions, err := transactions.GetTransactionsByCategoryAndDate(accountId, categoryId, date1, date2)
			if err != nil {
				log.Printf("Error: Get transactions by category and date")
				return graphqlBalances, err
			}

			// For every transaction, add the amount to the total value of the category
			for _, transaction := range transactions {
				balanceAmounts[index] = balanceAmounts[index] + transaction.Amount
			}
		}
	}

	// Iterate every category
	for index, categoryId := range categoryIds {
		// Get category
		category.ID = categoryId
		err := category.GetCategoryById(userAuth.ID)
		if err != nil {
			log.Printf("Error: Get category")
			return graphqlBalances, err
		}

		// Add category and amount to create the balance
		graphqlBalance := &model.Balance{
			Amount: balanceAmounts[index],
			Category: &model.Category{
				ID:   category.ID,
				Name: category.Name,
			},
		}

		graphqlBalances = append(graphqlBalances, graphqlBalance)
	}

	return graphqlBalances, nil
}

// TokenWithCode is the resolver for the tokenWithCode field.
func (r *queryResolver) TokenWithCode(ctx context.Context, code string) (string, error) {
	// Set variables
	var userAuth *users.User

	log.Printf("Route: GetTokenWithCode")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return "Error", fmt.Errorf("access denied")
	}

	return santander.GetTokenWithCode(userAuth.ID, code)
}

// Category is the resolver for the category field.
func (r *transactionResolver) Category(ctx context.Context, obj *model.Transaction) (*model.Category, error) {
	var userAuth *users.User
	var transaction transactions.Transaction
	var category categories.Category

	log.Printf("Resolver: Category-Transaction")

	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.Category{}, fmt.Errorf("access denied")
	}

	transaction.ID = obj.ID

	// Get transaction in db
	err := transaction.GetTransactionById("")
	if err != nil {
		log.Printf("Error: Get transaction")
		return &model.Category{}, err
	}

	if transaction.Category == "" {
		return &model.Category{}, nil
	}

	category.ID = transaction.Category

	err = category.GetCategoryById(userAuth.ID)
	if err != nil {
		log.Printf("Error: Get category")
		return &model.Category{}, err
	}

	graphqlCategory := &model.Category{
		ID:   category.ID,
		Name: category.Name,
	}

	return graphqlCategory, nil
}

// User is the resolver for the user field.
func (r *transactionResolver) User(ctx context.Context, obj *model.Transaction) (*model.User, error) {
	// Set variables
	var userAuth *users.User
	var transaction transactions.Transaction
	var account accounts.Account
	var user users.User

	log.Printf("Resolver: User-Transaction")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.User{}, fmt.Errorf("access denied")
	}

	transaction.ID = obj.ID

	// Get transaction in db
	err := transaction.GetTransactionById("")
	if err != nil {
		log.Printf("Error: Get transaction")
		return &model.User{}, err
	}

	account.ID = transaction.AccountID

	// Get account in db
	err = account.GetAccountById("")
	if err != nil {
		log.Printf("Error: Get account")
		return &model.User{}, err
	}

	user.ID = account.UserID

	// Get user
	err = user.GetUserById()
	if err != nil {
		log.Printf("Error: Get user")
		return &model.User{}, err
	}

	// Parse to model.User
	graphqlUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}

	return graphqlUser, nil
}

// Account is the resolver for the account field.
func (r *transactionResolver) Account(ctx context.Context, obj *model.Transaction) (*model.Account, error) {
	// Set variables
	var userAuth *users.User
	var transaction transactions.Transaction
	var account accounts.Account

	log.Printf("Resolver: Account-Transaction")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return &model.Account{}, fmt.Errorf("access denied")
	}

	transaction.ID = obj.ID

	// Get transaction in db
	err := transaction.GetTransactionById("")
	if err != nil {
		log.Printf("Error: Get transaction")
		return &model.Account{}, err
	}

	account.ID = transaction.AccountID

	// Get account in db
	err = account.GetAccountById("")
	if err != nil {
		log.Printf("Error: Get account")
		return &model.Account{}, err
	}

	// Parse to model.Account
	graphqlAccount := &model.Account{
		ID:       account.ID,
		Iban:     account.Iban,
		Name:     account.Name,
		Currency: account.Currency,
		Amount:   account.Amount,
		Bank:     account.Bank,
	}

	return graphqlAccount, nil
}

// Accounts is the resolver for the accounts field.
func (r *userResolver) Accounts(ctx context.Context, obj *model.User) ([]*model.Account, error) {
	// Set variables
	var userAuth *users.User
	var graphqlAccounts []*model.Account

	log.Printf("Resolver: Accounts-User")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlAccounts, fmt.Errorf("access denied")
	}

	// Get all accounts of user
	accounts, err := accounts.GetAllAccounts(obj.ID)
	if err != nil {
		log.Printf("Error: Get all accounts")
		return graphqlAccounts, err
	}

	// Parse accounts of user
	for _, account := range accounts {
		graphqlAccount := &model.Account{
			ID:       account.ID,
			Iban:     account.Iban,
			Name:     account.Name,
			Currency: account.Currency,
			Amount:   account.Amount,
			Bank:     account.Bank,
		}
		graphqlAccounts = append(graphqlAccounts, graphqlAccount)
	}

	return graphqlAccounts, nil
}

// Transactions is the resolver for the transactions field.
func (r *userResolver) Transactions(ctx context.Context, obj *model.User) ([]*model.Transaction, error) {
	// Set variables
	var userAuth *users.User
	var graphqlTransactions []*model.Transaction

	log.Printf("Resolver: Transactions-User")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlTransactions, fmt.Errorf("access denied")
	}

	// Get all accounts of user
	accounts, err := accounts.GetAllAccounts(obj.ID)
	if err != nil {
		log.Printf("Error: Get all accounts")
		return graphqlTransactions, err
	}

	// Iterate for every account
	for _, account := range accounts {
		// Get all transactions of account
		transactions, err := transactions.GetAllTransactions(account.ID)
		if err != nil {
			log.Printf("Error: Get all transactions")
			return graphqlTransactions, err
		}

		// Parse transactions of account
		for _, transaction := range transactions {
			graphqlTransaction := &model.Transaction{
				ID:          transaction.ID,
				Description: transaction.Description,
				Date:        int(transaction.Date.UnixMilli()),
				Amount:      transaction.Amount,
			}
			graphqlTransactions = append(graphqlTransactions, graphqlTransaction)
		}
	}

	return graphqlTransactions, nil
}

// Account returns AccountResolver implementation.
func (r *Resolver) Account() AccountResolver { return &accountResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Transaction returns TransactionResolver implementation.
func (r *Resolver) Transaction() TransactionResolver { return &transactionResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type accountResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type transactionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
