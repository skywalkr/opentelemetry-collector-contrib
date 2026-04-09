// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottlfuncs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"

import (
	"context"
	"errors"
	"net"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

type ParseIPArguments[K any] struct {
	IP ottl.StringGetter[K]
}

func NewParseIPFactory[K any]() ottl.Factory[K] {
	return ottl.NewFactory("ParseIP", &ParseIPArguments[K]{}, createParseIPFunction[K])
}

func createParseIPFunction[K any](_ ottl.FunctionContext, oArgs ottl.Arguments) (ottl.ExprFunc[K], error) {
	args, ok := oArgs.(*ParseIPArguments[K])
	if !ok {
		return nil, errors.New("ParseIpFactory args must be of type *ParseIPArguments[K]")
	}

	return parseIPFunc(args.IP), nil
}

// parseIPFunc returns a `pcommon.Map` struct that is a result of parsing the target string as an IP address
func parseIPFunc[K any](target ottl.StringGetter[K]) ottl.ExprFunc[K] {
	return func(ctx context.Context, tCtx K) (any, error) {
		t, err := target.Get(ctx, tCtx)
		if err != nil {
			return nil, err
		}
		if t == "" {
			return nil, errors.New("cannot parse from empty target")
		}

		result := net.ParseIP(t)

		if result == nil {
			return nil, errors.New("could not parse IP address from target")
		}

		var resultMap = pcommon.NewMap()
		resultMap.PutStr("Address", result.String())
		resultMap.PutBool("IsMulticast", result.IsMulticast())
		resultMap.PutBool("IsLoopback", result.IsLoopback())
		resultMap.PutBool("IsPrivate", result.IsPrivate())

		return resultMap, nil
	}
}
