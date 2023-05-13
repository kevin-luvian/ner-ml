package repo

//go:generate mockgen -source=./init.go -destination=mock_user_repo.go -package=user

import (
	"github.com/kevin-luvian/gomodify/pkg/db"
)

type Repo struct {
	db *db.DB
}

func New(db *db.DB) *Repo {
	return &Repo{
		db: db,
	}
}
