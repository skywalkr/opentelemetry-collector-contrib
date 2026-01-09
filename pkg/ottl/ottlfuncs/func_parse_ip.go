package ottlfuncs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"

import (
	"context"
	"errors"
	"net"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

type ParseIpArguments[K any] struct {
	Ip ottl.StringGetter[K]
}

func NewParseIpFactory[K any]() ottl.Factory[K] {
	return ottl.NewFactory("ParseIp", &ParseIpArguments[K]{}, createParseIpFunction[K])
}

func createParseIpFunction[K any](_ ottl.FunctionContext, oArgs ottl.Arguments) (ottl.ExprFunc[K], error) {
	args, ok := oArgs.(*ParseIpArguments[K])

	if !ok {
		return nil, errors.New("ParseIpFactory args must be of type *ParseIpArguments[K]")
	}

	return parseIpFunc(args.Ip)
}

func parseIpFunc[K any](target ottl.StringGetter[K]) (ottl.ExprFunc[K], error) {
	return func(ctx context.Context, tCtx K) (any, error) {
		targetValue, err := target.Get(ctx, tCtx)
		if err != nil {
			return nil, err
		}

		ip := net.ParseIP(targetValue)

		return ip, nil
	}, nil
}
