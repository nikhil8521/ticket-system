package routes

import (
	"ticket-system/controllers"
	"ticket-system/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authCtrl *controllers.AuthController, ticketCtrl *controllers.TicketController) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
	}

	// Ticket routes (Protected)
	tickets := r.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		tickets.POST("", ticketCtrl.CreateTicket)
		tickets.GET("", ticketCtrl.GetTickets)
		tickets.GET("/:id", ticketCtrl.GetTicketByID)
		tickets.PATCH("/:id/status", ticketCtrl.UpdateTicketStatus)
	}

	return r
}
