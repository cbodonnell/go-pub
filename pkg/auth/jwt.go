package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims holds the JWT payload for an authenticated user.
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	jwt.RegisteredClaims
}

// RefreshClaims holds the JWT payload for a refresh token.
type RefreshClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTHelper creates and validates JWTs and refresh tokens.
type JWTHelper struct {
	secret        []byte
	jwtMaxAge     time.Duration
	refreshMaxAge time.Duration
	secure        bool // set Secure flag on cookies
}

func NewJWTHelper(secret string, jwtMaxAgeSecs, refreshMaxAgeSecs int, secure bool) *JWTHelper {
	return &JWTHelper{
		secret:        []byte(secret),
		jwtMaxAge:     time.Duration(jwtMaxAgeSecs) * time.Second,
		refreshMaxAge: time.Duration(refreshMaxAgeSecs) * time.Second,
		secure:        secure,
	}
}

// CreateJWT mints a signed JWT for the given user.
func (j *JWTHelper) CreateJWT(userID int, username, userUUID string) (string, *Claims, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		UUID:     userUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.jwtMaxAge)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	return signed, claims, err
}

// CreateRefresh mints a signed refresh token for the given user.
func (j *JWTHelper) CreateRefresh(userID int) (string, string, error) {
	jti := uuid.New().String()
	claims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshMaxAge)),
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	return signed, jti, err
}

// ParseJWT reads and validates the jwt cookie from the request.
func (j *JWTHelper) ParseJWT(r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, err
	}
	claims := &Claims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ParseRefresh reads and validates the refresh cookie from the request.
func (j *JWTHelper) ParseRefresh(r *http.Request) (*RefreshClaims, error) {
	cookie, err := r.Cookie("refresh")
	if err != nil {
		return nil, err
	}
	claims := &RefreshClaims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// SetJWTCookie writes the jwt cookie to the response.
func (j *JWTHelper) SetJWTCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    value,
		Path:     "/",
		MaxAge:   int(j.jwtMaxAge.Seconds()),
		HttpOnly: true,
		Secure:   j.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// SetRefreshCookie writes the refresh cookie to the response.
func (j *JWTHelper) SetRefreshCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    value,
		Path:     "/",
		MaxAge:   int(j.refreshMaxAge.Seconds()),
		HttpOnly: true,
		Secure:   j.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearCookies removes both auth cookies from the response.
func (j *JWTHelper) ClearCookies(w http.ResponseWriter) {
	for _, name := range []string{"jwt", "refresh"} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   j.secure,
			SameSite: http.SameSiteLaxMode,
		})
	}
}
