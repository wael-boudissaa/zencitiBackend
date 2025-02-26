package feedback

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wael-boudissaa/marquinoBackend/services/auth"
	"github.com/wael-boudissaa/marquinoBackend/types"
	"github.com/wael-boudissaa/marquinoBackend/utils"
)

type Handler struct {
	store types.FeedBackStore
}

func NewHanlder(store types.FeedBackStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/feedback", h.getAllFeedBack).Methods("GET")
	router.HandleFunc("/feedback", h.createFeedBack).Methods("POST")
}

func (h *Handler) getAllFeedBack(w http.ResponseWriter, r *http.Request) {
	feedBacks, err := h.store.GetAllFeedBack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, http.StatusOK, feedBacks)
}

func (h *Handler) createFeedBack(w http.ResponseWriter, r *http.Request) {

	comment := r.FormValue("comment")

	idFeedBack, err := auth.CreateAnId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idCustomer := mux.Vars(r)["idCustomer"]
	err = h.store.CreateFeedBack(idFeedBack, idCustomer, comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, http.StatusCreated, nil)
}
