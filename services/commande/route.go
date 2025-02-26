package commande

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/marquinoBackend/services/auth"
	"github.com/wael-boudissaa/marquinoBackend/types"
	"github.com/wael-boudissaa/marquinoBackend/utils"
)

type Handler struct {
	s types.CommandeStore
}

func NewHandler(s types.CommandeStore) *Handler {
	return &Handler{s: s}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/commande", h.getAllCommandes).Methods("GET")
	router.HandleFunc("/commande/{id}", h.CreateCommande).Methods("POST")
}

func (h *Handler) getAllCommandes(w http.ResponseWriter, r *http.Request) {
	commandes, err := h.s.GetAllCommandes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, http.StatusOK, commandes)
}

func (h *Handler) CreateCommande(w http.ResponseWriter, r *http.Request) {
	var idCustomer = mux.Vars(r)["id"]
	var products []types.ProductBought
	err := utils.ParseJsonList(r, &products)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error parsing JSON: %v", err))
		return
	}
	idCommande, err := auth.CreateAnId()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error while creating commande id"))
		return
	}
	var price int
	for _, product := range products {
		price += product.Price
	}

	err = h.s.CreateCommande(idCommande, idCustomer, price)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error while  ddd creating commande"))
		return
	}
	for _, product := range products {
		_, err := h.s.InsertProductINCommande(product, idCommande)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error while creating commande produtct"))
			return
		}
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "commande created"})

}
