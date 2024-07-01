package orm

import (
	"context"
	"errors"
)

type Model interface{}

type Repository[M Model] struct {
	*Database
}

func NewRepository[M Model]() Repository[M] {
	return Repository[M]{}
}

func (r *Repository[M]) Create(ctx context.Context, item *M) error {
	_, err := r.NewInsert().Model(item).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository[M]) FindAll(ctx context.Context) ([]M, error) {
	var items []M = make([]M, 0)

	err := r.NewSelect().Model((*M)(nil)).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FindOneBy finds a single record from the database based on multiple field conditions.
// It accepts a context and a variadic list of arguments representing field-value pairs.
// The args should be passed in pairs like ("field1", value1, "field2", value2).
func (r *Repository[M]) FindOneBy(ctx context.Context, args ...any) (*M, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("arguments must be key-value pairs")
	}

	var item *M = new(M)
	query := r.NewSelect().Model((*M)(nil))

	for i := 0; i < len(args); i += 2 {
		field, ok := args[i].(string)
		if !ok {
			return nil, errors.New("field names must be strings")
		}
		value := args[i+1]
		query = query.Where(field+" = ?", value)
	}

	err := query.Scan(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *Repository[M]) UpdateOneWhere(ctx context.Context, field string, value any, item *M) error {
	_, err := r.NewUpdate().Model(item).Where(field+" = ?", value).Exec(ctx)

	return err
}

func (r *Repository[M]) DeleteOneWhere(ctx context.Context, field string, value any) error {
	_, err := r.NewDelete().Model((*M)(nil)).Where(field+" = ?", value).Exec(ctx)

	return err
}
