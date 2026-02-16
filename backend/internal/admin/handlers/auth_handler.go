package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sdv77/Linguture1/internal/admin/models"
	"github.com/sdv77/Linguture1/internal/admin/service"
)

type AdminAuthHandler struct {
	adminService *service.AdminService
}

func NewAdminAuthHandler(adminService *service.AdminService) *AdminAuthHandler {
	return &AdminAuthHandler{adminService: adminService}
}

func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.AdminLoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	tokenResp, err := h.adminService.Login(input)
	if err != nil {
		log.Printf("Ошибка входа администратора: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResp)
}
