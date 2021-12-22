package originip

import "context"

type ipCtxKey struct{}

type OriginIP struct {
	IP string
}

var ipKey = ipCtxKey{}

func NewContext(ctx context.Context, oip *OriginIP) context.Context {
	return context.WithValue(ctx, ipKey, oip)
}

func FromCtx(ctx context.Context) *OriginIP {
	oip, ok := ctx.Value(ipKey).(*OriginIP)
	if !ok {
		// FIXME: check what to do in case of error obtaining the Origin IP
		return &OriginIP{IP: "127.0.0.1"}
	}
	return oip
}
