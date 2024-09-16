package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/models"
	"github.com/gin-gonic/gin"
)

// возвращает список тендеров с фильтрацией по типу услуг
func (s *Server) GetAllTendersHandler(ctx *gin.Context) {
	serviceType := ctx.Query("serviceType")

	tenders, err := s.Db.GetAllTenders(serviceType)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get tenders")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch tenders", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List of tenders", "tenders": tenders})
}

func (s *Server) GetTendersByUser(ctx *gin.Context) {
	username := ctx.Query("username")
	tenders, err := s.Db.GetTendersByUser(username)
	if err != nil {
		s.log.Error().Err(err).Msg("Invalid username")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params"})
		return
	}
	ctx.JSON(http.StatusOK, tenders)

}

func (s *Server) CreateTenderHandler(ctx *gin.Context) {
	var tender models.Tender
	if err := ctx.ShouldBindBodyWithJSON(&tender); err != nil {
		s.log.Error().Err(err).Msg("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	if err := s.Valid.Struct(tender); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tender, err := s.Db.CreateTender(tender)
	if err != nil {
		if err.Error() == fmt.Sprintf("user %s is not responsible for organization %d", tender.CreatorUsername, tender.OrganizationID) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "User is not responsible for this organization"})
			return
		}
		s.log.Error().Err(err).Msg("Failed to add tender")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add tender", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Tender created successfully", "tender": tender})
}
func (s *Server) SetTenderStatusHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	// Извлечение нового статуса из тела запроса
	var requestBody struct {
		Status string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if requestBody.Status != "PUBLISHED" && requestBody.Status != "CLOSED" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	err = s.Db.SetTenderStatus(id, requestBody.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tender status updated successfully"})
}

func (s *Server) EditTenderHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := s.Valid.Struct(requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	query, err := s.Db.EditTender(id, requestBody.Name, requestBody.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, query)
}

func (s *Server) RollbackTenderHandler(ctx *gin.Context) {
	idStr := ctx.Param("tenderID")
	log.Println("tenderID:", idStr) // Логируем значение
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}
	versionStr := ctx.Param("version")
	log.Println("version:", versionStr) // Логируем значение
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
		return
	}
	updatedTender, err := s.Db.RollbackTender(id, version)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedTender)
}
