package postgres

import (
	"github.com/jackc/pgtype"
)

type TogoTask struct {
	ID          int64
	Name        string
	Description string
	CreatedOn   pgtype.Timestamptz
	CompletedOn pgtype.Timestamptz
	DueDate     pgtype.Timestamptz
}
