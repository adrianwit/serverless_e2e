package main

import (
	"fmt"
	"github.com/adrianwit/serverless_e2e/cloud_function/gcetoggle"
	"log"
	"net/http"
	"os"
)

func main()  {
	http.HandleFunc("/",gcetoggle.HTTPGCEToggleFn)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
