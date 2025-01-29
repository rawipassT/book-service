package repository

import (
	"context"
	"errors"

	"github.com/rawipassT/book-service/config"
	"github.com/rawipassT/book-service/models"
	"github.com/jackc/pgx/v4"
	
)

type BookRepository struct {
}

func NewBookRepository() *BookRepository {
	repo := BookRepository{}
	return &repo
}

func (r *BookRepository) ListBooks(title, author, category string) ([]models.Book, error) {
	query := `
		SELECT id, title, author, category, status, borrow_count, created_at
		FROM books
		WHERE ($1 = '' OR LOWER(title) LIKE LOWER('%' || $1 || '%'))
		AND ($2 = '' OR LOWER(author) LIKE LOWER('%' || $2 || '%'))
		AND ($3 = '' OR LOWER(category) LIKE LOWER('%' || $3 || '%'))
		ORDER BY created_at DESC
	`

	rows, err := config.DB.Query(context.Background(), query, title, author, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Category,
			&book.Status,
			&book.BorrowCount,
			&book.CreatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (r *BookRepository) FetchBookByID(id string) (*models.Book, error) {
	query := `
		SELECT id, title, author, category, status, borrow_count, created_at 
		FROM books 
		WHERE id = $1
	`
	var book models.Book
	err := config.DB.QueryRow(context.Background(), query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Category,
		&book.Status,
		&book.BorrowCount,
		&book.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) DeleteBook(id string) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := config.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}

func (r *BookRepository) CreateBook(book *models.Book) error {
	query := `
		INSERT INTO books (title, author, category, created_at)
		VALUES ($1, $2, $3, NOW())
	`

	_, err := config.DB.Exec(
		context.Background(),
		query,
		book.Title, book.Author, book.Category,
	)

	return err
}

func (r *BookRepository) UpdateBook(id string, book *models.Book) error {
	query := `
		UPDATE books 
		SET title = $1, author = $2, category = $3
		WHERE id = $4
	`

	result, err := config.DB.Exec(
		context.Background(),
		query,
		book.Title, book.Author, book.Category, id,
	)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}

func (r *BookRepository) BorrowBook(userID, bookID string) error {
	tx, err := config.DB.Begin(context.Background()) 
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) 

	var status string
	err = tx.QueryRow(context.Background(), "SELECT status FROM books WHERE id = $1", bookID).Scan(&status)
	if err != nil {
		return err
	}
	if status != "available" {
		return errors.New("book is not available for borrowing")
	}

	_, err = tx.Exec(
		context.Background(),
		`INSERT INTO borrow_records (id, user_id, book_id, borrowed_at) VALUES (uuid_generate_v4(), $1, $2, NOW())`,
		userID, bookID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE books SET status = 'borrowed', borrow_count = borrow_count + 1 WHERE id = $1`,
		bookID,
	)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(context.Background())
}

func (r *BookRepository) ReturnBook(userID, bookID string) error {
	tx, err := config.DB.Begin(context.Background()) 
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) 

	var recordID string
		err = tx.QueryRow(
		context.Background(),
		`SELECT id FROM borrow_records WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL`,
		userID, bookID,
	).Scan(&recordID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("no active borrow record found")
		}
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE borrow_records SET returned_at = NOW() WHERE id = $1`,
		recordID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE books SET status = 'available' WHERE id = $1`,
		bookID,
	)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(context.Background())
}

func (r *BookRepository) ListMostBorrowedBooks(limit int) ([]models.Book, error) {
	query := `
		SELECT id, title, author, category, status, borrow_count, created_at
		FROM books
		ORDER BY borrow_count DESC
		LIMIT $1
	`
	rows, err := config.DB.Query(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Category,
			&book.Status,
			&book.BorrowCount,
			&book.CreatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}


