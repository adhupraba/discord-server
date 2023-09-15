package utils

import (
	"context"
	"errors"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	"github.com/adhupraba/discord-server/lib"
)

func GetUserFromClerk(ctx context.Context) (user *clerk.User, errCode int, err error) {
	sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)

	if !ok {
		return nil, http.StatusUnauthorized, errors.New("Unauthorized")
	}

	user, err = lib.ClerkClient.Users().Read(sessClaims.Subject)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, 0, nil
}
