package server

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/models"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog"
)

// TODO: validation add
type TendersRepo interface {
	GetAllTenders(string) ([]models.Tender, error)
	GetTendersByUser(string) ([]models.Tender, error)
	CreateTender(models.Tender) (models.Tender, error)
	SetTenderStatus(int, string) error
	EditTender(int, string, string) (models.Tender, error)
	RollbackTender(int, int) (models.Tender, error)
}

type BidsRepo interface {
	GetBidsByUser(string) ([]models.Bid, error)
	GetBidsForTender(int) ([]models.Bid, error)
	CreateBid(models.Bid, string) (models.Bid, error)
	SetBidStatus(int, string) error
	EditBid(int, string, string) (models.Bid, error)
	RollbackBid(int, int) (models.Bid, error)
	SubmitDecision(int, string) error
	DeclineDecision(int, string) error
	CheckUserPermissionForBid(int, string) (bool, error)
}

type Repository interface {
	TendersRepo
	BidsRepo
}

type Server struct {
	Db    Repository
	log   zerolog.Logger
	Valid *validator.Validate
}

func New(ctx context.Context, db Repository, zlog *zerolog.Logger) *Server {
	validate := validator.New()
	return &Server{
		Db:    db,
		log:   *zlog,
		Valid: validate,
	}
}
