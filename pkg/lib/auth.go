package lib

import (
	"errors"
	"strings"
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type ReqHeaderAuth struct {
	Authorization string `reqHeader:"authorization"`
}

var Claim Claims

type Claims struct {
	jwt.StandardClaims
	IsAdmin bool
}

// ClaimsJWT func
func ClaimsJWT(accesToken *string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(*accesToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	timeNow := time.Now().Unix()
	timeSessions := int64(claims["exp"].(float64))
	if timeSessions < timeNow {
		return claims, err
	}

	return claims, nil
}

func ParseJwt(cookie string, secretKey ...string) (*Claims, error) {
	secret := utils.GetString("SECRET_KEY")
	if len(secretKey) > 0 {
		secret = secretKey[0]
	}

	claims := Claims{}
	if token, err := jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}); err != nil || !token.Valid {
		return nil, err
	}

	return &claims, nil
}

func GetToken(c *fiber.Ctx) string {
	token := ""

	bearerToken := new(ReqHeaderAuth)
	if err := c.ReqHeaderParser(bearerToken); err == nil {
		token, _ = bearerToken.GetBearerToken()
	}

	tokenArr := []string{
		"token",
	}

	if len(token) == 0 {
		for _, v := range tokenArr {
			if c.Cookies(v) != "" {
				token = c.Cookies(v)
				break
			}
		}
	}

	return token
}

func (h *ReqHeaderAuth) GetBearerToken() (string, error) {
	if len(h.Authorization) == 0 {
		err := errors.New("authorization header not found")
		return "", err
	}

	authorization := strings.Split(h.Authorization, " ")
	if strings.ToLower(authorization[0]) != "bearer" {
		err := errors.New("not a bearer token")
		return "", err
	}

	return authorization[1], nil
}

func GenerateRefreshJwt(issuer string) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokens.SignedString([]byte(utils.GetString("REFRESH_KEY")))
}

func GenerateJwt(claim *Claims, sessionID uuid.UUID) (string, error) {
	claim.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return tokens.SignedString([]byte(utils.GetString("SECRET_KEY")))
}
