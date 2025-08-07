package repositories

import (
	"database/sql"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/magzhan/library/internal/models"
)

type AuthorRepository struct {
	DB *sql.DB
}

func NewAuthorRepository(DB *sql.DB) AuthorRepository {
	return AuthorRepository{
		DB: DB,
	}
}

func (repo AuthorRepository) GetAuthors() ([]models.AuthorModel, error) {
	rows, err := repo.DB.Query("SELECT * FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var authors []models.AuthorModel
	for rows.Next() {
		var author models.AuthorModel
		err := rows.Scan(&author.AuthorID, &author.FirstName, &author.LastName, &author.ThirdName, &author.Nickname, &author.Specialization)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return authors, nil
}

func (repo AuthorRepository) GetAuthorsBooks(id uuid.UUID) ([]models.BookModel, error) {
	query := `SELECT b.bookID, b.title, b.genre, b.ISBN
				FROM books b
				INNER JOIN authors_books ab ON b.bookID = ab.bookID
				WHERE ab.AuthorID = $1`
	rows, err := repo.DB.Query(query, id)
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
func (repo AuthorRepository) LinkAuthorToBook(authorID, bookID uuid.UUID) error {
	_, err := repo.DB.Exec(`
		INSERT INTO authors_books (AuthorID, BookID)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, authorID, bookID)
	return err
}
func (repo AuthorRepository) UnlinkAuthorFromBook(authorID, bookID uuid.UUID) error {
	_, err := repo.DB.Exec(`
		DELETE FROM authors_books
		WHERE AuthorID = $1 AND BookID = $2
	`, authorID, bookID)
	return err
}

func (repo AuthorRepository) AddBooksToAuthor(authorID uuid.UUID, bookIDs []uuid.UUID) error {
	for _, bookID := range bookIDs {
		err := repo.LinkAuthorToBook(authorID, bookID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo AuthorRepository) CreateAuthor(author models.AuthorModel) (models.AuthorModel, error) {
	author.AuthorID = uuid.New()
	_, err := repo.DB.Exec(
		`INSERT INTO authors (AuthorID, FirstName, LastName, ThirdName, Nickname, Specialization)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		author.AuthorID,
		author.FirstName,
		author.LastName,
		author.ThirdName,
		author.Nickname,
		author.Specialization,
	)
	if err != nil {
		return models.AuthorModel{}, err
	}
	return author, nil
}

func (repo AuthorRepository) UpdateAuthor(id uuid.UUID, author models.AuthorModel) (models.AuthorModel, error) {
	oldAuthor := models.AuthorModel{}
	err := repo.DB.QueryRow(`SELECT * FROM authors WHERE AuthorID = $1`, id).Scan(
		&oldAuthor.AuthorID,
		&oldAuthor.FirstName,
		&oldAuthor.LastName,
		&oldAuthor.ThirdName,
		&oldAuthor.Nickname,
		&oldAuthor.Specialization,
	)
	if err != nil {
		return models.AuthorModel{}, err
	}
	oldAuthor.FirstName = author.FirstName
	oldAuthor.LastName = author.LastName
	oldAuthor.ThirdName = author.ThirdName
	oldAuthor.Nickname = author.Nickname
	oldAuthor.Specialization = author.Specialization
	result, err := repo.DB.Exec(`UPDATE authors SET FirstName = $1, LastName = $2, ThirdName = $3,Nickname=$4, Specialization=$5 WHERE AuthorID = $6`,
		oldAuthor.FirstName,
		oldAuthor.LastName,
		oldAuthor.ThirdName,
		oldAuthor.Nickname,
		oldAuthor.Specialization,
		id)
	if err != nil {
		return models.AuthorModel{}, err
	}
	rowsAffected, err := result.RowsAffected()
	log.Debug(rowsAffected)
	if err != nil {
		return models.AuthorModel{}, err
	}
	if rowsAffected == 0 {
		log.Warn("No rows updated for author")
	}
	return oldAuthor, nil
}

func (repo AuthorRepository) DeleteAuthor(id uuid.UUID) error {
	_, err := repo.DB.Exec(`DELETE FROM authors WHERE AuthorID = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
