package main

import (
	"github.com/magzhan/library/internal/data/repositories"
	"github.com/magzhan/library/internal/handlers"
	"github.com/magzhan/library/internal/initializers"
	"github.com/magzhan/library/internal/migration"
	"os"
)

func init() {
	//initializers.LoadEnvVar()
	initializers.LoadDatabase()
	migration.Migration()

}

func main() {
	bookRepo := repositories.NewBookRepository(initializers.DB)
	authorRepo := repositories.NewAuthorRepository(initializers.DB)
	customerRepo := repositories.NewCustomerRepository(initializers.DB)
	port := os.Getenv("PORT")
	server := handlers.NewHttpHandler(port, bookRepo, authorRepo, customerRepo)
	server.RepoBook = bookRepo
	server.RepoAuthor = authorRepo
	server.RepoCustomer = customerRepo
	server.Run()
}
