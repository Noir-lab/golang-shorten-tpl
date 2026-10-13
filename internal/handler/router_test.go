package handler

import (
	"golang-shorten-tpl/internal/service/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestRouter(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		body         string
		mockSetup    func(m *mocks.MockService)
		wantStatus   int
		wantLocation string
		wantBody     string
	}{
		{
			name:   "POST / router to POSTHandler",
			method: http.MethodPost,
			path:   "/",
			body:   "https://practicum.yandex.ru",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Shorten("https://practicum.yandex.ru").
					Return("abc12345", nil).
					Times(1)
			},
			wantStatus: http.StatusCreated,
			wantBody:   testBaseURL + "/abc12345",
		},
		{
			name:   "GET /{id} routes to GETHandler and passes id",
			method: http.MethodGet,
			path:   "/abc12345",
			mockSetup: func(m *mocks.MockService) {

				m.EXPECT().
					GetOriginal("abc12345").
					Return("https://practicum.yandex.ru/", true).
					Times(1)
			},
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "https://practicum.yandex.ru/",
		},
		{
			name:   "GET unknown id -> 400",
			method: http.MethodGet,
			path:   "/unknown",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					GetOriginal("unknown").
					Return("", false).
					Times(1)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "PUT / -> 405",
			method:     http.MethodPut,
			path:       "/",
			mockSetup:  func(m *mocks.MockService) {},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "GET / -> 405",
			method:     http.MethodGet,
			path:       "/",
			mockSetup:  func(m *mocks.MockService) {},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "POST /{id} -> 405",
			method:     http.MethodPost,
			path:       "/abc12345",
			body:       "https://practicum.yandex.ru/",
			mockSetup:  func(m *mocks.MockService) {},
			wantStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := mocks.NewMockService(ctrl)
			tt.mockSetup(mockSvc)

			srv := httptest.NewServer(NewRouter(NewURLHandler(mockSvc, testBaseURL)))
			defer srv.Close()

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			req, err := http.NewRequest(tt.method, srv.URL+tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			req.Header.Set("Content-Type", "text/plain")

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if got := res.Header.Get("Location"); got != tt.wantLocation {
				t.Errorf("Location: got %q, want %q", got, tt.wantLocation)
			}
			if tt.wantBody != "" {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}
				if string(body) != tt.wantBody {
					t.Errorf("body: got %q, want %q", body, tt.wantBody)
				}
			}
		})
	}
}
