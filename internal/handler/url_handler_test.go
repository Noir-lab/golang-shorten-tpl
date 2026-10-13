package handler

import (
	"errors"
	"golang-shorten-tpl/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang-shorten-tpl/internal/service/mocks"

	"go.uber.org/mock/gomock"
)

const testBaseURL = "http://localhost: 8080"

func TestURLHandler_POSTHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(m *mocks.MockService)
		wantStatusCode int
		wantBody       string
		wantCT         string
	}{
		{
			name:        "valid url -> 201 and short url",
			requestBody: "https://practicum.yandex.ru/",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Shorten("https://practicum.yandex.ru/").
					Return("abc12345", nil).
					Times(1)
			},
			wantStatusCode: http.StatusCreated,
			wantBody:       testBaseURL + "/abc12345",
			wantCT:         "text/plain",
		},
		{
			name:        "invalid url -> 400",
			requestBody: "not-a-url",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Shorten("not-a-url").
					Return("", service.ErrInvalidURL).
					Times(1)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:        "empty body -> 400",
			requestBody: "",
			mockSetup: func(m *mocks.MockService) {
				// Валидацию делает service, handler лишь передаёт тело как есть.
				m.EXPECT().
					Shorten("").
					Return("", service.ErrInvalidURL).
					Times(1)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:        "unexpected service error -> 500",
			requestBody: "https://practicum.yandex.ru/",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					Shorten("https://practicum.yandex.ru/").
					Return("", errors.New("storage is down")).
					Times(1)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := mocks.NewMockService(ctrl)
			tt.mockSetup(mockSvc)
			h := NewURLHandler(mockSvc, testBaseURL)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()

			h.POSTHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatusCode {
				t.Fatalf("status: got %d, want %d", res.StatusCode, tt.wantStatusCode)
			}

			if tt.wantCT != "" {
				if got := res.Header.Get("Content-Type"); !strings.HasPrefix(got, tt.wantCT) {
					t.Errorf("Content-Type: got %q, want prefix %q", got, tt.wantCT)
				}
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

func TestURLHandler_GETHandler(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockSetup    func(m *mocks.MockService)
		wantStatus   int
		wantLocation string
	}{
		{
			name: "existing id -> 307 and Location",
			id:   "abc12345",
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
			name: "unknown id -> 400 and no Location",
			id:   "missing",
			mockSetup: func(m *mocks.MockService) {
				m.EXPECT().
					GetOriginal("missing").
					Return("", false).
					Times(1)
			},
			wantStatus:   http.StatusBadRequest,
			wantLocation: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := mocks.NewMockService(ctrl)
			tt.mockSetup(mockSvc)

			h := NewURLHandler(mockSvc, testBaseURL)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			// Хендлер вызывается напрямую, мимо ServeMux, поэтому {id}
			// никто не разберёт: подставляем path value вручную (Go 1.22+).
			req.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			h.GETHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if got := res.Header.Get("Location"); got != tt.wantLocation {
				t.Errorf("Location: got %q, want %q", got, tt.wantLocation)
			}
		})
	}
}
