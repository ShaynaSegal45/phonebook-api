package contactsmanaging

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ShaynaSegal45/phonebook-api/contact"
)

type Pagination struct {
	next  *paginationValues `json:"next,omitempty"`
	prev  *paginationValues `json:"prev,omitempty"`
	count int               `json:"count"`
}

type paginationValues struct {
	limit  int64 `json:"limit"`
	offset int64 `json:"offset"`
}

type SearchContactsResponse struct {
	Contacts           []contact.Contact `json:"contacts"`
	Pagination         Pagination        `json:"pagination"`
	TotalContactsCount int               `json:"count"`
}

type Endpoints struct {
	AddContactEndpoint    http.HandlerFunc
	GetContactsEndpoint   http.HandlerFunc
	UpdateContactEndpoint http.HandlerFunc
	DeleteContactEndpoint http.HandlerFunc
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		AddContactEndpoint:    makeAddContactEndpoint(s),
		GetContactsEndpoint:   makeGetContactsEndpoint(s),
		UpdateContactEndpoint: makeUpdateContactEndpoint(s),
		DeleteContactEndpoint: makeDeleteContactEndpoint(s),
	}
}

func makeAddContactEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newContact contact.Contact
		if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.AddContact(context.Background(), newContact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Generated ID: %s", id)

		response := map[string]string{
			"id": id,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func makeGetContactsEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get(fullTextParam)
		limitStr := r.URL.Query().Get(limitParam)
		offsetStr := r.URL.Query().Get(offsetParam)

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = defaultCount
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = defaultOffset
		}

		contacts, err := s.GetContacts(context.Background(), limit, offset, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		totalContacts, err := s.CountContacts(context.Background(), query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pagination := createPagination(offset, limit, totalContacts)

		// Extract query parameters
		queryParams := r.URL.Query()

		// Convert query parameters to a map
		queryParamsMap := make(map[string]string)
		for key, values := range queryParams {
			if len(values) > 0 {
				queryParamsMap[key] = values[0]
			}
		}

		//encodeSearchContactsPagination(context.Background(), pagination, totalContacts, queryParamsMap)

		// Create response
		response := SearchContactsResponse{
			Contacts:           contacts,
			Pagination:         pagination,
			TotalContactsCount: totalContacts,
		}

		encodeSearchContactsHandlerResponse(w, response, queryParamsMap)
	}
}

// func makeGetContactsEndpoint(s Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		query := r.URL.Query().Get(fullTextParam)
// 		limitStr := r.URL.Query().Get(limitParam)
// 		offsetStr := r.URL.Query().Get(offsetParam)

// 		limit, err := strconv.Atoi(limitStr)
// 		if err != nil || limit <= 0 {
// 			limit = defaultCount
// 		}

// 		offset, err := strconv.Atoi(offsetStr)
// 		if err != nil || offset < 0 {
// 			offset = defaultOffset
// 		}

// 		contacts, err := s.GetContacts(context.Background(), limit, offset, query)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		totalContacts, err := s.CountContacts(context.Background(), query)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		pagination := createPagination(offset, limit, totalContacts)

// 		response := SearchContactsResponse{
// 			Contacts:           contacts,
// 			Pagination:         pagination,
// 			TotalContactsCount: totalContacts,
// 		}

// 		encodeSearchContactsHandlerResponse(w, response)
// 	}
// }

func createPagination(offset, limit, totalContacts int) Pagination {
	var next, prev *paginationValues

	if offset+limit < totalContacts {
		next = &paginationValues{
			limit:  int64(limit),
			offset: int64(offset + limit),
		}
	}

	if offset > 0 {
		prev = &paginationValues{
			limit:  int64(limit),
			offset: int64(max(0, offset-limit)),
		}
	}

	return Pagination{
		next:  next,
		prev:  prev,
		count: totalContacts,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func makeUpdateContactEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var updatedContact contact.Contact
		if err := json.NewDecoder(r.Body).Decode(&updatedContact); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.UpdateContact(context.Background(), id, updatedContact); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func makeDeleteContactEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := s.DeleteContact(context.Background(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
