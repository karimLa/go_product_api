package handlers

import (
	"log"
	"net/http"

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

func (p *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	// Catch all
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}
