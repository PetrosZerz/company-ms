package adapters

import (
	"company-ms/internal"
	"company-ms/internal/application"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type CompanyHandler struct {
	service *application.CompanyService
	logger  *zap.Logger
}

func NewCompanyHandler(service *application.CompanyService, logger *zap.Logger) *CompanyHandler {
	return &CompanyHandler{service, logger}
}

func (h *CompanyHandler) respondWithError(w http.ResponseWriter, code int, err error) {
	h.logger.Error("Responding with error:")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company internal.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err)
		return
	}
	company.ID = uuid.New().String()
	if err := h.service.Create(&company); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err)
		return
	}

	h.logger.Info("Company created")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(company); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *CompanyHandler) GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.GetAll()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	h.logger.Info("Companies retrieved")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(companies); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	company, err := h.service.GetByID(idStr)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, err)
		return
	}
	h.logger.Info("Company retrieved")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(company); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var company internal.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err)
		return
	}

	company.ID = id
	if err := h.service.Update(&company); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("Company updated")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(company); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if err := h.service.Delete(idStr); err != nil {
		h.respondWithError(w, http.StatusNotFound, err)
		return
	}
	h.logger.Info("Company deleted")
	w.WriteHeader(http.StatusNoContent)
}
