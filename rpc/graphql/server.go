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

func NewGraphqlServer(es graphql.ExecutableSchema, prod bool) *handler.Server {
	h := handler.NewDefaultServer(es)
	h.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		if err != nil {
			inner := err.Unwrap()

			if inner != nil {
				if inner == RequireAuthenticatedError {
					err.Message = inner.Error()
					err.Extensions = map[string]interface{}{
						"code": 401,
						"msg": inner.Error(),
					}
					logrus.Debugf("unauthenticated: %+v", e)
				} else if ex, ok := inner.(interface{ Business() (uint32, string) }); ok {
					c, m := ex.Business()
					logrus.Debugf("business error : [%d] %s", c, m)
					err.Message = inner.Error()
					err.Extensions = map[string]interface{}{
						"code": c,
						"msg":  m,
					}
				}  else {
					logrus.Debugf("general error  :  %+v", e)

					if prod {
						msg := "service unavailable currently"
						err.Message = msg
						err.Extensions = map[string]interface{}{
							"code": 500,
							"msg": msg,
						}
					}
				}
			} else {
				logrus.Errorf("error          : %v", e)
			}
		}

		return err
	})

	return h
}
