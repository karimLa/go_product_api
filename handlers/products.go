package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/karimLa/microservices/data"
)

// Product is handler for products
type Product struct {
	l *log.Logger
}

// NewProduct creates a new product
func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Print("Handle GET Product")

	lp := data.GetProducts()

	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Product) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Print("Handle Post Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}

func (p *Product) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Print("Handle PUT Product")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(&prod, id)
	if err == data.ErrProductNotFound {
		http.Error(w, data.ErrProductNotFound.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, data.ErrProductNotFound.Error(), http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Product) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
