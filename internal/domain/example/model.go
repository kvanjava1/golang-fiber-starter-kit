package example

import (
	"time"

	"github.com/uptrace/bun"
)

type Crud struct {
	bun.BaseModel `bun:"table:cruds,alias:c"`

	ID          string    `bun:"id,pk" json:"id"`
	Title       string    `bun:"title,notnull" json:"title"`
	Description string    `bun:"description" json:"description"`
	Status      string    `bun:"status,notnull,default:'active'" json:"status"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}
