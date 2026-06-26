package controllers

import (
	"net/http"
	"strconv"
	"ticket-system/models"
	"ticket-system/services"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService services.TicketService
}

func NewTicketController(ticketService services.TicketService) *TicketController {
	return &TicketController{ticketService: ticketService}
}

func (ctrl *TicketController) CreateTicket(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)
	ticket, err := ctrl.ticketService.CreateTicket(userID, input.Title, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (ctrl *TicketController) GetTickets(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	tickets, err := ctrl.ticketService.GetTicketsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (ctrl *TicketController) GetTicketByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	userID := c.MustGet("userID").(uint)
	ticket, err := ctrl.ticketService.GetTicketByID(uint(id), userID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func (ctrl *TicketController) UpdateTicketStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var input struct {
		Status models.TicketStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status enum
	if input.Status != models.StatusOpen && input.Status != models.StatusInProgress && input.Status != models.StatusClosed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	userID := c.MustGet("userID").(uint)
	err = ctrl.ticketService.UpdateTicketStatus(uint(id), userID, input.Status)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
