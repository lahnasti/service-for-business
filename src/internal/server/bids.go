package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetBidsByUserHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	bids, err := s.Db.GetBidsByUser(username)
	if err != nil {
		s.log.Error().Err(err).Msg("Invalid username")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid params"})
		return
	}
	ctx.JSON(http.StatusOK, bids)
}

func (s *Server) GetBidsForTenderHandler(ctx *gin.Context) {
	idStr := ctx.Param("tenderID")
	s.log.Info().Msgf("Received tenderId: %s", idStr)

	tenderID, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}
	bids, err := s.Db.GetBidsForTender(tenderID)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get bids for tender")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get bids for tender", "error": err})
		return
	}
	// Если предложений нет
	if len(bids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No bids found for this tender"})
		return
	}
	ctx.JSON(http.StatusOK, bids)
}
func (s *Server) CreateBidHandler(ctx *gin.Context) {
	var bid models.Bid
	if err := ctx.ShouldBindJSON(&bid); err != nil {
		s.log.Error().Err(err).Msg("Invalid JSON payload")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON payload", "error": err.Error()})
		return
	}
	if err := s.Valid.Struct(bid); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creatorUsername := bid.CreatorUsername
	id, err := s.Db.CreateBid(bid, creatorUsername)
	if err != nil {
		if err.Error() == fmt.Sprintf("user %s does not have permission to create bid", bid.CreatorUsername) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "user does not have permission to create bid"})
			return
		}
		s.log.Error().Err(err).Msg("Failed to add bid")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to add bid", "error": err})
		return
	}
	ctx.JSON(http.StatusOK, id)
}

func (s *Server) SetBidStatusHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
	}
	var requestBody struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if requestBody.Status != "PUBLISHED" && requestBody.Status != "CANCELED" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	err = s.Db.SetBidStatus(id, requestBody.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Bid status updated successfully"})
}

func (s *Server) EditBidHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
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

	query, err := s.Db.EditBid(id, requestBody.Name, requestBody.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, query)
}

func (s *Server) RollbackBidHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}
	versionStr := ctx.Param("version")
	log.Println("version:", versionStr)
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
		return
	}
	updateBid, err := s.Db.RollbackBid(id, version)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updateBid)
}

func (s *Server) SubmitDecisionHandler(ctx *gin.Context) {
	bidIDStr := ctx.Param("id")
	bidID, err := strconv.Atoi(bidIDStr)
	if err != nil {
		s.log.Error().Err(err).Msg("Invalid bid ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid bid ID", "error": err.Error()})
		return
	}

	// Получаем имя пользователя из тела запроса
	var requestBody struct {
		Username string `json:"username" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		s.log.Error().Err(err).Msg("Invalid JSON payload")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON payload", "error": err.Error()})
		return
	}

	// Проверка прав пользователя и согласование предложения
	err = s.Db.SubmitDecision(bidID, requestBody.Username)
	if err != nil {
		// Обработка ошибок связанных с проверкой статуса
		if err.Error() == "bid must be in PUBLISHED status to submit decision" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err.Error() == "user does not have permission to approve or decline this bid" {
			s.log.Error().Err(err).Msg("User does not have permission")
			ctx.JSON(http.StatusForbidden, gin.H{"message": "User does not have permission"})
			return
		}
		s.log.Error().Err(err).Msg("Failed to approve bid")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to approve bid", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Bid approved and tender closed"})
}

func (s *Server) SubmitDeclinedHandler(ctx *gin.Context) {
	// Получаем идентификатор предложения (bidID) из параметров URL
	bidIDStr := ctx.Param("bidID")
	bidID, err := strconv.Atoi(bidIDStr)
	if err != nil {
		s.log.Error().Err(err).Msg("Invalid bid ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid bid ID", "error": err.Error()})
		return
	}

	// Получаем имя пользователя из тела запроса (например, переданное в JSON)
	var requestBody struct {
		Username string `json:"username" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		s.log.Error().Err(err).Msg("Invalid JSON payload")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON payload", "error": err.Error()})
		return
	}

	// Проверка прав пользователя и отклонение предложения
	err = s.Db.DeclineDecision(bidID, requestBody.Username)
	if err != nil {
		// Обработка ошибок связанных с проверкой статуса
		if err.Error() == "bid must be in PUBLISHED status to submit decision" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "user does not have permission to approve or decline this bid" {
			s.log.Error().Err(err).Msg("User does not have permission")
			ctx.JSON(http.StatusForbidden, gin.H{"message": "User does not have permission"})
			return
		}
		s.log.Error().Err(err).Msg("Failed to decline bid")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to decline bid", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Bid declined"})
}
