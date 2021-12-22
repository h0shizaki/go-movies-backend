package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			//Anonymous user
		}

		headerPart := strings.Split(authHeader, " ")

		if len(headerPart) != 2 {
			app.logger.Println("error :" + authHeader)
			app.errorJSON(w, errors.New("Invalid auth header"))
			return
		}

		if headerPart[0] != "Bearer" {
			app.errorJSON(w, errors.New("Unauthorized - no bearer"))
			return
		}

		token := headerPart[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))

		if err != nil {
			app.errorJSON(w, errors.New("Unauthorized - failed HMACCheck"))
			return
		}

		if !claims.Valid(time.Now()) {
			app.errorJSON(w, errors.New("Unauthorized - token expired"))
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			app.errorJSON(w, errors.New("Unauthorized - invalid audience"))
			return
		}

		if claims.Issuer != "mydomain.com" {
			app.errorJSON(w, errors.New("Unauthorized - invalid issuer"))
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)

		if err != nil {
			app.errorJSON(w, errors.New("Unauthorized - Unauthorized"))
			return
		}

		log.Println(userID)

		next.ServeHTTP(w, r)
	})
}
