package api

import (
	"errors"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/joe/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenClaims       = errors.New("can't extract claims from jwt token")
	ErrTokenInvalid      = errors.New("jwt token not valid")
	ErrTokenExpired      = errors.New("jwt token expired")
	ErrTokenIncomplete   = errors.New("jwt token is missing required fields")
	ErrTokenUnauthorized = errors.New("jwt token authorized field is false (?!)")
)

func validateToken(header string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenClaims
	} else if !token.Valid {
		return nil, ErrTokenInvalid
	}

	required := []string{
		"expires_at",
		"username",
		"authorized",
	}
	for _, req := range required {
		if _, found := claims[req]; !found {
			return nil, ErrTokenIncomplete
		}
	}

	log.Debug("%+v", claims)

	if expiresAt, err := time.Parse(time.RFC3339, claims["expires_at"].(string)); err != nil {
		return nil, ErrTokenExpired
	} else if expiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	} else if claims["authorized"].(bool) != true {
		return nil, ErrTokenUnauthorized
	}
	return claims, err
}

func getUser(r *http.Request) *models.User {
	client := clientIP(r)
	tokenHeader := reqToken(r)
	if tokenHeader == "" {
		log.Debug("unauthenticated request from %s", client)
		return nil
	}

	claims, err := validateToken(tokenHeader)
	if err != nil {
		log.Debug("token error for %s: %v", client, err)
		return nil
	}

	if u, found := models.Users.Load(claims["username"].(string)); found {
		return u.(*models.User)
	}

	return nil
}

// GET|POST /api/v1/auth
func (api *API) Authenticate(w http.ResponseWriter, r *http.Request) {
	client := clientIP(r)
	params := parseParameters(r)
	if username, found := params["user"]; !found {
		ERROR(w, http.StatusBadRequest, ErrEmpty)
	} else if password, found := params["pass"]; !found {
		ERROR(w, http.StatusBadRequest, ErrEmpty)
	} else if u, found := models.Users.Load(username); !found {
		log.Warning("%s: user %s not found", client, username)
		ERROR(w, http.StatusUnauthorized, ErrEmpty)
	} else if user := u.(*models.User); user.ValidPassword(password.(string)) == false {
		log.Warning("%s: invalid password", client)
		ERROR(w, http.StatusUnauthorized, ErrEmpty)
	} else {
		claims := jwt.MapClaims{}
		claims["authorized"] = true
		claims["username"] = username
		claims["expires_at"] = time.Now().Add(time.Hour * time.Duration(user.TokenTTL)).Format(time.RFC3339)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if token, err := token.SignedString([]byte(os.Getenv("API_SECRET"))); err != nil {
			log.Error("error creating token for %s: %v", client, err)
			ERROR(w, http.StatusInternalServerError, ErrEmpty)
		} else {
			log.Info("%s authenticated successfully", client)
			JSON(w, http.StatusOK, map[string]string{
				"token": token,
			})
		}
	}
}
