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

type Endpoints struct {
	AddContactEndpoint  http.HandlerFunc
	GetContactsEndpoint http.HandlerFunc
	//SearchContactEndpoint http.HandlerFunc
	UpdateContactEndpoint http.HandlerFunc
	DeleteContactEndpoint http.HandlerFunc
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		AddContactEndpoint:  makeAddContactEndpoint(s),
		GetContactsEndpoint: makeGetContactsEndpoint(s),
		//SearchContactEndpoint: makeSearchContactEndpoint(s),
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
		query := r.URL.Query().Get("text")
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		contacts, err := s.GetContacts(context.Background(), query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contacts); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

// func makeSearchContactEndpoint(s Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		query := r.URL.Query().Get("q")
// 		contacts, err := s.SearchContact(context.Background(), query)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		if err := json.NewEncoder(w).Encode(contacts); err != nil {
// 			log.Printf("Error encoding response: %v", err)
// 		}
// 	}
// }

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
