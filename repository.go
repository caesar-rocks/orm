package orm

import (
	"context"
)

type Model interface{}

type Application struct {
	ID   int    `json:"id" column:"id" primary:"true"`
	Name string `json:"name" column:"name"`
	Slug string `json:"slug" column:"slug"`
}

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

func (r *Repository[M]) FindOne(ctx context.Context, id int) (*M, error) {
	var item *M = new(M)

	err := r.NewSelect().Model((*M)(nil)).Where("id = ?", id).Scan(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *Repository[M]) FindOneBy(ctx context.Context, field string, value interface{}) (*M, error) {
	var item *M = new(M)

	err := r.NewSelect().Model((*M)(nil)).Where(field+" = ?", value).Scan(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *Repository[M]) UpdateOne(ctx context.Context, item M) error {
	_, err := r.NewUpdate().Model(item).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[M]) DeleteOne(ctx context.Context, id int) error {
	_, err := r.NewDelete().Model((*M)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
