package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/knbr13/company-service-go/internal/repositories"
	"github.com/knbr13/company-service-go/internal/services"
	"github.com/knbr13/company-service-go/internal/validator"
	"github.com/knbr13/company-service-go/pkg/util"
)

type CompanyHandler struct {
	Services *services.Services
}

func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req repositories.Company
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrJsonResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	compID, err := h.Services.Companies.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, repositories.ErrCompanyNameAlreadyExists) {
			util.ErrJsonResponse(w, http.StatusConflict, fmt.Sprintf("name: %s", err.Error()))
			return
		}
		if e, ok := err.(validator.ValidationError); ok {
			util.ErrJsonResponse(w, http.StatusBadRequest, e.Error())
			return
		}
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error creating company: %s\n", err.Error())
		return
	}

	util.JsonResponse(w, http.StatusCreated, map[string]any{
		"company_id": compID,
	})
}

func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	compID := chi.URLParam(r, "id")
	comp, err := h.Services.Companies.Get(ctx, compID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.ErrJsonResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error getting company: %s\n", err.Error())
		return
	}

	util.JsonResponse(w, http.StatusOK, comp)
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	compID := chi.URLParam(r, "id")
	comp, err := h.Services.Companies.Get(ctx, compID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.ErrJsonResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error getting company: %s\n", err.Error())
		return
	}

	var req repositories.Company
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrJsonResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if req.Name != "" {
		comp.Name = req.Name
	}
	if req.Description != "" {
		comp.Description = req.Description
	}
	if req.Type != "" {
		comp.Type = req.Type
	}
	if req.AmountOfEmployees != nil {
		comp.AmountOfEmployees = req.AmountOfEmployees
	}
	if req.Registered != nil {
		comp.Registered = req.Registered
	}

	err = h.Services.Companies.UpdateCompany(ctx, comp)
	if err != nil {
		if e, ok := err.(validator.ValidationError); ok {
			util.ErrJsonResponse(w, http.StatusBadRequest, e.Error())
			return
		}
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error updating company: %s\n", err.Error())
		return
	}

	util.JsonResponse(w, http.StatusOK, comp)
}

func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	compID := chi.URLParam(r, "id")
	err := h.Services.Companies.Delete(ctx, compID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.ErrJsonResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error deleting company: %s\n", err.Error())
		return
	}

	util.JsonResponse(w, http.StatusNoContent, nil)
}
