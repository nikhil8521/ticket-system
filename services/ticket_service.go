package services

import (
	"errors"
	"ticket-system/models"
	"ticket-system/repository"
)

type TicketService interface {
	CreateTicket(userID uint, title, description string) (*models.Ticket, error)
	GetTicketsByUserID(userID uint) ([]models.Ticket, error)
	GetTicketByID(id, userID uint) (*models.Ticket, error)
	UpdateTicketStatus(id, userID uint, status models.TicketStatus) error
}

type ticketService struct {
	ticketRepo repository.TicketRepository
}

func NewTicketService(ticketRepo repository.TicketRepository) TicketService {
	return &ticketService{ticketRepo: ticketRepo}
}

func (s *ticketService) CreateTicket(userID uint, title, description string) (*models.Ticket, error) {
	ticket := &models.Ticket{
		Title:       title,
		Description: description,
		Status:      models.StatusOpen,
		UserID:      userID,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (s *ticketService) GetTicketsByUserID(userID uint) ([]models.Ticket, error) {
	return s.ticketRepo.FindAllByUserID(userID)
}

func (s *ticketService) GetTicketByID(id, userID uint) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	if ticket.UserID != userID {
		return nil, errors.New("forbidden")
	}

	return ticket, nil
}

func (s *ticketService) UpdateTicketStatus(id, userID uint, status models.TicketStatus) error {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return errors.New("ticket not found")
	}

	if ticket.UserID != userID {
		return errors.New("forbidden")
	}

	// Validate transitions
	// open -> in_progress -> closed
	// closed cannot be reopened
	// open cannot directly become closed
	
	if ticket.Status == models.StatusClosed {
		return errors.New("closed tickets cannot be reopened")
	}

	if ticket.Status == models.StatusOpen {
		if status != models.StatusInProgress {
			return errors.New("open tickets can only move to in_progress")
		}
	} else if ticket.Status == models.StatusInProgress {
		if status != models.StatusClosed {
			return errors.New("in_progress tickets can only move to closed")
		}
	}

	ticket.Status = status
	return s.ticketRepo.Update(ticket)
}
