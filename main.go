package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/auth0-community/auth0"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	jose "gopkg.in/square/go-jose.v2"
)

type Product struct {
	Id          int
	Name        string
	Slug        string
	Description string
}

var products = []Product{
	Product{Id: 1, Name: "Hover Shooters", Slug: "hover-shooters", Description: "Shoot your way to the top on 14 different hoverboards"},
	Product{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description: "Explore the depths of the sea in this one of a kind underwater experience"},
	Product{Id: 3, Name: "Dinosaur Park", Slug: "dinosaur-park", Description: "Go back 65 million years in the past and ride a T-Rex"},
	Product{Id: 4, Name: "Cars VR", Slug: "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
	Product{Id: 5, Name: "Robin Hood", Slug: "robin-hood", Description: "Pick up the bow and arrow and master the art of archery"},
	Product{Id: 6, Name: "Real World VR", Slug: "real-world-vr", Description: "Explore the seven wonders of the world in VR"},
}

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/products", authMiddleware(ProductsHandler)).Methods("GET")
	r.Handle("/products/{slug}/feedback", authMiddleware(AddFeedbackHandler)).Methods("POST")

	fmt.Println("Ready to listen for request on port 3000")
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// secret, _ := base64.URLEncoding.DecodeString("l9rqay1dsOYc4D7SQqHKDTA1rY0FuvfO")
		//
		// secretProvider := auth0.NewKeyProvider(secret)
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: "https://avantidev.auth0.com/.well-known/jwks.json"})
		audience := []string{"http://localhost:8080", "https://avantidev.auth0.com/userinfo"}
		//iss, _ := base64.URLEncoding.DecodeString("https://avantidev.auth0.com")

		configuration := auth0.NewConfiguration(client, audience, "https://avantidev.auth0.com/", jose.RS256)

		validator := auth0.NewValidator(configuration)

		spew.Dump(validator)

		token, err := validator.ValidateRequest(r)
		fmt.Println("##########Token###############")
		spew.Dump(token)
		fmt.Println(token)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
		} else {
			fmt.Println("Token validated successfully")
			fmt.Println(token)
			next.ServeHTTP(w, r)
		}
	})
}

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	slug := vars["slug"]

	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if product.Slug != "" {
		payload, _ := json.Marshal(product)
		w.Write([]byte(payload))
	} else {
		w.Write([]byte("Product Not Found"))
	}
})
