package contactsmanaging

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ShaynaSegal45/phonebook-api/contact"
	"github.com/ShaynaSegal45/phonebook-api/errors"
)

const operationName = "contactsmanaging"

type ContactsRepo interface {
	InsertContact(ctx context.Context, c contact.Contact) *errors.Error
	GetContact(ctx context.Context, id string) (contact.Contact, *errors.Error)
	SearchContacts(ctx context.Context, limit, offset int, query string) ([]contact.Contact, *errors.Error)
	CountContacts(ctx context.Context, query string) (int, *errors.Error)
	UpdateContact(ctx context.Context, id string, c contact.Contact) *errors.Error
	DeleteContact(ctx context.Context, id string) *errors.Error
	ContactExists(ctx context.Context, firstName, lastName string) (bool, *errors.Error)
}

type service struct {
	repo ContactsRepo
}

func NewService(repo ContactsRepo) Service {
	return &service{repo: repo}
}

func (s *service) Ping(ctx context.Context) string {
	return "pong"
}

func (s *service) AddContact(ctx context.Context, c contact.Contact) (string, *errors.Error) {
	exists, err := s.repo.ContactExists(ctx, c.FirstName, c.LastName)
	if err != nil {
		return "", err.ErrorWrapper(operationName, "AddContact")
	}

	if exists {
		conflictErr := fmt.Errorf("contact with name %s %s already exists", c.FirstName, c.LastName)
		return "", errors.CreateError(operationName, "AddContact", conflictErr, errors.ConflictError)
	}

	id := generateUniqueID()
	c.ID = id

	if err := s.repo.InsertContact(ctx, c); err != nil {
		return "", err.ErrorWrapper(operationName, "AddContact")
	}

	return id, nil
}

func (s *service) GetContacts(ctx context.Context, limit, offset int, query string) ([]contact.Contact, *errors.Error) {
	contacts, err := s.repo.SearchContacts(ctx, limit, offset, query)
	if err != nil {
		return nil, err.ErrorWrapper(operationName, "GetContacts")
	}
	return contacts, nil
}

func (s *service) CountContacts(ctx context.Context, query string) (int, *errors.Error) {
	count, err := s.repo.CountContacts(ctx, query)
	if err != nil {
		return 0, err.ErrorWrapper(operationName, "CountContacts")

	}
	return count, nil
}

func (s *service) GetContact(ctx context.Context, id string) (contact.Contact, *errors.Error) {
	c, err := s.repo.GetContact(ctx, id)
	if err != nil {
		return contact.Contact{}, err.ErrorWrapper(operationName, "GetContact")
	}
	return c, nil
}

func (s *service) UpdateContact(ctx context.Context, id string, updatedContact contact.Contact) *errors.Error {
	//todo check if firstname lastname exists first
	if err := s.repo.UpdateContact(ctx, id, updatedContact); err != nil {
		return err.ErrorWrapper(operationName, "UpdateContact")
	}
	return nil
}

func (s *service) DeleteContact(ctx context.Context, id string) *errors.Error {
	if err := s.repo.DeleteContact(ctx, id); err != nil {
		return err.ErrorWrapper(operationName, "DeleteContact")
	}
	return nil
}

func generateUniqueID() string {
	return uuid.New().String()
}
