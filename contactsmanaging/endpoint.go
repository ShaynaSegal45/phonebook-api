package contactsmanaging

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ShaynaSegal45/phonebook-api/contact"
)

type Pagination struct {
	next        *paginationValues
	prev        *paginationValues
	count       int
	queryParams map[string]string
}

type paginationValues struct {
	limit  int64
	offset int64
}

type SearchContactsResponse struct {
	Contacts           []contact.Contact `json:"contacts"`
	Pagination         Pagination        `json:"pagination"`
	TotalContactsCount int               `json:"count"`
}

type Endpoints struct {
	AddContactEndpoint    http.HandlerFunc
	GetContactsEndpoint   http.HandlerFunc
	GetContactEndpoint    http.HandlerFunc
	UpdateContactEndpoint http.HandlerFunc
	DeleteContactEndpoint http.HandlerFunc
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		AddContactEndpoint:    makeAddContactEndpoint(s),
		GetContactsEndpoint:   makeGetContactsEndpoint(s),
		GetContactEndpoint:    makeGetContactEndpoint(s),
		UpdateContactEndpoint: makeUpdateContactEndpoint(s),
		DeleteContactEndpoint: makeDeleteContactEndpoint(s),
	}
}

func makeAddContactEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeAddContactRequest(r)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		req, ok := request.(CreateContactRequest)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		validationErr := req.Validate()
		if validationErr != nil {
			http.Error(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.AddContact(context.Background(), req.toContact())
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		encodeAddContactResponse(w, id)
	}
}

func makeGetContactEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeGetContactRequest(r)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		req, ok := request.(GetContactRequest)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		contact, err := s.GetContact(context.Background(), req.ID)
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}
		response := GetContactResponse{
			ID:        contact.ID,
			FirstName: contact.FirstName,
			LastName:  contact.LastName,
			Phone:     contact.Phone,
			Address:   contact.Address,
		}

		encodeGetContactResponse(w, response)
	}
}

func makeUpdateContactEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeUpdateContactRequest(r)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		req, ok := request.(UpdateContactRequest)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if validationErr := req.Validate(); validationErr != nil {
			http.Error(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		if err := s.UpdateContact(context.Background(), req.toContact()); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}
		encodeUpdateContactResponse(w)
	}
}

func makeDeleteContactEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeDeleteContactRequest(r)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}

		req, ok := request.(DeleteContactRequest)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := s.DeleteContact(context.Background(), req.ID); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}
		encodeDeleteContactResponse(w)
	}
}

func makeGetContactsEndpoint(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeSearchContactsRequest(r)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		req, ok := request.(SearchContactsRequest)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		contacts, err := s.GetContacts(context.Background(), req.toFilters())
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		totalContacts, err := s.CountContacts(context.Background(), req.Text)
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		pagination := createPagination(req.Offset, req.Limit, totalContacts, r.URL)

		response := SearchContactsResponse{
			Contacts:           contacts,
			Pagination:         pagination,
			TotalContactsCount: totalContacts,
		}

		encodeSearchContactsHandlerResponse(w, response)
	}
}

func createPagination(offset, limit, totalContacts int, url *url.URL) Pagination {
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

	queryParams := url.Query()
	queryParamsMap := make(map[string]string)
	for key, values := range queryParams {
		if len(values) > 0 {
			queryParamsMap[key] = values[0]
		}
	}

	return Pagination{
		next:        next,
		prev:        prev,
		count:       totalContacts,
		queryParams: queryParamsMap,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r CreateContactRequest) toContact() contact.Contact {
	return contact.Contact{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Phone:     r.Phone,
		Address:   r.Address,
	}
}

func (r UpdateContactRequest) toContact() contact.Contact {
	return contact.Contact{
		ID:        r.ID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Phone:     r.Phone,
		Address:   r.Address,
	}
}

func (r SearchContactsRequest) toFilters() contact.Filters {
	return contact.Filters{
		FullText: r.Text,
		Limit:    r.Limit,
		Offset:   r.Offset,
	}
}

func (r CreateContactRequest) Validate() error {
	if r.FirstName == "" && r.LastName == "" {
		return fmt.Errorf("CreateContactRequest.Validate: must include firstname or lastname")
	}

	return nil
}

func (r UpdateContactRequest) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("UpdateContactRequest.Validate: missing id")
	}

	return nil
}
