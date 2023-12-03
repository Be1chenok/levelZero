package handler

import (
	"context"
	"errors"
	"net/http"
	"text/template"

	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/gorilla/mux"
)

const (
	homeHtml         = "../../web/template/home.html"
	orderHtml        = "../../web/template/order.html"
	nothingFoundHtml = "../../web/template/nothingFound.html"
)

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("orderUID")
	http.Redirect(w, r, "/order/"+orderUID, http.StatusFound)
}

func (h Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(homeHtml)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	w.Header().Set(contentType, textHtml)

	if err = tmpl.Execute(w, nil); err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}
}

func (h Handler) FindOrderByUID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.conf.Server.RequestTime)
	defer cancel()

	vars := mux.Vars(r)
	uid := vars["uid"]

	order, err := h.service.Order.FindByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, domain.ErrNothingFound) {
			h.NothingFound(w, r)
			return
		}
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	tmpl, err := template.ParseFiles(orderHtml)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	w.Header().Set(contentType, textHtml)
	w.WriteHeader(http.StatusOK)
	if err = tmpl.Execute(w, order); err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}
}

func (h Handler) NothingFound(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(nothingFoundHtml)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}

	w.Header().Set(contentType, textHtml)
	w.WriteHeader(http.StatusBadRequest)
	if err = tmpl.Execute(w, nil); err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrSomethingWentWrong)
		return
	}
}
