package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
)

func write500(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "500 Server Error\n")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get Project ID
		projectId, err := metadata.ProjectID()
		if err != nil {
			write500(w, err)
			return
		}

		// Get ID Token
		audience := os.Getenv("ID_TOKEN_AUDIENCE")
		idToken, err := metadata.Get("instance/service-accounts/default/identity?audience=" + audience)
		if err != nil {
			write500(w, err)
			return
		}
		fmt.Printf("ID Token: %s\n", idToken)

		// Call backend service
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://inter-service-auth-backend-dot-%s.appspot.com", projectId), nil)
		req.Header.Add("Authorization", "Bearer "+idToken)
		if err != nil {
			write500(w, err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			write500(w, err)
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			write500(w, err)
			return
		}

		fmt.Fprintf(w, "Response from backend:\n  %s", string(b))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
