package main

import (
	_ "golang-shorten-tpl/api/docs"
	"golang-shorten-tpl/internal/handler"
	"golang-shorten-tpl/internal/repository"
	"golang-shorten-tpl/internal/service"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           URL Shortener API
// @version         1.0
// @description     Сервис сокращения URL
// @host            localhost:8080
// @BasePath        /
func main() {
	const (
		baseURL    = "http://localhost:8080"
		serverAddr = "localhost:8080"
	)
	//mux := http.NewServeMux()
	repo := repository.NewMemoryRepository()
	svc := service.NewURLService(repo)
	h := handler.NewURLHandler(svc, baseURL)
	r := handler.NewRouter(h)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	//mux.HandleFunc("POST /{$}", h.POSTHandler)
	//mux.HandleFunc("GET /{id}", h.GETHandler)
	//mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	//err := http.ListenAndServe(serverAddr, mux)
	//if err != nil {
	//	panic(err)
	//}

	err := http.ListenAndServe(serverAddr, r)
	if err != nil {
		return
	}
}
