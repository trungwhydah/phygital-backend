package author

import (
	"backend-service/internal/core_backend/entity"
)

// Author interface
type Author interface {
	CreateAuthor(author *entity.Author) (*entity.Author, error)
	GetListAuthor() (*[]entity.Author, error)
	GetAuthorByID(authorID *string) (*entity.Author, error)
}

// Repository interface
type Repository interface {
	Author
}

// UseCase interface
type UseCase interface {
	CreateAuthor(author *entity.Author) (*entity.Author, int, error)
	GetListAuthor() (*[]entity.Author, int, error)
	GetAuthorDetail(authorID *string) (*entity.Author, int, error)
}
