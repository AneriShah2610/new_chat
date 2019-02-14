package directive

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
)

func EnumLogging(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	rc := graphql.GetResolverContext(ctx)
	log.Printf("enum logging: %v, %s, %T, %+v", rc.Path(), rc.Field.Name, obj, obj)
	return next(ctx)
}
