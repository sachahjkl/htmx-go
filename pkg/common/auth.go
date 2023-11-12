package common

import (
	"fmt"
	"sachahjkl/htmx_go/pkg/model"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func UserFromJwt(db *gorm.DB, userJWT *jwt.Token) (*model.User, error) {

	if userJWT == nil {
		return nil, fmt.Errorf("user JWT is nil")
	}

	claims := userJWT.Claims.(jwt.MapClaims)
	userIdStr := claims[USER_CLAIM_KEY]
	if userIdStr == nil {
		return nil, fmt.Errorf("user id not in claims")
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%v", userIdStr), 10, 32)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse user id from claims")
	}

	return model.GetUser(db, uint(userId))

}
