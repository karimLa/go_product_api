package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		// expect the id in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if g == nil || len(g) != 1 {
			http.Error(w, "invalid uri", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Printf("gourp 0 %v", g[0])
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "invalid uri", http.StatusBadRequest)
			return
		}

		p.l.Println("ID is: ", id)

		p.updateProduct(w, r, id)
	}

	// Catch all
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Print("Handle GET Product")

	lp := data.GetProducts()

	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Product) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Print("Handle Post Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
		return
	}

	data.AddProduct(prod)
}

func (p *Product) updateProduct(w http.ResponseWriter, r *http.Request, id int) {
	p.l.Print("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(prod, id)
	if err == data.ErrProductNotFound {
		http.Error(w, data.ErrProductNotFound.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, data.ErrProductNotFound.Error(), http.StatusInternalServerError)
		return
	}
}
