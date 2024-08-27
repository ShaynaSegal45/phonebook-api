package contactsmanaging

import (
	"context"
	"encoding/json"

	"fmt"
	"net/http"

	"github.com/ShaynaSegal45/phonebook-api/contact"
	"github.com/go-chi/chi/v5"
)

const (
	idParam       = "id"
	fullTextParam = "fullText"

	offsetParam = "offset"
	countParam  = "count"
	limitParam  = "limit"

	defaultOffsetStr = "0"
	defaultCountStr  = "2"

	defaultOffset = 0
	defaultCount  = 2
)

type Service interface {
	Ping(ctx context.Context) string
	AddContact(ctx context.Context, c contact.Contact) (string, error)
	GetContacts(ctx context.Context, limit, offset int, query string) ([]contact.Contact, error)
	CountContacts(ctx context.Context, query string) (int, error)
	GetContact(ctx context.Context, id string) (contact.Contact, error)
	UpdateContact(ctx context.Context, id string, c contact.Contact) error
	DeleteContact(ctx context.Context, id string) error
}

func NewHTTPHandler(s Service) http.Handler {
	router := chi.NewRouter()
	endpoint := MakeEndpoints(s)

	router.Post("/contact", endpoint.AddContactEndpoint)
	router.Get("/contacts", endpoint.GetContactsEndpoint)
	router.Get("/contact/{id}", endpoint.GetContactEndpoint)
	router.Put("/contact/{id}", endpoint.UpdateContactEndpoint)
	router.Delete("/contact/{id}", endpoint.DeleteContactEndpoint)
	router.Get("/ping", pingHandler)

	return router
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func encodeSearchContactsHandlerResponse(w http.ResponseWriter, response interface{}) {
	res, ok := response.(SearchContactsResponse)
	if !ok {
		http.Error(w, "Invalid response type", http.StatusInternalServerError)
		return
	}

	pagination := encodeSearchContactsPagination(context.Background(), res.Pagination, res.TotalContactsCount)

	contactsjson := map[string]interface{}{
		"contacts":   res.Contacts,
		"pagination": pagination,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(contactsjson); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func encodeSearchContactsPagination(ctx context.Context, pagination Pagination, totalContacts int) map[string]interface{} {
	baseURL := "/contacts"
	queryParamsStr := ""
	paginationDefaultQueries := make(map[string]bool)
	paginationDefaultQueries["count"] = true
	paginationDefaultQueries["limit"] = true
	paginationDefaultQueries["offset"] = true

	for k, v := range pagination.queryParams {
		_, ok := paginationDefaultQueries[k]
		if ok {
			continue
		}
		if queryParamsStr == "" {
			queryParamsStr = fmt.Sprintf("%s=%s", k, v)
		} else {
			queryParamsStr = fmt.Sprintf("%s&%s=%s", queryParamsStr, k, v)
		}
	}

	var next, prev string

	if pagination.next != nil {
		if queryParamsStr != "" {
			next = fmt.Sprintf("%s?%s&limit=%v&offset=%v&count=%v", baseURL, queryParamsStr, pagination.next.limit, pagination.next.offset, totalContacts)
		} else {
			next = fmt.Sprintf("%s?limit=%v&offset=%v&count=%v", baseURL, pagination.next.limit, pagination.next.offset, totalContacts)
		}
	}

	if pagination.prev != nil {
		if queryParamsStr != "" {
			prev = fmt.Sprintf("%s?%s&limit=%v&offset=%v&count=%v", baseURL, queryParamsStr, pagination.prev.limit, pagination.prev.offset, totalContacts)
		} else {
			prev = fmt.Sprintf("%s?limit=%v&offset=%v&count=%v", baseURL, pagination.prev.limit, pagination.prev.offset, totalContacts)
		}
	}

	return map[string]interface{}{
		"next":  next,
		"prev":  prev,
		"count": totalContacts,
	}
}
