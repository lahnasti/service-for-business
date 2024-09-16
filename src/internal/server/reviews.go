package server

import (
	"net/http"
	"strconv"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetReviewsHandler(ctx *gin.Context) {
	tenderIDStr := ctx.Param("tenderID")
	username := ctx.Query("username")
	organizationIDStr := ctx.Query("organizationId")

	tenderID, err := strconv.Atoi(tenderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	organizationID, err := strconv.Atoi(organizationIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	reviews, err := s.Db.GetReviewsByAuthorAndTender(tenderID, username, organizationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}

func (s *Server) AddFeedbackHandler(ctx *gin.Context) {
	var review models.Review
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Извлекаем данные из тела запроса
	bidID := review.BidID
	username := review.Username // Предполагается, что `username` также передается в теле запроса

	// Получаем tenderID по bidID
	tenderID, err := s.Db.GetTenderIDByBidID(bidID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tender ID"})
		return
	}

	// Проверяем права доступа пользователя
	hasPermission, err := s.Db.CheckUserPermissionForTender(tenderID, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user permission"})
		return
	}
	if !hasPermission {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to leave feedback"})
		return
	}

	// Добавляем отзыв в базу данных
	if err := s.Db.AddFeedback(review, username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add feedback"})
		return
	}

	ctx.JSON(http.StatusOK, review)
}
