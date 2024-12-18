package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.60

import (
	"context"
	"fmt"

	"github.com/skeetcha/pf2ql/graph/model"
)

// FindSource is the resolver for the findSource field.
func (r *queryResolver) FindSource(ctx context.Context, id *string) (*model.Source, error) {
	rows, err := r.Db.Query("select * from sources where id=" + *id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		val, err := GetData(rows)

		if err != nil {
			return nil, err
		}

		return val, nil
	}

	return nil, fmt.Errorf("id %s not found", *id)
}

// FindSources is the resolver for the findSources field.
func (r *queryResolver) FindSources(ctx context.Context, filter *model.SourceFilter) ([]*model.Source, error) {
	str := "select * from sources where "
	val, err := GetFilterString(filter)

	if err != nil {
		return nil, err
	}

	str += val

	rows, err := r.Db.Query(str)

	if err != nil {
		return nil, err
	}

	sources := []*model.Source{}

	for rows.Next() {
		data, err := GetData(rows)

		if err != nil {
			return nil, err
		}

		sources = append(sources, data)
	}

	return sources, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
