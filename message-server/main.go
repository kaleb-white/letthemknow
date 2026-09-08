package main

import (
	"fmt"
	"net/http"

	"github.com/kaleb-white/letthemknow/handlers"
	"github.com/kaleb-white/letthemknow/schemas"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloWorld)
	fmt.Println("Starting server on 8080...")
	http.ListenAndServe(":8080", mux)
}
