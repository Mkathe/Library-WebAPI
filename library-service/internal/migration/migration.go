package migration

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/magzhan/library/internal/initializers"
)

func Migration() {
	db := initializers.DB

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS books (
		BookID UUID PRIMARY KEY,
		Title VARCHAR(255) NOT NULL,
		Genre VARCHAR(255) NOT NULL,
		ISBN VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		log.Fatal("Error creating books table:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS authors (
		AuthorID UUID PRIMARY KEY,
		FirstName VARCHAR(255) NOT NULL,
		LastName VARCHAR(255) NOT NULL,
		ThirdName VARCHAR(255) NOT NULL,
    	Nickname VARCHAR(255),
    	Specialization VARCHAR(255)
	)`)
	if err != nil {
		log.Fatal("Error creating authors table:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS customers (
		CustomerID UUID PRIMARY KEY,
		FirstName VARCHAR(255) NOT NULL,
		LastName VARCHAR(255) NOT NULL,
		ThirdName VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		log.Fatal("Error creating customers table:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS customer_books (
		CustomerID UUID,
		BookID UUID,
		PRIMARY KEY (CustomerID, BookID),
		FOREIGN KEY (CustomerID) REFERENCES customers(CustomerID) ON DELETE CASCADE,
		FOREIGN KEY (BookID) REFERENCES books(BookID) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal("Error creating customer_books table:", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS authors_books (
		AuthorID UUID,
		BookID UUID,
		PRIMARY KEY (AuthorID, BookID),
		FOREIGN KEY (AuthorID) REFERENCES authors(AuthorID) ON DELETE CASCADE,
		FOREIGN KEY (BookID) REFERENCES books(BookID) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal("Error creating authors_books table:", err)
	}
	log.Info("✅ Migration completed successfully")
}
