package core

import (
	"context"
	"reflect"

	"github.com/nkbhasker/go-auth-starter/internal/uid"
)

type Identity interface {
	ID() string
	UserID() uid.Identifier
}

type indentity struct {
	jti    string
	userId uid.Identifier
}

type identityContextKey struct{}

func NewIdentity(jti, sub string) (Identity, error) {
	userId, err := uid.FromIdString(sub)
	if err != nil {
		return nil, err
	}

	return &indentity{jti: jti, userId: userId}, nil
}

func (u *indentity) ID() string {
	return u.jti
}

func (u *indentity) UserID() uid.Identifier {
	return u.userId
}

func IdentityFromContext(ctx context.Context) Identity {
	ctxValue, ok := ctx.Value(identityContextKey{}).(Identity)
	if !ok {
		return &indentity{}
	}

	return ctxValue
}

func IdentityToContext(ctx context.Context, identity Identity) context.Context {
	if reflect.ValueOf(identity).IsNil() {
		return ctx
	}

	return context.WithValue(ctx, identityContextKey{}, identity)
}
