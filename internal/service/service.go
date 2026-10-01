package service

type Repository interface {
	Save(id string, originalURL string) error
	Get(id string) (originalURL string, err bool)
}

type Service interface {
	Shorten(originalURL string) (shortID string, err error)
	GetOriginal(id string) (originalURL string, err bool)
}
