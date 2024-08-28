package contactsmanaging

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ShaynaSegal45/phonebook-api/contact"
	"github.com/ShaynaSegal45/phonebook-api/errors"
)

type MockContactsRepo struct {
	mock.Mock
}

func (m *MockContactsRepo) GetContact(ctx context.Context, id string) (contact.Contact, *errors.Error) {
	args := m.Called(ctx, id)
	var c contact.Contact
	var err *errors.Error

	if contactArg := args.Get(0); contactArg != nil {
		c, _ = contactArg.(contact.Contact)
	}

	if errArg := args.Get(1); errArg != nil {
		err, _ = errArg.(*errors.Error)
		return c, err
	}

	return c, nil
}

func (m *MockContactsRepo) ContactExists(ctx context.Context, firstName, lastName string) (bool, *errors.Error) {
	args := m.Called(ctx, firstName, lastName)
	var err *errors.Error
	if e := args.Get(1); e != nil {
		err = e.(*errors.Error)
	}
	return args.Bool(0), err
}

func (m *MockContactsRepo) SearchContacts(ctx context.Context, f contact.Filters) ([]contact.Contact, *errors.Error) {
	args := m.Called(ctx, f.Limit, f.Offset, f.FullText)
	return args.Get(0).([]contact.Contact), args.Get(1).(*errors.Error)
}

func (m *MockContactsRepo) CountContacts(ctx context.Context, query string) (int, *errors.Error) {
	args := m.Called(ctx, query)
	return args.Int(0), args.Get(1).(*errors.Error)
}

func (m *MockContactsRepo) UpdateContact(ctx context.Context, c contact.Contact) *errors.Error {
	args := m.Called(ctx, c)
	return args.Get(0).(*errors.Error)
}

func (m *MockContactsRepo) DeleteContact(ctx context.Context, id string) *errors.Error {
	args := m.Called(ctx, id)
	return args.Get(0).(*errors.Error)
}

func (m *MockContactsRepo) InsertContact(ctx context.Context, c contact.Contact) *errors.Error {
	args := m.Called(ctx, c)
	return args.Get(0).(*errors.Error)
}

func TestGetContact_NotFound(t *testing.T) {
	repo := new(MockContactsRepo)
	service := NewService(repo)

	repo.On("GetContact", mock.Anything, "123").Return(contact.Contact{}, errors.CreateError("contactsmanaging", "GetContact", fmt.Errorf("not found"), errors.NotFoundError))

	c, err := service.GetContact(context.Background(), "123")

	assert.Error(t, err)
	assert.Equal(t, contact.Contact{}, c)
	repo.AssertExpectations(t)
}

func TestAddContact_Conflict(t *testing.T) {
	repo := new(MockContactsRepo)
	service := NewService(repo)

	contactToAdd := contact.Contact{FirstName: "John", LastName: "Doe"}
	repo.On("ContactExists", mock.Anything, contactToAdd.FirstName, contactToAdd.LastName).Return(true, nil)

	id, err := service.AddContact(context.Background(), contactToAdd)

	assert.Error(t, err)
	assert.Empty(t, id)
	repo.AssertExpectations(t)
}

//add more tests
// func TestAddContact_Success(t *testing.T) {
// 	repo := new(MockContactsRepo)
// 	service := NewService(repo)

// 	contactToAdd := contact.Contact{FirstName: "John", LastName: "Doe"}
// 	repo.On("ContactExists", mock.Anything, contactToAdd.FirstName, contactToAdd.LastName).Return(false, nil)
// 	repo.On("InsertContact", mock.Anything, contactToAdd).Return(nil)

// 	id, err := service.AddContact(context.Background(), contactToAdd)

// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, id)
// 	repo.AssertExpectations(t)
// }

// func TestGetContact_Success(t *testing.T) {
// 	repo := new(MockContactsRepo)
// 	service := NewService(repo)

// 	expectedContact := contact.Contact{ID: "123", FirstName: "John", LastName: "Doe"}
// 	repo.On("GetContact", mock.Anything, "123").Return(expectedContact, nil)

// 	c, err := service.GetContact(context.Background(), "123")

// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedContact, c)
// 	repo.AssertExpectations(t)
// }
