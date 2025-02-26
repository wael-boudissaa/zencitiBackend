package product

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/marquinoBackend/services/auth"
	"github.com/wael-boudissaa/marquinoBackend/types"
	"github.com/wael-boudissaa/marquinoBackend/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// // customer
	router.HandleFunc("/products", h.GetAllProduct).Methods("GET")
	router.HandleFunc("/product/{id}", h.GetOneProduct).Methods("GET")
	router.HandleFunc("/products/{idCategorie}", h.GetProductByCategorie).Methods("GET")

	// router.HandleFunc("/products", h.store).Methods("GET")
	// router.HandleFunc("/product/{id}", h.getProduct).Methods("GET")
	// router.HandleFunc("/product", h.createProduct).Methods("POST")
	// router.HandleFunc("/product/{id}", h.updateProduct).Methods("PUT")
	// router.HandleFunc("/product/{id}", h.deleteProduct).Methods("DELETE")
	// //admin
	router.HandleFunc("/product", h.CreateProduct).Methods("POST")
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// product,
	var product types.ProductCreate
	err := utils.ParseJson(r, &product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	idProduct, err := auth.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error while creating product"))
		return
	}
	err = h.store.CreateProduct(product, idProduct)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error while creating product"))
		return
	}
	utils.WriteJson(w, http.StatusOK, "product created")
}

func (h *Handler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetAllProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, products)
}

func (h *Handler) GetOneProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	product, err := h.store.GetProductById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error while getting product"))
		return
	}
	utils.WriteJson(w, http.StatusOK, *product)
}

func (h *Handler) GetProductByCategorie(w http.ResponseWriter, r *http.Request) {
	idCategorie := mux.Vars(r)["idCategorie"]
	products, err := h.store.GetProductByCategorie(idCategorie)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error while getting products"))
		return
	}
	utils.WriteJson(w, http.StatusOK, products)
}
