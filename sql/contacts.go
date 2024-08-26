package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ShaynaSegal45/phonebook-api/contact"
)

type ContactsRepo struct {
	db *sql.DB
	// ContactsRepo
}

// NewContactsRepo initializes the SQLite database connection and returns a contactsRepo instance
func NewContactsRepo(db *sql.DB) *ContactsRepo {
	return &ContactsRepo{
		db: db,
	}
}

// InsertContact inserts a new contact into the contacts table
func (r *ContactsRepo) InsertContact(ctx context.Context, c contact.Contact) error {
	query := `INSERT INTO contacts (id, firstname, lastname, address, phone) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.FirstName, c.LastName, c.Address, c.Phone)
	if err != nil {
		return fmt.Errorf("failed to insert contact: %w", err)
	}
	return nil
}

// GetContact retrieves a contact by ID from the contacts table
func (r *ContactsRepo) GetContact(ctx context.Context, id string) (contact.Contact, error) {
	query := `SELECT id, firstname, lastname, address, phone FROM contacts WHERE id = ?`
	var c contact.Contact
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Address, &c.Phone)
	if err != nil {
		return contact.Contact{}, fmt.Errorf("failed to get contact: %w", err)
	}
	return c, nil
}

// SearchContacts searches for contacts by first name or last name
func (r *ContactsRepo) SearchContacts(ctx context.Context, query string) ([]contact.Contact, error) {
	queryLike := `%` + query + `%`
	sqlQuery := `SELECT id, firstname, 
                lastname, phone, address FROM contacts WHERE ? = "" 
                OR firstname LIKE ? OR lastname LIKE ? OR phone LIKE ?`
	rows, err := r.db.QueryContext(ctx, sqlQuery, query, queryLike, queryLike, queryLike)
	if err != nil {
		return nil, fmt.Errorf("failed to search contacts: %w", err)
	}
	defer rows.Close()

	var contacts []contact.Contact
	for rows.Next() {
		var c contact.Contact
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Address, &c.Phone); err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

// UpdateContact updates an existing contact in the contacts table
func (r *ContactsRepo) UpdateContact(ctx context.Context, id string, c contact.Contact) error {
	query := `UPDATE contacts SET firstname = ?, lastname = ?, address = ?, phone = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, c.FirstName, c.LastName, c.Address, c.Phone, id)
	if err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}
	return nil
}

// DeleteContact deletes a contact by ID from the contacts table
func (r *ContactsRepo) DeleteContact(ctx context.Context, id string) error {
	query := `DELETE FROM contacts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	return nil
}

func (r *ContactsRepo) ContactExists(ctx context.Context, firstName, lastName string) (bool, error) {

	query := `DELETE FROM contacts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, firstName)
	if err != nil {
		return true, nil
	}
	return true, nil
}

//type ContactsRepo interface {
//InsertContact(ctx context.Context, c contact.Contact) error
//GetContact(ctx context.Context, id string) (contact.Contact, error)
//SearchContacts(ctx context.Context, query string) ([]contact.Contact, error)
//UpdateContact(ctx context.Context, id string, c contact.Contact) error
//DeleteContact(ctx context.Context, id string) error
//ContactExists(ctx context.Context, firstName, lastName string) (bool, error)
//  GetContacts(ctx context.Context, limit, offset int) ([]contact.Contact, error)
//}
