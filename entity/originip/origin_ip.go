package originip

import (
	"context"

	"github.com/labstack/gommon/log"
)

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
		log.Warn("incoming request IP not found")
		return &OriginIP{IP: "192.168.0.1"}
	}

	log.Debugf("incoming request IP: %v", oip.IP)
	return oip
}
