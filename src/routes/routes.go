package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rawipassT/book-service/internal/http"
)

func SetupRoutes(bookHandler *http.BookHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/books", bookHandler.ListBooks)
	router.GET("/books/:id", bookHandler.FetchBookByID)
	router.DELETE("/books/:id", bookHandler.DeleteBook)
	router.POST("/books", bookHandler.CreateBook)
	router.PUT("/books/:id", bookHandler.UpdateBook)
	router.POST("/books/borrow", bookHandler.BorrowBook)
	router.POST("/books/return", bookHandler.ReturnBook)
	router.GET("/books/most_borrowed", bookHandler.GetMostBorrowedBooks)



	return router
}
