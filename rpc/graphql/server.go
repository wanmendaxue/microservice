package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var RequireAuthenticatedError = errors.New("operation requires authenticated")

func NewGraphqlServer(es graphql.ExecutableSchema) *handler.Server {
	h := handler.NewDefaultServer(es)
	h.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		formattedErr := e
		if formattedErr != nil {
			if formattedErr == RequireAuthenticatedError {
				formattedErr = &gqlerror.Error{
					Message: e.Error(),
					Extensions: map[string]interface{}{
						"code": 401,
					},
				}
			} else if _, ok := e.(*gqlerror.Error); ok {
				// if err is GQL error object, then do nothing
				logrus.WithFields(logrus.Fields{"error": e, "type": "gql-error"}).Debugf("service response gql error: %v", e)
			} else if ex, ok := e.(interface{ Business() (uint32, string) }); ok {
				c, m := ex.Business()
				logrus.WithFields(logrus.Fields{"error": e, "type": "business"}).Debugf("service response business error: [%d] %s", c, m)
				formattedErr = &gqlerror.Error{
					Message: e.Error(),
					Extensions: map[string]interface{}{
						"code": c,
						"msg":  m,
					},
				}
			} else if ex, ok := e.(interface{ Demand() string }); ok {
				logrus.WithFields(logrus.Fields{"error": e, "type": "demand"}).Debugf("service response demand error: %s", ex.Demand())
			} else {
				logrus.WithFields(logrus.Fields{"error": e, "type": "general"}).Debugf("service response general error: %+v", e)
			}
		}
		return graphql.DefaultErrorPresenter(ctx, formattedErr)
	})

	return h
}
