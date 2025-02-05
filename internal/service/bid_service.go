package service

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type BidService interface {
	CreateBid(ctx context.Context, bid *model.CreateBidRequest) (*model.Bid, error)
	GetCurrentUserBids(ctx context.Context, limit int, offset int, username string) ([]model.Bid, error)
	GetTenderBids(ctx context.Context, tenderID string, limit int, offset int, username string) ([]model.Bid, error)
	GetBidStatus(ctx context.Context, bidID string, username string) (model.BidStatus, error)
	UpdateBidStatus(ctx context.Context, bidID string, username string, status string) (model.BidStatus, error)
	EditBid(ctx context.Context, bidID string, username string, updateData model.UpdateData) (*model.Bid, error)
	SubmitBidDecision(ctx context.Context, bidID string, username string, decision string) (*model.Bid, error)
	RollbackBidVersion(ctx context.Context, bidID string, username string, version int) (*model.Bid, error)
}

type bidService struct {
	BidRepository          repository.BidRepository
	tenderRepository       repository.TenderRepository
	organizationRepository repository.OrganizationRepository
	userRepository         repository.UserRepository
	logger                 *slog.Logger
}

func NewBidService(bidRepository repository.BidRepository, tenderRepository repository.TenderRepository, organizationRepository repository.OrganizationRepository, userRepository repository.UserRepository, logger *slog.Logger) BidService {
	return &bidService{bidRepository, tenderRepository, organizationRepository, userRepository, logger}
}

func (s *bidService) CreateBid(ctx context.Context, bidRequest *model.CreateBidRequest) (*model.Bid, error) {

	_, err := s.tenderRepository.GetTenderById(ctx, bidRequest.TenderID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tender", slog.Any("error", err))
		if errors.Is(err, model.ErrTenderNotFound) {
			return nil, model.ErrTenderNotFound
		}
		return nil, fmt.Errorf("Error getting tender, %w", err)
	}

	var authorID string
	authorType := model.BidAuthorTypeUser
	if bidRequest.OrganizationID != "" {
		authorType = model.BidAuthorTypeOrganization
		authorID = bidRequest.OrganizationID
		_, err = s.organizationRepository.GetOrganizationById(ctx, bidRequest.OrganizationID)
		if err != nil {
			s.logger.ErrorContext(ctx, "Error getting organization", slog.Any("error", err))
			return nil, fmt.Errorf("Error getting organization, %w", err)
		}
	} else {
		_, err := s.userRepository.GetUserByUsername(ctx, bidRequest.CreatorUsername)
		if err != nil {
			s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
			if errors.Is(err, model.ErrUserNotFound) {
				return nil, model.ErrUserNotFound
			}
			return nil, fmt.Errorf("Error getting user, %w", err)
		}
		authorID = bidRequest.CreatorUsername
	}

	bid := &model.Bid{}

	bid.ID = uuid.NewString()
	bid.Name = bidRequest.Name
	bid.Description = bidRequest.Description
	bid.Status = bidRequest.Status
	bid.TenderID = bidRequest.TenderID
	bid.AuthorType = authorType
	bid.AuthorID = authorID
	bid.CreatorUsername = bidRequest.CreatorUsername
	bid.Version = 1
	bid.CreatedAt = time.Now()
	bid.UpdatedAt = time.Now()

	bidResponse, err := s.BidRepository.CreateBid(ctx, bid)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error creating bid", slog.Any("error", err))
		return nil, fmt.Errorf("Error creating bid, %w", err)
	}

	return bidResponse, nil

}

func (s *bidService) GetCurrentUserBids(ctx context.Context, limit int, offset int, username string) ([]model.Bid, error) {
	_, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error getting user: %w", err)
	}

	bid, err := s.BidRepository.GetBidByUsername(ctx, limit, offset, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
		return nil, fmt.Errorf("Error getting bids, %w", err)
	}
	return bid, nil
}

func (s *bidService) GetTenderBids(ctx context.Context, tenderID string, limit int, offset int, username string) ([]model.Bid, error) {

	_, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error getting user, %w", err)
	}

	_, err = s.tenderRepository.GetTenderById(ctx, tenderID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tender", slog.Any("error", err))
		if errors.Is(err, model.ErrTenderNotFound) {
			return nil, model.ErrTenderNotFound
		}
		return nil, fmt.Errorf("Error getting tender, %w", err)
	}

	bids, err := s.BidRepository.GetTenderBids(ctx, tenderID, limit, offset, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
		if errors.Is(err, model.ErrBidNotFound) {
			return nil, model.ErrBidNotFound
		}
		return nil, fmt.Errorf("Error getting bids, %w", err)
	}
	return bids, nil
}

func (s *bidService) GetBidStatus(ctx context.Context, bidID string, username string) (model.BidStatus, error) {

	_, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrUserNotFound
		}
		return "", fmt.Errorf("Error getting user: %w", err)
	}

	status, err := s.BidRepository.GetBidStatus(ctx, bidID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting bid status", slog.Any("error", err))
		if errors.Is(err, model.ErrBidNotFound) {
			return "", model.ErrBidNotFound
		}
		return "", fmt.Errorf("Error getting bid status, %w", err)
	}
	return status, nil
}

func (s *bidService) UpdateBidStatus(ctx context.Context, bidID string, username string, status string) (model.BidStatus, error) {
	_, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrUserNotFound
		}
		return "", fmt.Errorf("Error getting user: %w", err)
	}

	bid, err := s.BidRepository.GetBidById(ctx, bidID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting bid", slog.Any("error", err))
		if errors.Is(err, model.ErrBidNotFound) {
			return "", model.ErrBidNotFound
		}
		return "", fmt.Errorf("Error getting bid, %w", err)
	}

	isResponsible, err := s.organizationRepository.IsUserResponsibleForOrganization(ctx, bid.AuthorID, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting responsible for organization", slog.Any("error", err))
		return "", fmt.Errorf("Error getting responsible for organization, %w", err)
	}
	if !isResponsible {
		s.logger.ErrorContext(ctx, "User is not responsible for organization", slog.Any("error", err))
		return "", model.ErrForbidden
	}

	bid.Status = model.BidStatus(status)

	bid, err = s.BidRepository.UpdateBid(ctx, bid)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error updating bid", slog.Any("error", err))
		return "", fmt.Errorf("Error updating bid, %w", err)
	}
	return bid.Status, nil

}

func (s *bidService) EditBid(ctx context.Context, bidID string, username string, updateData model.UpdateData) (*model.Bid, error) {
	_, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error getting user: %w", err)
	}

	bid, err := s.BidRepository.GetBidById(ctx, bidID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting bid", slog.Any("error", err))
		if errors.Is(err, model.ErrBidNotFound) {
			return nil, model.ErrBidNotFound
		}
		return nil, fmt.Errorf("Error getting bid, %w", err)
	}

	if bid.CreatorUsername != username {
		s.logger.ErrorContext(ctx, "User is not responsible for bid", slog.Any("error", err))
		return nil, model.ErrForbidden
	}

	if updateData.Name != nil {
		if *updateData.Name != "" {
			bid.Name = *updateData.Name
		}
	}

	if updateData.Description != nil {
		if *updateData.Description != "" {
			bid.Description = *updateData.Description
		}
	}

	bid, err = s.BidRepository.UpdateBid(ctx, bid)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error updating bid", slog.Any("error", err))
		return nil, fmt.Errorf("Error updating bid, %w", err)
	}
	return bid, nil
}

func (s *bidService) RollbackBidVersion(ctx context.Context, bidID string, username string, version int) (*model.Bid, error) {
	_, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error getting user: %w", err)
	}

	bid, err := s.BidRepository.GetBidById(ctx, bidID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting bid", slog.Any("error", err))
		if errors.Is(err, model.ErrBidNotFound) {
			return nil, model.ErrBidNotFound
		}
		return nil, fmt.Errorf("Error getting bid, %w", err)
	}

	if bid.CreatorUsername != username {
		s.logger.ErrorContext(ctx, "User is not responsible for bid", slog.Any("error", err))
		return nil, model.ErrForbidden
	}

	bid, err = s.BidRepository.RollbackBidVersion(ctx, bidID, version)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error rolling back bid version", slog.Any("error", err))
		return nil, fmt.Errorf("Error rolling back bid version, %w", err)
	}
	return bid, nil
}

func (s *bidService) SubmitBidDecision(ctx context.Context, bidID string, username string, decision string) (*model.Bid, error) {
	_, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error getting user: %w", err)
	}

	bid, err := s.BidRepository.GetBidById(ctx, bidID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting bid", slog.Any("error", err))
		if errors.Is(err, model.ErrBidNotFound) {
			return nil, model.ErrBidNotFound
		}
		return nil, fmt.Errorf("Error getting bid, %w", err)
	}

	tender, err := s.tenderRepository.GetTenderById(ctx, bid.TenderID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting tender", slog.Any("error", err))
		if errors.Is(err, model.ErrTenderNotFound) {
			return nil, model.ErrTenderNotFound
		}
		return nil, fmt.Errorf("Error getting tender, %w", err)
	}

	isResponsible, err := s.tenderRepository.IsUserResponsibleForTender(ctx, tender.ID, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error checking user responsibility for tender", slog.Any("error", err))
		return nil, fmt.Errorf("Error checking user responsibility for tender: %w", err)
	}
	if !isResponsible {
		s.logger.ErrorContext(ctx, "User is not responsible for this tender", slog.String("username", username), slog.String("tenderID", tender.ID))
		return nil, model.ErrForbidden
	}

	if bid.Status != model.BidStatusPublished && bid.Status != model.BidStatusCreated {
		s.logger.ErrorContext(ctx, "Cannot submit decision for bid with status", slog.String("status", string(bid.Status)))
		return nil, model.ErrDecisionSubmit
	}

	if decision == "Approved" {
		bid.Status = model.BidStatusApproved
		tender.Status = model.TenderStatusClosed
		_, err = s.tenderRepository.UpdateTender(ctx, tender)
		if err != nil {
			s.logger.ErrorContext(ctx, "Error updating tender status to closed", slog.Any("error", err))
			return nil, fmt.Errorf("Error updating tender status to closed: %w", err)
		}

	} else if decision == "Rejected" {
		bid.Status = model.BidStatusRejected
	} else {
		s.logger.ErrorContext(ctx, "Invalid decision parameter", slog.String("decision", decision))
		return nil, model.ErrDecisionSubmit
	}

	updatedBid, err := s.BidRepository.UpdateBid(ctx, bid)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error updating bid status on decision", slog.Any("error", err))
		return nil, fmt.Errorf("Error updating bid status on decision: %w", err)
	}

	return updatedBid, nil
}
