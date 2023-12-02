package handler

import (
	"context"
	"errors"
	"net/http"
	"text/template"

	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/gorilla/mux"
)

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("orderUID")
	http.Redirect(w, r, "/order/"+orderUID, http.StatusFound)
}

func (h Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../web/home.html")
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, SomethingWentWrong)
		return
	}

	w.Header().Set(contentType, textHtml)

	if err = tmpl.Execute(w, nil); err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, SomethingWentWrong)
		return
	}
}

func (h Handler) FindOrderByUID(w http.ResponseWriter, r *http.Request) {
	ctx, cancle := context.WithTimeout(r.Context(), h.conf.Server.RequestTime)
	defer cancle()

	vars := mux.Vars(r)
	uid := vars["uid"]

	order, err := h.service.Order.FindByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, domain.NothingFound) {
			h.NothingFound(w, r)
			return
		}
		writeJsonErrorResponse(w, http.StatusInternalServerError, SomethingWentWrong)
		return
	}

	tmpl, err := template.ParseFiles("../../web/order.html")
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, SomethingWentWrong)
		return
	}

	w.Header().Set(contentType, textHtml)
	w.WriteHeader(http.StatusOK)
	if err = tmpl.Execute(w, order); err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, SomethingWentWrong)
		return
	}
}

func (h Handler) NothingFound(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../web/nothingFound.html")
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, SomethingWentWrong)
		return
	}

	w.Header().Set(contentType, textHtml)
	w.WriteHeader(http.StatusBadRequest)
	if err = tmpl.Execute(w, nil); err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, SomethingWentWrong)
		return
	}
}
