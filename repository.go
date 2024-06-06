package orm

import (
	"context"
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

func (r *Repository[M]) FindOneBy(ctx context.Context, field string, value any) (*M, error) {
	var item *M = new(M)

	err := r.NewSelect().Model((*M)(nil)).Where(field+" = ?", value).Scan(ctx, item)
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
