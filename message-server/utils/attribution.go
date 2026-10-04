package utils

import (
	"context"
	"fmt"
	"net/http"
)

const key string = "userAttribution"
const dfault string = "Unknown"

func SetAttribution(ctx context.Context, attribution string) context.Context {
	if attribution == "" {
		attribution = dfault
	}
	return context.WithValue(ctx, key, attribution)
}

func GetAttribution(ctx context.Context) (s string) {
	defer func() {
		if r := recover(); r != nil {
			s = dfault	
		}
	}()

	s = ctx.Value(key).(string)
	if s == "" {
		return dfault
	}
	return s
}

func AttributeLog(msg string, r *http.Request) string {
	return fmt.Sprintf(" %s for user %s at addr %s", msg, GetAttribution(r.Context()), r.RemoteAddr)
}
