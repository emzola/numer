package service_test

import (
	"context"
	"testing"

	"github.com/emzola/numer/user-service/internal/models"
	"github.com/emzola/numer/user-service/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// Mock the repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, email, password, role string) (*models.User, error) {
	args := m.Called(ctx, email, password, role)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) CreateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	args := m.Called(ctx, customer)
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockUserRepository) GetCustomerByID(ctx context.Context, customerID int64) (*models.Customer, error) {
	args := m.Called(ctx, customerID)
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockUserRepository) UpdateCustomer(ctx context.Context, customer *models.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteCustomer(ctx context.Context, customerID int64) error {
	args := m.Called(ctx, customerID)
	return args.Error(0)
}

// Unit tests for the UserService

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	email := "test@example.com"
	password := "password123"
	role := "user"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	expectedUser := &models.User{
		ID:             1,
		Email:          email,
		HashedPassword: string(hashedPassword),
		Role:           role,
	}

	// Mock the repository response
	mockRepo.On("CreateUser", mock.Anything, email, mock.AnythingOfType("string"), role).Return(expectedUser, nil)

	// Call the CreateUser method
	user, err := userService.CreateUser(context.Background(), email, password, role)

	// Assertions
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := int64(1)
	expectedUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
		Role:  "user",
	}

	// Mock the repository response
	mockRepo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)

	// Call the GetUserByID method
	user, err := userService.GetUserByID(context.Background(), userID)

	// Assertions
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	updatedUser := &models.User{
		ID:             1,
		Email:          "updated@example.com",
		HashedPassword: "newpassword",
		Role:           "admin",
	}

	// Mock the repository response
	mockRepo.On("UpdateUser", mock.Anything, updatedUser).Return(nil)

	// Call the UpdateUser method
	err := userService.UpdateUser(context.Background(), updatedUser)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := int64(1)

	// Mock the repository response
	mockRepo.On("DeleteUser", mock.Anything, userID).Return(nil)

	// Call the DeleteUser method
	err := userService.DeleteUser(context.Background(), userID)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateCustomer(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	customer := &models.Customer{
		ID:    1,
		Name:  "Customer Name",
		Email: "customer@example.com",
	}

	// Mock the repository response
	mockRepo.On("CreateCustomer", mock.Anything, customer).Return(customer, nil)

	// Call the CreateCustomer method
	createdCustomer, err := userService.CreateCustomer(context.Background(), customer)

	// Assertions
	require.NoError(t, err)
	require.Equal(t, customer, createdCustomer)
	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	customerID := int64(1)
	expectedCustomer := &models.Customer{
		ID:    customerID,
		Name:  "Customer Name",
		Email: "customer@example.com",
	}

	// Mock the repository response
	mockRepo.On("GetCustomerByID", mock.Anything, customerID).Return(expectedCustomer, nil)

	// Call the GetCustomerByID method
	customer, err := userService.GetCustomerByID(context.Background(), customerID)

	// Assertions
	require.NoError(t, err)
	require.Equal(t, expectedCustomer, customer)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomer(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	updatedCustomer := &models.Customer{
		ID:    1,
		Name:  "Updated Customer",
		Email: "updated@example.com",
	}

	// Mock the repository response
	mockRepo.On("UpdateCustomer", mock.Anything, updatedCustomer).Return(nil)

	// Call the UpdateCustomer method
	err := userService.UpdateCustomer(context.Background(), updatedCustomer)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCustomer(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	customerID := int64(1)

	// Mock the repository response
	mockRepo.On("DeleteCustomer", mock.Anything, customerID).Return(nil)

	// Call the DeleteCustomer method
	err := userService.DeleteCustomer(context.Background(), customerID)

	// Assertions
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
