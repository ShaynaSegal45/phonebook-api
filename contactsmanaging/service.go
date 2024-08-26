package contactsmanaging

import (
	"context"

	"github.com/google/uuid"

	"github.com/ShaynaSegal45/phonebook-api/contact"
)

const pageSize = 2 // TODO: Change to 10

type ContactsRepo interface {
	InsertContact(ctx context.Context, c contact.Contact) error
	GetContact(ctx context.Context, id string) (contact.Contact, error)
	SearchContacts(ctx context.Context, query string) ([]contact.Contact, error)
	UpdateContact(ctx context.Context, id string, c contact.Contact) error
	DeleteContact(ctx context.Context, id string) error
	ContactExists(ctx context.Context, firstName, lastName string) (bool, error)
}

// service struct implements the Service interface
type service struct {
	repo ContactsRepo
}

// NewService creates a new service instance with the given repository
func NewService(repo ContactsRepo) Service {
	return &service{repo: repo}
}

// Ping responds with "pong"
func (s *service) Ping(ctx context.Context) string {
	return "pong"
}

// AddContact adds a new contact and returns the generated ID
func (s *service) AddContact(ctx context.Context, c contact.Contact) (string, error) {
	// exists, err := s.repo.ContactExists(ctx, c.FirstName, c.FirstName)
	// if err != nil {
	// 	return "", err
	// }

	// if exists {
	// 	return "", fmt.Errorf("contact with name %s %s already exists", c.FirstName, c.LastName)
	// }

	id := generateUniqueID()
	c.ID = id

	err := s.repo.InsertContact(ctx, c)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetContacts retrieves contacts for the specified page
func (s *service) GetContacts(ctx context.Context, query string) ([]contact.Contact, error) {
	contacts, err := s.repo.SearchContacts(ctx, query)
	if err != nil {

	}
	//return s.repo.SearchContacts(ctx,)
	return contacts, nil
}

// SearchContact searches for contacts matching the query
// func (s *service) SearchContact(ctx context.Context, query string) ([]contact.Contact, error) {
// 	//return s.repo.SearchContact(ctx,query)
// 	return nil, nil
// }

// UpdateContact updates the contact with the given ID
func (s *service) UpdateContact(ctx context.Context, id string, updatedContact contact.Contact) error {
	return s.repo.UpdateContact(ctx, id, updatedContact)
}

// DeleteContact deletes the contact with the given ID
func (s *service) DeleteContact(ctx context.Context, id string) error {
	return s.repo.DeleteContact(ctx, id)
}

func generateUniqueID() string {
	return uuid.New().String()
}
