package service

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type BidService interface {
	CreateBid(ctx context.Context, bid *model.CreateBidRequest) (*model.Bid, error)
	GetCurrentUserBids(ctx context.Context, limit int, offset int, username string) ([]model.Bid, error)
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
		_, err := s.userRepository.GetUserById(ctx, bidRequest.CreatorUsername)
		if err != nil {
			s.logger.ErrorContext(ctx, "Error getting user", slog.Any("error", err))
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
	bid, err := s.BidRepository.GetBidByUsername(ctx, limit, offset, username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
		return nil, fmt.Errorf("Error getting bids, %w", err)
	}
	return bid, nil
}
