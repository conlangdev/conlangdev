package conlangdev

import "context"

func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, "conlangdev_user", user)
}

func GetUserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value("conlangdev_user").(*User)
	return user
}
