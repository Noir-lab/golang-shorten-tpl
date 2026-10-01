package handler

import (
	"errors"
	"golang-shorten-tpl/internal/service"
	"io"
	"log"
	"net/http"
)

type URLHandler struct {
	service service.Service
	baseURL string
}

func NewURLHandler(svc service.Service, baseURL string) *URLHandler {
	return &URLHandler{service: svc, baseURL: baseURL}
}

// POSTHandler
// @Summary      Сократить URL
// @Description  Принимает оригинальный URL в теле запроса, возвращает короткий URL
// @Tags         shorten
// @Accept       text/plain
// @Produce      text/plain
// @Param        url body string true "Оригинальный URL"
// @Success      201 {string} string "Сокращённый URL"
// @Failure      400 {string} string "Некорректный запрос"
// @Router       / [post]
func (h *URLHandler) POSTHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}
	originalURL := string(body)
	log.Printf("получено тело запроса: %q", originalURL)

	id, err := h.service.Shorten(originalURL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		http.Error(w, "invalid error", http.StatusBadRequest)
		return
	}

	shortURL := h.baseURL + "/" + id
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// GETHandler
// @Summary      Перейти по короткому URL
// @Description  Редиректит на оригинальный URL по идентификатору
// @Tags         shorten
// @Param        id path string true "Идентификатор короткого URL"
// @Success      307 "Temporary Redirect"
// @Failure      400 {string} string "URL не найден"
// @Router       /{id} [get]
func (h *URLHandler) GETHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	originalURL, found := h.service.GetOriginal(id)

	if !found {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
