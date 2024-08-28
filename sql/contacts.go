package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/ShaynaSegal45/phonebook-api/contact"
	"github.com/ShaynaSegal45/phonebook-api/errors"
)

const (
	operationName = "contactsmanaging"
	ttl           = 300 * time.Second
)

type ContactsRepo struct {
	db    *sql.DB
	cache *redis.Client
}

func NewContactsRepo(db *sql.DB, cache *redis.Client) *ContactsRepo {
	return &ContactsRepo{
		db:    db,
		cache: cache,
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

func (r *ContactsRepo) SearchContacts(ctx context.Context, f contact.Filters) ([]contact.Contact, *errors.Error) {
	queryLike := `%` + f.FullText + `%`
	sqlQuery := `SELECT id, firstname, lastname, phone, address FROM contacts WHERE ? = "" 
                OR firstname LIKE ? OR lastname LIKE ? OR phone LIKE ?
				ORDER BY lastname, firstname 
                LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, sqlQuery, f.FullText, queryLike, queryLike, queryLike, f.Limit, f.Offset)
	if err != nil {
		errMsg := "ContactsRepo.SearchContacts"
		if err == sql.ErrNoRows {
			log.Printf("%s: contact not found", errMsg)
			return nil, errors.CreateError(operationName, errMsg, err, errors.NotFoundError)
		}
		log.Printf("%s: failed to search contacts with query %s: %v", errMsg, f.FullText, err)
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

func (r *ContactsRepo) UpdateContact(ctx context.Context, c contact.Contact) *errors.Error {
	query, args := buildUpdateQuery(c)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		errMsg := "ContactsRepo.UpdateContact"
		log.Printf("%s: failed to update contact with id %s: %v", errMsg, c.ID, err)
		return errors.CreateError("UpdateContact", errMsg, err, errors.InternalError)
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

func buildUpdateQuery(c contact.Contact) (string, []interface{}) {
	query := `UPDATE contacts SET`
	var args []interface{}

	if c.FirstName != "" {
		query += ` firstname = ?,`
		args = append(args, c.FirstName)
	}
	if c.LastName != "" {
		query += ` lastname = ?,`
		args = append(args, c.LastName)
	}
	if c.Address != "" {
		query += ` address = ?,`
		args = append(args, c.Address)
	}
	if c.Phone != "" {
		query += ` phone = ?,`
		args = append(args, c.Phone)
	}

	query = query[:len(query)-1] + ` WHERE id = ?`
	args = append(args, c.ID)

	return query, args
}

func (r *ContactsRepo) GetContact(ctx context.Context, id string) (contact.Contact, *errors.Error) {
	cachedContact, err := r.cache.Get(ctx, id).Result()
	if err == redis.Nil {

		var c contact.Contact
		query := `SELECT id, firstname, lastname, address, phone FROM contacts WHERE id = ?`
		err = r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Address, &c.Phone)
		if err != nil {
			errMsg := fmt.Sprintf("ContactsRepo.GetContact: failed to get contact with id %s", id)
			if err == sql.ErrNoRows {
				log.Printf("%s: contact not found", errMsg)
				return contact.Contact{}, errors.CreateError(operationName, errMsg, err, errors.NotFoundError)
			}
			log.Printf("%s: %v", errMsg, err)
			return contact.Contact{}, errors.CreateError(operationName, errMsg, err, errors.InternalError)
		}

		err = r.cache.Set(ctx, id, fmt.Sprintf("%s,%s,%s,%s,%s", c.ID, c.FirstName, c.LastName, c.Address, c.Phone), ttl).Err()
		if err != nil {
			log.Printf("Failed to set cache for contact id %s: %v", id, err)
		}

		return c, nil

	} else if err != nil {
		return contact.Contact{}, errors.CreateError(operationName, "failed to get cache", err, errors.InternalError)
	}

	contact := DeserializeContact(cachedContact)
	return contact, nil
}

func DeserializeContact(data string) contact.Contact {

	fields := strings.Split(data, ",")
	return contact.Contact{
		ID:        fields[0],
		FirstName: fields[1],
		LastName:  fields[2],
		Address:   fields[3],
		Phone:     fields[4],
	}
}
