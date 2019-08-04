package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/compute/metadata"
	jwt "github.com/dgrijalva/jwt-go"
)

type IdTokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func write401(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "401 Unauthorized\n")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Fetch Authorization Header
		bearerHeader := r.Header.Get("Authorization")
		if bearerHeader == "" {
			write401(w, fmt.Errorf("No Authorization header found"))
			return
		}

		re := regexp.MustCompile(`^\s*Bearer\s+(.+)$`)
		matched := re.FindStringSubmatch(bearerHeader)
		if len(matched) != 2 {
			write401(w, fmt.Errorf("Authorization header is invalid format"))
			return
		}
		bearerToken := matched[1]

		// Verify ID Token
		token, err := jwt.ParseWithClaims(bearerToken, &IdTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			kid := token.Header["kid"].(string)

			// Get certificate
			resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			decoder := json.NewDecoder(resp.Body)
			var jsonBody interface{}
			if err := decoder.Decode(&jsonBody); err != nil {
				return nil, err
			}
			cert := jsonBody.(map[string]interface{})[kid].(string)

			return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		})
		if err != nil {
			write401(w, fmt.Errorf("Invalid token: %s", err))
			return
		}

		claims, ok := token.Claims.(*IdTokenClaims)
		if !(ok && token.Valid) {
			write401(w, fmt.Errorf("Invalid token"))
			return
		}

		// Get Project ID
		projectId, err := metadata.ProjectID()
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Error\n")
			return
		}

		// Check if the request came from same application
		if claims.Email != fmt.Sprintf("%s@appspot.gserviceaccount.com", projectId) {
			write401(w, fmt.Errorf("Invalid token: email is invalid, %s\n", claims.Email))
			return
		}

		fmt.Printf("%#v\n", token)
		fmt.Fprintf(w, "Request by: %s\n", claims.Email)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
