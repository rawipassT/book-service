package usecase

import (
	"github.com/rawipassT/book-service/internal/repository"
	"github.com/rawipassT/book-service/models"
)

type BookUsecase struct {
	repo repository.BookRepository
}

func NewBookUseCase() *BookUsecase {
	usecase := BookUsecase{}
	return &usecase
}

func (u *BookUsecase) FetchBookByID(id string) (*models.Book, error) {
	return u.repo.FetchBookByID(id)
}

func (u *BookUsecase) ListBooks(title, author, category string) ([]models.Book, error) {
	return u.repo.ListBooks(title, author, category)
}

func (u *BookUsecase) DeleteBook(id string) error {
	return u.repo.DeleteBook(id)
}

func (u *BookUsecase) CreateBook(book *models.Book) error {
	return u.repo.CreateBook(book)
}

func (u *BookUsecase) UpdateBook(id string, book *models.Book) error {
	return u.repo.UpdateBook(id, book)
}

func (u *BookUsecase) BorrowBook(userID, bookID string) error {
	return u.repo.BorrowBook(userID, bookID)
}

func (u *BookUsecase) ReturnBook(userID, bookID string) error {
	return u.repo.ReturnBook(userID, bookID)
}

func (u *BookUsecase) ListMostBorrowedBooks(limit int) ([]models.Book, error) {
	return u.repo.ListMostBorrowedBooks(limit)
}

