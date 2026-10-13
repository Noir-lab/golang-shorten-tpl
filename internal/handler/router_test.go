package handler

import (
	"net/http/httptest"
	"testing"

	"golang-shorten-tpl/internal/service"
)

const testBaseURL = "http://localhost: 8080"

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := service.NewMockService(ctrl)

	h := NewURLHandler(mockSvc, testBaseURL)

	srv := httptest.NewServer(NewRouter(h))
	defer srv.Close()
}
