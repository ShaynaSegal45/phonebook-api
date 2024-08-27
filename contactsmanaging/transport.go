package contactsmanaging

import (
	"context"
	"encoding/json"
	"strconv"

	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ShaynaSegal45/phonebook-api/contact"
	"github.com/ShaynaSegal45/phonebook-api/errors"
)

const (
	idParam       = "id"
	fullTextParam = "fullText"

	offsetParam = "offset"
	countParam  = "count"
	limitParam  = "limit"

	defaultOffset int = 0
	defaultCount  int = 2
)

type Service interface {
	Ping(ctx context.Context) string
	AddContact(ctx context.Context, c contact.Contact) (string, *errors.Error)
	GetContacts(ctx context.Context, limit, offset int, query string) ([]contact.Contact, *errors.Error)
	CountContacts(ctx context.Context, query string) (int, *errors.Error)
	GetContact(ctx context.Context, id string) (contact.Contact, *errors.Error)
	UpdateContact(ctx context.Context, id string, c contact.Contact) *errors.Error
	DeleteContact(ctx context.Context, id string) *errors.Error
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

type CreateContactRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

type UpdateContactRequest struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

type GetContactRequest struct {
	ID string `json:"id"`
}

type DeleteContactRequest struct {
	ID string `json:"id"`
}

type SearchContactsRequest struct {
	Text   string `json: "fullText"`
	Offset int    `json: "offset"`
	Limit  int    `json: "limit"`
}

type GetContactResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

func decodeAddContactRequest(r *http.Request) (interface{}, error) {
	var req CreateContactRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	return req, err
}

func decodeGetContactRequest(r *http.Request) (interface{}, error) {
	return GetContactRequest{
		ID: r.URL.Query().Get(idParam),
	}, nil
}

func decodeUpdateContactRequest(r *http.Request) (interface{}, error) {
	var req UpdateContactRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	req.ID = r.URL.Query().Get(idParam)
	return req, err
}

func decodeDeleteContactRequest(r *http.Request) (interface{}, error) {
	return DeleteContactRequest{
		ID: r.URL.Query().Get(idParam),
	}, nil
}

func decodeSearchContactsRequest(r *http.Request) (interface{}, error) {
	query := r.URL.Query().Get(fullTextParam)
	limitStr := r.URL.Query().Get(limitParam)
	offsetStr := r.URL.Query().Get(offsetParam)

	limit, parsingErr := strconv.Atoi(limitStr)
	if parsingErr != nil || limit <= 0 {
		limit = defaultCount
	}

	offset, parsingErr := strconv.Atoi(offsetStr)
	if parsingErr != nil || offset < 0 {
		offset = defaultOffset
	}

	return SearchContactsRequest{
		Text:   query,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func encodeAddContactResponse(w http.ResponseWriter, id string) {
	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func encodeGetContactResponse(w http.ResponseWriter, response interface{}) {
	res, ok := response.(GetContactResponse)
	if !ok {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// todo return id in service
func encodeUpdateContactResponse(w http.ResponseWriter, id string) {
	response := map[string]string{"id": id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func encodeDeleteContactResponse(w http.ResponseWriter) {
	response := map[string]string{}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
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
