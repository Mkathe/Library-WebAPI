package repositories

import (
	"database/sql"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/magzhan/library/internal/models"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(DB *sql.DB) CustomerRepository {
	return CustomerRepository{
		DB: DB,
	}
}

func (repo CustomerRepository) GetCustomers() ([]models.CustomerModel, error) {
	rows, err := repo.DB.Query("SELECT * FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var customers []models.CustomerModel
	for rows.Next() {
		var customer models.CustomerModel
		err := rows.Scan(&customer.CustomerID, &customer.FirstName, &customer.LastName, &customer.ThirdName)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}
func (repo CustomerRepository) GetCustomersBooks(id uuid.UUID) ([]models.BookModel, error) {
	query := `SELECT b.bookID, b.title, b.genre, b.ISBN
				FROM books b
				INNER JOIN customer_books cb ON b.bookID = cb.bookID
				WHERE cb.CustomerID = $1`
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
func (repo CustomerRepository) CreateCustomer(customer models.CustomerModel) (models.CustomerModel, error) {
	customer.CustomerID = uuid.New()
	_, err := repo.DB.Exec(
		`INSERT INTO customers (CustomerID, FirstName, LastName, ThirdName)
	VALUES ($1, $2, $3, $4)`,
		customer.CustomerID,
		customer.FirstName,
		customer.LastName,
		customer.ThirdName,
	)
	if err != nil {
		return models.CustomerModel{}, err
	}
	return customer, nil
}

func (repo CustomerRepository) UpdateCustomer(id uuid.UUID, customer models.CustomerModel) (models.CustomerModel, error) {
	oldCustomer := models.CustomerModel{}
	err := repo.DB.QueryRow(`SELECT * FROM customers WHERE CustomerID = $1`, id).Scan(
		&oldCustomer.CustomerID,
		&oldCustomer.FirstName,
		&oldCustomer.LastName,
		&oldCustomer.ThirdName,
	)
	if err != nil {
		return models.CustomerModel{}, err
	}
	oldCustomer.FirstName = customer.FirstName
	oldCustomer.LastName = customer.LastName
	oldCustomer.ThirdName = customer.ThirdName
	result, err := repo.DB.Exec(`UPDATE customers SET FirstName = $1, LastName = $2, ThirdName = $3 WHERE AuthorID = $4`,
		oldCustomer.FirstName,
		oldCustomer.LastName,
		oldCustomer.ThirdName,
		id)
	if err != nil {
		return models.CustomerModel{}, err
	}
	rowsAffected, err := result.RowsAffected()
	log.Debug(rowsAffected)
	if err != nil {
		return models.CustomerModel{}, err
	}
	if rowsAffected == 0 {
		log.Warn("No rows updated for author")
	}
	return oldCustomer, nil
}
func (repo CustomerRepository) LinkCustomerToBook(customerID, bookID uuid.UUID) error {
	_, err := repo.DB.Exec(`
		INSERT INTO customer_books (CustomerID, BookID)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, customerID, bookID)
	return err
}
func (repo CustomerRepository) UnlinkCustomerFromBook(customerID, bookID uuid.UUID) error {
	_, err := repo.DB.Exec(`
		DELETE FROM customer_books
		WHERE CustomerID = $1 AND BookID = $2
	`, customerID, bookID)
	return err
}

func (repo CustomerRepository) AddBooksToCustomer(customerID uuid.UUID, bookIDs []uuid.UUID) error {
	for _, bookID := range bookIDs {
		err := repo.LinkCustomerToBook(customerID, bookID)
		if err != nil {
			return err
		}
	}
	return nil
}
func (repo CustomerRepository) DeleteCustomer(id uuid.UUID) error {
	_, err := repo.DB.Exec(`DELETE FROM customers WHERE CustomerID = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
