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
	var account *accounts.Account
	var user *users.User

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
			Category:    transaction.Category,
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
				Category:    transaction.Category,
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
	var graphqlTransaction *model.Transaction
	var userId string

	log.Printf("Route: TransactionByID")

	// Get user from context
	if userAuth = middleware.ForContext(ctx); userAuth == nil {
		log.Printf("Error: Access denied")
		return graphqlTransaction, fmt.Errorf("access denied")
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
		return graphqlTransaction, err
	}

	// Iterate for every account
	for _, account := range accounts {
		// Get transaction by id
		err := transaction.GetTransactionById(account.ID)
		if err != nil {
			log.Printf("Error: Get account")
			continue
		}

		// Parse to model.Transaction
		graphqlTransaction = &model.Transaction{
			ID:          transaction.ID,
			Description: transaction.Description,
			Date:        int(transaction.Date.UnixMilli()),
			Amount:      transaction.Amount,
			Category:    transaction.Category,
		}

		break
	}

	if graphqlTransaction.Account == nil {
		log.Printf("Error: Get account")
		return graphqlTransaction, fmt.Errorf("not found")
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
				Category:    transaction.Category,
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
				Category:    transaction.Category,
			}
			graphqlTransactions = append(graphqlTransactions, graphqlTransaction)
		}
	}

	return graphqlTransactions, nil
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

// User is the resolver for the user field.
func (r *transactionResolver) User(ctx context.Context, obj *model.Transaction) (*model.User, error) {
	// Set variables
	var userAuth *users.User
	var transaction *transactions.Transaction
	var account *accounts.Account
	var user *users.User

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
	var transaction *transactions.Transaction
	var account *accounts.Account

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
				Category:    transaction.Category,
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
