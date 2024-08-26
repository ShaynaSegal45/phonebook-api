package contactsmanaging

import (
	"context"
	"net/http"

	"github.com/ShaynaSegal45/phonebook-api/contact"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	Ping(ctx context.Context) string
	AddContact(ctx context.Context, c contact.Contact) (string, error)
	GetContacts(ctx context.Context, query string) ([]contact.Contact, error)
	//SearchContact(ctx context.Context, query string) ([]contact.Contact, error)
	UpdateContact(ctx context.Context, id string, c contact.Contact) error
	DeleteContact(ctx context.Context, id string) error
}

func NewHTTPHandler(s Service) http.Handler {
	router := chi.NewRouter()
	endpoint := MakeEndpoints(s)

	router.Post("/contact", endpoint.AddContactEndpoint)
	router.Get("/contacts", endpoint.GetContactsEndpoint)
	//router.Get("/contact/search", endpoint.SearchContactEndpoint)
	router.Put("/contact/{id}", endpoint.UpdateContactEndpoint)
	router.Delete("/contact/{id}", endpoint.DeleteContactEndpoint)
	router.Get("/ping", pingHandler)

	return router
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
