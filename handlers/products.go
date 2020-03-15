package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dzonint/go-microservice/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp, err := data.GetProducts()
	if err != nil {
		http.Error(rw, "Unable to get products", http.StatusInternalServerError)
		return
	}
	err = lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall JSON", http.StatusInternalServerError)
		return
	}
}

func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	product, err := data.GetProduct(id)
	if err != nil {
		http.Error(rw, "Unable to get products", http.StatusInternalServerError)
		return
	}
	err = product.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshall JSON", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(KeyProduct{}).(*data.Product)
	err := data.AddProduct(product)
	if err == data.ErrFailedToOpenDB {
		http.Error(rw, "Unable to access database", http.StatusInternalServerError)
		return
	} else if err == data.ErrFailedToUpdateDB {
		http.Error(rw, "Unable to add product", http.StatusBadRequest)
		return
	}
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	} else if err == data.ErrFailedToOpenDB {
		http.Error(rw, "Unable to access database", http.StatusInternalServerError)
		return
	} else if err == data.ErrFailedToUpdateDB {
		http.Error(rw, "Unable to update database", http.StatusBadRequest)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := &data.Product{}
		err := product.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		err = product.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(rw, fmt.Sprintf("Error validating product: &s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		nextHandler.ServeHTTP(rw, r)
	})
}
