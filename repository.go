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

// UpdateOneWhere updates a single record in the database based on multiple field conditions.
// It accepts a context, a variadic list of arguments representing field-value pairs, and the item to be updated.
// The args should be passed in pairs like ("field1", value1, "field2", value2).
func (r *Repository[M]) UpdateOneWhere(ctx context.Context, item *M, args ...any) error {
	// Ensure args are in key-value pairs
	if len(args)%2 != 0 {
		return errors.New("arguments must be key-value pairs")
	}

	// Start building the update query
	query := r.NewUpdate().Model(item)

	// Iterate through the arguments to build the WHERE clause
	for i := 0; i < len(args); i += 2 {
		// Ensure the key is a string representing the field name
		field, ok := args[i].(string)
		if !ok {
			return errors.New("field names must be strings")
		}
		// Get the value associated with the field
		value := args[i+1]
		// Add the condition to the query
		query = query.Where(field+" = ?", value)
	}

	// Execute the update query
	_, err := query.Exec(ctx)

	return err
}

// DeleteOneWhere deletes a single record from the database based on multiple field conditions.
// It accepts a context and a variadic list of arguments representing field-value pairs.
// The args should be passed in pairs like ("field1", value1, "field2", value2).
func (r *Repository[M]) DeleteOneWhere(ctx context.Context, args ...any) error {
	// Ensure args are in key-value pairs
	if len(args)%2 != 0 {
		return errors.New("arguments must be key-value pairs")
	}

	// Initialize the model item
	var item *M = new(M)

	// Start building the delete query
	query := r.NewDelete().Model(item)

	// Iterate through the arguments to build the WHERE clause
	for i := 0; i < len(args); i += 2 {
		// Ensure the key is a string representing the field name
		field, ok := args[i].(string)
		if !ok {
			return errors.New("field names must be strings")
		}
		// Get the value associated with the field
		value := args[i+1]
		// Add the condition to the query
		query = query.Where(field+" = ?", value)
	}

	// Execute the delete query
	_, err := query.Exec(ctx)

	return err
}
