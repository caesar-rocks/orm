package orm

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type BaseModel struct {
	bun.BaseModel
	CreatedAt time.Time
	UpdatedAt time.Time
}

var _ bun.BeforeAppendModelHook = (*BaseModel)(nil)

func (m *BaseModel) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
