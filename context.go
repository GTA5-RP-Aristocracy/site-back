package siteback

import (
	"context"

	"github.com/google/uuid"
)

type (
	CtxKeyUserID struct{}
)

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, CtxKeyUserID{}, userID)
}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	if userID, ok := ctx.Value(CtxKeyUserID{}).(uuid.UUID); ok {
		return userID
	}
	return uuid.Nil
}
