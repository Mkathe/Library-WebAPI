package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/magzhan/library/internal/dto"
	"github.com/magzhan/library/internal/interfaces"
	"github.com/magzhan/library/internal/models"
)

type HttpHandler struct {
	addr         string
	RepoBook     interfaces.IBookRepository
	RepoAuthor   interfaces.IAuthorRepository
	RepoCustomer interfaces.ICustomerRepository
}

func NewHttpHandler(addr string, repoBook interfaces.IBookRepository, repoAuthor interfaces.IAuthorRepository, repoCustomer interfaces.ICustomerRepository) *HttpHandler {
	return &HttpHandler{
		addr:         addr,
		RepoBook:     repoBook,
		RepoAuthor:   repoAuthor,
		RepoCustomer: repoCustomer,
	}
}

func (h *HttpHandler) Run() {
	app := fiber.New()
	app.Get("/books", h.GetBooks)
	app.Post("/books", h.CreateBook)
	app.Patch("/books/:id", h.UpdateBook)
	app.Delete("/books/:id", h.DeleteBook)
	app.Get("/authors", h.GetAuthors)
	app.Get("/authors/:id/books", h.GetAuthorsBooks)
	app.Post("/authors", h.CreateAuthor)
	app.Post("/authors/:authorid/books", h.AddBooksToAuthor)
	app.Patch("/authors/:id", h.UpdateAuthor)
	app.Delete("/authors/:id", h.DeleteAuthor)
	app.Get("/members", h.GetCustomers)
	app.Post("/members", h.CreateCustomer)
	app.Get("/members/:id/books", h.GetCustomersBooks)
	app.Post("/members/:customerid/books", h.AddBooksToCustomer)
	app.Patch("/members/:id", h.UpdateCustomer)
	app.Delete("/members/:id", h.DeleteCustomer)
	log.Fatal(app.Listen(h.addr))
}

func (h *HttpHandler) GetBooks(ctx fiber.Ctx) error {
	books, err := h.RepoBook.GetBooks()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ctx.JSON(fiber.Map{
		"books": books,
	})
	return nil
}

func (h *HttpHandler) DeleteBook(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	if err := h.RepoBook.DeleteBook(parsedID); err != nil {
		log.Error(err)
		return err
	}
	ctx.SendStatus(fiber.StatusNoContent)
	return nil
}

func (h *HttpHandler) CreateBook(ctx fiber.Ctx) error {
	var book models.BookModel
	err := ctx.Bind().Body(&book)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	_, err = h.RepoBook.CreateBook(book)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ctx.SendStatus(fiber.StatusCreated)
	return nil
}

func (h *HttpHandler) UpdateBook(ctx fiber.Ctx) error {
	var book models.BookModel
	err := ctx.Bind().Body(&book)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	id := ctx.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	_, err = h.RepoBook.UpdateBook(parsedID, book)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ctx.SendStatus(fiber.StatusNoContent)
	return nil
}

func (h *HttpHandler) GetAuthors(ctx fiber.Ctx) error {
	authors, err := h.RepoAuthor.GetAuthors()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ctx.JSON(fiber.Map{
		"authors": authors,
	})
	return nil
}

func (h *HttpHandler) GetAuthorsBooks(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
	}
	books, err := h.RepoAuthor.GetAuthorsBooks(parsedID)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ctx.JSON(fiber.Map{
		"Author's books": books,
	})
	return nil
}

func (h *HttpHandler) CreateAuthor(ctx fiber.Ctx) error {
	var author models.AuthorModel
	err := ctx.Bind().Body(&author)
	if err != nil {
		log.Error(err)
	}
	_, err = h.RepoAuthor.CreateAuthor(author)
	if err != nil {
		log.Error(err)
	}
	ctx.SendStatus(fiber.StatusCreated)
	return nil
}
func (h *HttpHandler) UpdateAuthor(ctx fiber.Ctx) error {
	var author models.AuthorModel
	err := ctx.Bind().Body(&author)
	if err != nil {
		log.Error(err)
	}
	id := ctx.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
	}
	_, err = h.RepoAuthor.UpdateAuthor(parsedID, author)
	if err != nil {
		log.Error(err)
	}
	ctx.SendStatus(fiber.StatusNoContent)
	return nil
}
func (h *HttpHandler) DeleteAuthor(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
	}
	if err := h.RepoAuthor.DeleteAuthor(parsedID); err != nil {
		log.Error(err)
		return err
	}
	ctx.SendStatus(fiber.StatusNoContent)
	return nil
}

func (h *HttpHandler) AddBooksToAuthor(ctx fiber.Ctx) error {
	authorIDStr := ctx.Params("authorid")
	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		log.Error(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid author id",
		})
	}
	var req dto.AddBooksToAuthorRequest
	if err := ctx.Bind().Body(&req); err != nil {
		log.Error(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if len(req.BookIDs) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bookIDs list cannot be empty",
		})
	}

	err = h.RepoAuthor.AddBooksToAuthor(authorID, req.BookIDs)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to add books to author",
		})
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *HttpHandler) GetCustomers(ctx fiber.Ctx) error {
	customers, err := h.RepoCustomer.GetCustomers()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ctx.JSON(fiber.Map{
		"customers": customers,
	})
	return nil
}

func (h *HttpHandler) CreateCustomer(ctx fiber.Ctx) error {
	var customer models.CustomerModel
	err := ctx.Bind().Body(&customer)
	if err != nil {
		log.Error(err)
	}
	_, err = h.RepoCustomer.CreateCustomer(customer)
	if err != nil {
		log.Error(err)
	}
	ctx.SendStatus(fiber.StatusCreated)
	return nil
}
func (h *HttpHandler) GetCustomersBooks(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
	}
	books, err := h.RepoCustomer.GetCustomersBooks(parsedID)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ctx.JSON(fiber.Map{
		"Author's books": books,
	})
	return nil
}

func (h *HttpHandler) AddBooksToCustomer(ctx fiber.Ctx) error {
	ctx.Request().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	customerIDStr := ctx.Params("customerid")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		log.Error(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid author id",
		})
	}
	var req dto.AddBooksToCustomerRequest
	if err := ctx.Bind().Body(&req); err != nil {
		log.Error(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if len(req.BookIDs) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bookIDs list cannot be empty",
		})
	}

	err = h.RepoCustomer.AddBooksToCustomer(customerID, req.BookIDs)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to add books to customer",
		})
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *HttpHandler) UpdateCustomer(ctx fiber.Ctx) error {
	var customer models.CustomerModel
	err := ctx.Bind().Body(&customer)
	if err != nil {
		log.Error(err)
	}
	id := ctx.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
	}
	_, err = h.RepoCustomer.UpdateCustomer(parsedID, customer)
	if err != nil {
		log.Error(err)
	}
	ctx.SendStatus(fiber.StatusNoContent)
	return nil
}

func (h *HttpHandler) DeleteCustomer(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Error(err)
	}
	if err := h.RepoCustomer.DeleteCustomer(parsedID); err != nil {
		log.Error(err)
		return err
	}
	ctx.SendStatus(fiber.StatusNoContent)
	return nil
}
