package interfaces

import (
	"github.com/google/uuid"
	"github.com/magzhan/library/internal/models"
)

type IBookRepository interface {
	GetBooks() ([]models.BookModel, error)
	CreateBook(book models.BookModel) (models.BookModel, error)
	UpdateBook(id uuid.UUID, book models.BookModel) (models.BookModel, error)
	DeleteBook(id uuid.UUID) error
}

type IAuthorRepository interface {
	GetAuthors() ([]models.AuthorModel, error)
	GetAuthorsBooks(id uuid.UUID) ([]models.BookModel, error)
	CreateAuthor(author models.AuthorModel) (models.AuthorModel, error)
	AddBooksToAuthor(authorID uuid.UUID, bookIDs []uuid.UUID) error
	UpdateAuthor(id uuid.UUID, author models.AuthorModel) (models.AuthorModel, error)
	DeleteAuthor(id uuid.UUID) error
}

type ICustomerRepository interface {
	GetCustomers() ([]models.CustomerModel, error)
	GetCustomersBooks(id uuid.UUID) ([]models.BookModel, error)
	CreateCustomer(customer models.CustomerModel) (models.CustomerModel, error)
	AddBooksToCustomer(authorID uuid.UUID, bookIDs []uuid.UUID) error
	UpdateCustomer(id uuid.UUID, customer models.CustomerModel) (models.CustomerModel, error)
	DeleteCustomer(id uuid.UUID) error
}
