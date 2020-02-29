package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

const Address = ":8080"

func main() {
	r := chi.NewRouter()
	r.Get("/register", handler.GetRegister)
	http.ListenAndServe(Address, r)
}
