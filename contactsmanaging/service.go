package contactsmanaging

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ShaynaSegal45/phonebook-api/contact"
)

//const pageSize = 2 // TODO: Change to 10

type ContactsRepo interface {
	InsertContact(ctx context.Context, c contact.Contact) error
	GetContact(ctx context.Context, id string) (contact.Contact, error)

	SearchContacts(ctx context.Context, limit, offset int, query string) ([]contact.Contact, error)
	CountContacts(ctx context.Context, query string) (int, error)

	UpdateContact(ctx context.Context, id string, c contact.Contact) error
	DeleteContact(ctx context.Context, id string) error
	ContactExists(ctx context.Context, firstName, lastName string) (bool, error)
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

func (s *service) AddContact(ctx context.Context, c contact.Contact) (string, error) {
	exists, err := s.repo.ContactExists(ctx, c.FirstName, c.LastName)
	if err != nil {
		return "", err
	}

	if exists {
		return "", fmt.Errorf("contact with name %s %s already exists", c.FirstName, c.LastName)
	}

	id := generateUniqueID()
	c.ID = id

	err = s.repo.InsertContact(ctx, c)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *service) GetContacts(ctx context.Context, limit, offset int, query string) ([]contact.Contact, error) {
	contacts, err := s.repo.SearchContacts(ctx, limit, offset, query)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *service) CountContacts(ctx context.Context, query string) (int, error) {
	count, err := s.repo.CountContacts(ctx, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *service) UpdateContact(ctx context.Context, id string, updatedContact contact.Contact) error {
	return s.repo.UpdateContact(ctx, id, updatedContact)
}

func (s *service) DeleteContact(ctx context.Context, id string) error {
	return s.repo.DeleteContact(ctx, id)
}

func generateUniqueID() string {
	return uuid.New().String()
}
