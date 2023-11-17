package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sash20m/go-api-template/internal/model"
	"github.com/sash20m/go-api-template/pkg/logger"
	"github.com/sirupsen/logrus"
)

// GetBooks godoc
//
//	@Summary		Get all books
//	@Tags			Books
//	@Produce		json
//	@Success		200	{object} []model.GetBookResponse
//
// @Router			/api/books [get]
func (h *Handlers) GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	books, err := h.Storage.GetBooks(ctx)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	bookResponse := []model.GetBookResponse{}
	for _, book := range books {
		bookResponse = append(bookResponse, model.GetBookResponse{
			ID:       book.ID,
			Title:    book.Title,
			Author:   book.Author,
			CoverURL: book.CoverURL,
			PostURL:  fmt.Sprint(book.PostURL),
		})

	}

	err = h.Sender.JSON(w, http.StatusOK, bookResponse)
	if err != nil {
		logger.OutputLog.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("Error when requesting /books")

		panic(err)
	}
}

// GetBook godoc
//
//	@Summary		Get a specific book
//	@Tags			Books
//	@Produce		json
//	@Param			id path int	true "Book ID"
//	@Success		200	{object} model.GetBookResponse
//
// @Router			/api/book/{id} [get]
func (h *Handlers) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	var book model.Book
	book, err = h.Storage.GetBook(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if (model.Book{}) == book {
		h.Sender.JSON(w, http.StatusBadRequest, "Book with id="+fmt.Sprint(book.ID)+" not found")
		if err != nil {
			panic(err)
		}
		return
	}

	bookResponse := model.GetBookResponse{
		ID:       book.ID,
		Title:    book.Title,
		Author:   book.Author,
		CoverURL: book.CoverURL,
		PostURL:  book.PostURL.String,
	}

	err = h.Sender.JSON(w, http.StatusOK, bookResponse)
	if err != nil {
		logger.OutputLog.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal(fmt.Sprint("Error when requesting /book/", book.ID))

		panic(err)
	}
}

// AddBook godoc
//
//		@Summary		Add a specific book
//		@Tags			Books
//		@Produce		json
//	 	@Accept			json
//		@Param			title body string true "Book title"
//		@Param			author body string true "Book author"
//		@Param			coverUrl body string true "Book coverUrl"
//		@Param			postUrl body string true "Book post url"
//		@Success		200	{object} model.IDResponse
//
// @Router			/api/book/add [post]
func (h *Handlers) AddBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var book model.AddBookRequest

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(book)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	id, err := h.Storage.AddBook(ctx, book)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := model.IDResponse{ID: id}
	err = h.Sender.JSON(w, http.StatusOK, response)
	if err != nil {
		panic(err)
	}
}

// UpdateBook godoc
//
//		@Summary		Update a specific book
//		@Tags			Books
//		@Produce		json
//	 	@Accept			json
//		@Param			title body string false "Book title"
//		@Param			author body string false "Book author"
//		@Param			coverUrl body string false "Book coverUrl"
//		@Param			postUrl body string false "Book post url"
//		@Success		200	{object} model.IDResponse
//
// @Router			/api/book/update [patch]
func (h *Handlers) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var book model.UpdateBookRequest

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(book)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	exists, err := h.Storage.VerifyBookExists(ctx, book.ID)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		h.Sender.JSON(w, http.StatusBadRequest, "Book with id="+fmt.Sprint(book.ID)+" not found")
		return
	}

	id, err := h.Storage.UpdateBook(ctx, book)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := model.IDResponse{ID: id}
	err = h.Sender.JSON(w, http.StatusOK, response)
	if err != nil {
		panic(err)
	}
}

// DeleteBook godoc
//
//	@Summary		Delete a specific book
//	@Tags			Books
//	@Produce		json
//	@Param			id path int	true "Book ID"
//	@Success		200	{object} model.GetBookResponse
//
// @Router			/api/book/delete/{id} [delete]
func (h *Handlers) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	exists, err := h.Storage.VerifyBookExists(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		err = h.Sender.JSON(w, http.StatusBadRequest, "Book with id="+fmt.Sprint(id)+" not found")
		if err != nil {
			panic(err)
		}
		return
	}

	err = h.Storage.DeleteBook(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Sender.JSON(w, http.StatusOK, map[string]bool{"success": true})
	if err != nil {
		panic(err)
	}
}

// cSpell:ignore godoc logrus
