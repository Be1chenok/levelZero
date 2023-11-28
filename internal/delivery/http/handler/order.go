package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/gorilla/mux"
)

func (h Handler) FindOrderByUID(w http.ResponseWriter, r *http.Request) {
	ctx, cancle := context.WithTimeout(r.Context(), h.conf.Server.RequestTime)
	defer cancle()

	vars := mux.Vars(r)
	uid := vars["uid"]

	order, err := h.service.Order.FindByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, domain.NothingFound) {
			writeJsonErrorResponse(w, http.StatusBadRequest, domain.NothingFound)
			return
		}
		writeJsonErrorResponse(w, http.StatusInternalServerError, SomethingWentWrong)
	}

	writeJsonResponse(w, http.StatusOK, order)
}
