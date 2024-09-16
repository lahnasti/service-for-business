package routes

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/server"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(s *server.Server) *gin.Engine {
	r := gin.Default()

	pingGroup := r.Group("/api")
	{
		pingGroup.GET("/ping", s.PingHandler)
	}

	tenderGroup := r.Group("/api/tenders")
	{
		tenderGroup.GET("/", s.GetAllTendersHandler)
		tenderGroup.GET("/my", s.GetTendersByUser)
		tenderGroup.POST("/new", s.CreateTenderHandler)
		tenderGroup.PATCH("/status/:id", s.SetTenderStatusHandler)
		tenderGroup.PATCH("/:id/edit", s.EditTenderHandler)
		tenderGroup.PUT("/:tenderID/rollback/:version", s.RollbackTenderHandler)
	}

	bidsGroup := r.Group("/api/bids")
	{
		bidsGroup.GET("/:tenderID/list", s.GetBidsForTenderHandler)
		bidsGroup.GET("/my", s.GetBidsByUserHandler)
		bidsGroup.POST("/new", s.CreateBidHandler)
		bidsGroup.PATCH("/status/:id", s.SetBidStatusHandler)
		bidsGroup.PATCH("/:id/edit", s.EditBidHandler)
		bidsGroup.PUT("/:bidID/rollback/:version", s.RollbackBidHandler)
		bidsGroup.PATCH("/:id/submit_decision", s.SubmitDecisionHandler)
		bidsGroup.PATCH("/:id/decline_decision", s.SubmitDeclinedHandler)

		// GET /api/bids/1/reviews?authorUsername=user2&organizationId=1
		//}
		return r
	}
}
