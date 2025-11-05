package ctx

import (
	"context"
	"errors"
)

const key string = "userId"

func SetUserIdCtx(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, key, val)
}

func GetUserIdCtx(ctx context.Context) (string, error) {
	id, ok := ctx.Value(key).(string)
	if !ok {
		return "", errors.New("failed getting the user id in the contex")
	}
	return id, nil
}
