package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/arnav0418/ecomm_go/ecomm-api/server"
	"github.com/arnav0418/ecomm_go/ecomm-api/storer"
	"github.com/go-chi/chi"
)

type handler struct {
	ctx context.Context
	server *server.Server
}

func NewHandler(server *server.Server) *handler {
	return &handler{ctx: context.Background(), server: server}
}

func (h *handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var p ProductReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error creating product", http.StatusBadRequest)
		return
	}
	
	product, err := h.server.CreateProduct(h.ctx, toStorerProduct(&p))
	if err != nil {
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id,10,64)
	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return 
	}

	product, err := h.server.GetProduct(h.ctx, i)
	if err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) listProducts(w http.ResponseWriter, r *http.Request) {
	products,  err := h.server.ListProducts(h.ctx)	
	if err != nil {
		http.Error(w, "error listing products", http.StatusInternalServerError)
		return
	}

	var res []ProductRes
	for _, p := range products {
		res = append(res, toProductRes(&p))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func toProductRes(p *storer.Product) ProductRes {
	return ProductRes{
		ID:           p.ID,	
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,	
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func toStorerProduct(p *ProductReq) *storer.Product {
	return &storer.Product{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}