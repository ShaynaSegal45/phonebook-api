package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ShaynaSegal45/phonebook-api/contact"
	"github.com/ShaynaSegal45/phonebook-api/errors"
)

const operationName = "contactsmanaging"

type ContactsRepo struct {
	db *sql.DB
}

func NewContactsRepo(db *sql.DB) *ContactsRepo {
	return &ContactsRepo{
		db: db,
	}
}

func (r *ContactsRepo) InsertContact(ctx context.Context, c contact.Contact) *errors.Error {
	query := `INSERT INTO contacts (id, firstname, lastname, address, phone) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.FirstName, c.LastName, c.Address, c.Phone)
	if err != nil {
		errMsg := fmt.Sprintf("ContactsRepo.InsertContact: failed to create contact with id %s", c.ID)
		log.Printf("%s: %v", errMsg, err)
		return errors.CreateError(operationName, errMsg, err, errors.InternalError)
	}

	return nil
}

func (r *ContactsRepo) GetContact(ctx context.Context, id string) (contact.Contact, *errors.Error) {
	query := `SELECT id, firstname, lastname, address, phone FROM contacts WHERE id = ?`
	var c contact.Contact
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Address, &c.Phone)
	if err != nil {
		errMsg := fmt.Sprintf("ContactsRepo.GetContact: failed to get contact with id %s", id)
		log.Printf("%s: %v", errMsg, err)
		return contact.Contact{}, errors.CreateError(operationName, errMsg, err, errors.InternalError)
	}
	return c, nil
}

func (r *ContactsRepo) SearchContacts(ctx context.Context, limit, offset int, query string) ([]contact.Contact, *errors.Error) {
	queryLike := `%` + query + `%`
	sqlQuery := `SELECT id, firstname, lastname, phone, address FROM contacts WHERE ? = "" 
                OR firstname LIKE ? OR lastname LIKE ? OR phone LIKE ?
                LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, sqlQuery, query, queryLike, queryLike, queryLike, limit, offset)
	if err != nil {
		errMsg := "ContactsRepo.SearchContacts"
		log.Printf("%s: failed to search contacts with query %s: %v", errMsg, query, err)
		return nil, errors.CreateError(operationName, errMsg, err, errors.InternalError)
	}
	defer rows.Close()

	var contacts []contact.Contact
	for rows.Next() {
		var c contact.Contact
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Address, &c.Phone); err != nil {
			errMsg := "ContactsRepo.SearchContacts error scanning rows"
			log.Printf("%s: failed to scan contact: %v", errMsg, err)
			return nil, errors.CreateError(operationName, errMsg, err, errors.InternalError)

		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *ContactsRepo) CountContacts(ctx context.Context, query string) (int, *errors.Error) {
	queryLike := `%` + query + `%`
	sqlQuery := `SELECT count(id) FROM contacts WHERE ? = "" OR
                  (firstname LIKE ? OR lastname LIKE ? OR phone LIKE ?)`

	var count int
	err := r.db.QueryRowContext(ctx, sqlQuery, query, queryLike, queryLike, queryLike).Scan(&count)
	if err != nil {
		errMsg := "ContactsRepo.CountContacts"
		log.Printf("%s: failed to count contacts with query %s: %v", errMsg, query, err)
		return 0, errors.CreateError(operationName, errMsg, err, errors.InternalError)

	}
	return count, nil
}

func (r *ContactsRepo) UpdateContact(ctx context.Context, id string, c contact.Contact) *errors.Error {
	query := `UPDATE contacts SET firstname = ?, lastname = ?, address = ?, phone = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, c.FirstName, c.LastName, c.Address, c.Phone, id)
	if err != nil {
		errMsg := "ContactsRepo.UpdateContact"
		log.Printf("%s: failed to update contact with id %s: %v", errMsg, id, err)
		return errors.CreateError(operationName, errMsg, err, errors.InternalError)

	}
	return nil
}

func (r *ContactsRepo) DeleteContact(ctx context.Context, id string) *errors.Error {
	query := `DELETE FROM contacts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		errMsg := "ContactsRepo.DeleteContact"
		log.Printf("%s: failed to delete contact with id %s: %v", errMsg, id, err)
		return errors.CreateError(operationName, errMsg, err, errors.InternalError)

	}
	return nil
}

func (r *ContactsRepo) ContactExists(ctx context.Context, firstName, lastName string) (bool, *errors.Error) {
	query := `SELECT 1 FROM contacts WHERE firstName = ? AND lastName = ?`
	var exists int
	err := r.db.QueryRowContext(ctx, query, firstName, lastName).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		errMsg := "ContactsRepo.ContactExists"
		log.Printf("%s: %v", errMsg, err)
		return false, errors.CreateError(operationName, "contacExists", err, errors.InternalError)
	}
	return true, nil
}
