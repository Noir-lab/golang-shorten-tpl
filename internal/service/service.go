package service

//go:generate go tool mockgen -source=service.go -destination=mocks/mock_service.go -package=mocks

type Repository interface {
	Save(id string, originalURL string) error
	Get(id string) (originalURL string, err bool)
}

type Service interface {
	Shorten(originalURL string) (shortID string, err error)
	GetOriginal(id string) (originalURL string, err bool)
}
