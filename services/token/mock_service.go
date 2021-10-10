// Code generated by mockery v1.0.0. DO NOT EDIT.

package token

import (
	sql "database/sql"

	models "github.com/artem-malko/auth-and-go/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// Create provides a mock function with given fields: tokenType, clientID, accountID, identityID, tx
func (_m *MockService) Create(tokenType models.TokenType, clientID models.ClientID, accountID uuid.UUID, identityID uuid.UUID, tx *sql.Tx) (*models.Token, error) {
	ret := _m.Called(tokenType, clientID, accountID, identityID, tx)

	var r0 *models.Token
	if rf, ok := ret.Get(0).(func(models.TokenType, models.ClientID, uuid.UUID, uuid.UUID, *sql.Tx) *models.Token); ok {
		r0 = rf(tokenType, clientID, accountID, identityID, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.TokenType, models.ClientID, uuid.UUID, uuid.UUID, *sql.Tx) error); ok {
		r1 = rf(tokenType, clientID, accountID, identityID, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteExpiredTokens provides a mock function with given fields: tokenType, tx
func (_m *MockService) DeleteExpiredTokens(tokenType models.TokenType, tx *sql.Tx) ([]*models.Token, error) {
	ret := _m.Called(tokenType, tx)

	var r0 []*models.Token
	if rf, ok := ret.Get(0).(func(models.TokenType, *sql.Tx) []*models.Token); ok {
		r0 = rf(tokenType, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.TokenType, *sql.Tx) error); ok {
		r1 = rf(tokenType, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUsedTokens provides a mock function with given fields:
func (_m *MockService) DeleteUsedTokens() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetActiveTokenByIdentityID provides a mock function with given fields: identityID, tokenType
func (_m *MockService) GetActiveTokenByIdentityID(identityID uuid.UUID, tokenType models.TokenType) (*models.Token, error) {
	ret := _m.Called(identityID, tokenType)

	var r0 *models.Token
	if rf, ok := ret.Get(0).(func(uuid.UUID, models.TokenType) *models.Token); ok {
		r0 = rf(identityID, tokenType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, models.TokenType) error); ok {
		r1 = rf(identityID, tokenType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Use provides a mock function with given fields: tokenID, tx
func (_m *MockService) Use(tokenID uuid.UUID, tx *sql.Tx) (*models.Token, error) {
	ret := _m.Called(tokenID, tx)

	var r0 *models.Token
	if rf, ok := ret.Get(0).(func(uuid.UUID, *sql.Tx) *models.Token); ok {
		r0 = rf(tokenID, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, *sql.Tx) error); ok {
		r1 = rf(tokenID, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}