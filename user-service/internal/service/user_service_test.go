package service_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/emzola/numer/user-service/internal/models"
// 	"github.com/emzola/numer/user-service/internal/service"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // Mock repository
// type mockUserRepo struct {
// 	mock.Mock
// }

// func (m *mockUserRepo) CreateUser(ctx context.Context, email, password, role string) (models.User, error) {
// 	args := m.Called(ctx, email, password, role)
// 	return args.Get(0).(models.User), args.Error(1)
// }

// func (m *mockUserRepo) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
// 	args := m.Called(ctx, userID)
// 	return args.Get(0).(models.User), args.Error(1)
// }

// func (m *mockUserRepo) UpdateUser(ctx context.Context, user models.User) error {
// 	return m.Called(ctx, user).Error(0)
// }

// func (m *mockUserRepo) DeleteUser(ctx context.Context, userID int64) error {
// 	return m.Called(ctx, userID).Error(0)
// }

// func (m *mockUserRepo) CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
// 	args := m.Called(ctx, customer)
// 	return args.Get(0).(models.Customer), args.Error(1)
// }

// func (m *mockUserRepo) GetCustomerByID(ctx context.Context, customerID int64) (models.Customer, error) {
// 	args := m.Called(ctx, customerID)
// 	return args.Get(0).(models.Customer), args.Error(1)
// }

// func (m *mockUserRepo) UpdateCustomer(ctx context.Context, customer models.Customer) error {
// 	return m.Called(ctx, customer).Error(0)
// }

// func (m *mockUserRepo) DeleteCustomer(ctx context.Context, customerID int64) error {
// 	return m.Called(ctx, customerID).Error(0)
// }

// func TestUserService_CreateUser(t *testing.T) {
// 	mockRepo := new(mockUserRepo)
// 	svc := service.NewUserService(mockRepo)

// 	user := models.User{
// 		ID:    1,
// 		Email: "test@example.com",
// 		Role:  "admin",
// 	}

// 	mockRepo.On("CreateUser", mock.Anything, user.Email, mock.Anything, user.Role).Return(user, nil)

// 	result, err := svc.CreateUser(context.Background(), user.Email, "password123", user.Role)

// 	assert.NoError(t, err)
// 	assert.Equal(t, user.ID, result.ID)
// 	assert.Equal(t, user.Email, result.Email)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUserService_GetUserByID(t *testing.T) {
// 	mockRepo := new(mockUserRepo)
// 	svc := service.NewUserService(mockRepo)

// 	user := models.User{
// 		ID:    1,
// 		Email: "test@example.com",
// 		Role:  "admin",
// 	}

// 	mockRepo.On("GetUserByID", mock.Anything, user.ID).Return(user, nil)

// 	result, err := svc.GetUserByID(context.Background(), user.ID)

// 	assert.NoError(t, err)
// 	assert.Equal(t, user.ID, result.ID)
// 	assert.Equal(t, user.Email, result.Email)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUserService_UpdateUser(t *testing.T) {
// 	mockRepo := new(mockUserRepo)
// 	svc := service.NewUserService(mockRepo)

// 	user := models.User{
// 		ID:    1,
// 		Email: "updated@example.com",
// 		Role:  "admin",
// 	}

// 	mockRepo.On("UpdateUser", mock.Anything, user).Return(nil)

// 	err := svc.UpdateUser(context.Background(), user)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUserService_DeleteUser(t *testing.T) {
// 	mockRepo := new(mockUserRepo)
// 	svc := service.NewUserService(mockRepo)

// 	userID := int64(1)

// 	mockRepo.On("DeleteUser", mock.Anything, userID).Return(nil)

// 	err := svc.DeleteUser(context.Background(), userID)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUserService_CreateCustomer(t *testing.T) {
// 	mockRepo := new(mockUserRepo)
// 	svc := service.NewUserService(mockRepo)

// 	customer := models.Customer{
// 		ID:     1,
// 		UserID: 1,
// 		Name:   "Customer A",
// 		Email:  "customer@example.com",
// 	}

// 	mockRepo.On("CreateCustomer", mock.Anything, customer).Return(customer, nil)

// 	result, err := svc.CreateCustomer(context.Background(), customer)

// 	assert.NoError(t, err)
// 	assert.Equal(t, customer.ID, result.ID)
// 	assert.Equal(t, customer.Email, result.Email)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUserService_GetCustomerByID(t *testing.T) {
// 	mockRepo := new(mockUserRepo)
// 	svc := service.NewUserService(mockRepo)

// 	customer := models.Customer{
// 		ID:     1,
// 		UserID: 1,
// 		Name:   "Customer A",
// 		Email:  "customer@example.com",
// 	}

// 	mockRepo.On("GetCustomerByID", mock.Anything, customer.ID).Return(customer, nil)

// 	result, err := svc.GetCustomerByID(context.Background(), customer.ID)

// 	assert.NoError(t, err)
// 	assert.Equal(t, customer.ID, result.ID)
// 	assert.Equal(t, customer.Email, result.Email)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUserService_UpdateCustomer(t *testing.T) {
// 	mockRepo := new(mockUserRepo)
// 	svc := service.NewUserService(mockRepo)

// 	customer := models.Customer{
// 		ID:    1,
// 		Name:  "Updated Customer",
// 		Email: "updated@example.com",
// 	}

// 	mockRepo.On("UpdateCustomer", mock.Anything, customer).Return(nil)

// 	err := svc.UpdateCustomer(context.Background(), customer)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUserService_DeleteCustomer(t *testing.T) {
// 	mockRepo := new(mockUserRepo)
// 	svc := service.NewUserService(mockRepo)

// 	customerID := int64(1)

// 	mockRepo.On("DeleteCustomer", mock.Anything, customerID).Return(nil)

// 	err := svc.DeleteCustomer(context.Background(), customerID)

// 	assert.NoError(t, err)
// 	mockRepo.AssertExpectations(t)
// }
