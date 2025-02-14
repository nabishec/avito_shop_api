package token

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var signingKey = []byte(os.Getenv("SIGNING_KEY"))

func CreateJWT(userID uuid.UUID) (token string, err error) {
	op := "internal.pkg.token.CreateJWT()"

	//jti = uuid.New().String()
	sub := userID.String()
	exp := time.Now().Unix() + 10800 // 3 hour

	claims := jwt.StandardClaims{
		//Id:        jti,
		Subject:   sub,
		ExpiresAt: exp,
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodES512, claims).SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	return
}

func CheckJWT(tokenString string) (userID string, err error) {
	op := "internal.pkg.token.CheckJWT()"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})
	if err != nil || !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {

				return "", fmt.Errorf("%s:%s", op, "that's not even a token")

			} else if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
				// Token is either expired
				return "", fmt.Errorf("%s:%s", op, "timing is everything")
			}
			return "", fmt.Errorf("%s,%s", op, "invalid token")
		}
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		return claims.Subject, nil
	} else {
		return "", fmt.Errorf("%s:%s", op, "failed conversion of jwt claims")
	}
}
