package repositories

import (
	"database/sql"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/magzhan/library/internal/models"
)

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(DB *sql.DB) BookRepository {
	return BookRepository{
		DB: DB,
	}
}

func (repo BookRepository) GetBooks() ([]models.BookModel, error) {
	rows, err := repo.DB.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []models.BookModel
	for rows.Next() {
		var book models.BookModel
		err := rows.Scan(&book.BookID, &book.Title, &book.Genre, &book.ISBN)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (repo BookRepository) CreateBook(book models.BookModel) (models.BookModel, error) {
	book.BookID = uuid.New()
	_, err := repo.DB.Exec(`INSERT INTO books (BookID, Title, Genre, ISBN) VALUES ($1, $2, $3, $4)`, book.BookID, book.Title, book.Genre, book.ISBN)
	if err != nil {
		return models.BookModel{}, err
	}
	return book, nil
}

func (repo BookRepository) UpdateBook(id uuid.UUID, book models.BookModel) (models.BookModel, error) {
	oldBook := models.BookModel{}
	err := repo.DB.QueryRow(`SELECT * FROM books WHERE BookID = $1`, id).Scan(&oldBook.BookID, &oldBook.Title, &oldBook.Genre, &oldBook.ISBN)
	if err != nil {
		return models.BookModel{}, err
	}
	oldBook.Title = book.Title
	oldBook.Genre = book.Genre
	oldBook.ISBN = book.ISBN
	result, err := repo.DB.Exec("UPDATE books SET Title = $1, Genre = $2, ISBN = $3 WHERE BookID = $4", oldBook.Title, oldBook.Genre, oldBook.ISBN, id)
	if err != nil {
		return models.BookModel{}, err
	}
	rowsAffected, err := result.RowsAffected()
	log.Debug(rowsAffected)
	if err != nil {
		return models.BookModel{}, err
	}
	return oldBook, nil
}

func (repo BookRepository) DeleteBook(id uuid.UUID) error {
	_, err := repo.DB.Exec(`DELETE FROM books WHERE BookID = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
