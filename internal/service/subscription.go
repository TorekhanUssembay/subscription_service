package service

import (
	"context"
	"fmt"
	"time"

	"github.com/TorekhanUssembay/subscription_service/internal/model"
	"github.com/TorekhanUssembay/subscription_service/internal/repository"
)

type CreateSubscriptionDTO struct {
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int     `json:"price" example:"400"`
	UserID      string  `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	StartDate   string  `json:"start_date" example:"07-2025"`
	EndDate     *string `json:"end_date" example:"12-2025"`
}

type SubscriptionService struct {
	repo *repository.SubscriptionRepo
}

func NewSubscriptionService(repo *repository.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, dto CreateSubscriptionDTO) (*model.Subscription, error) {
    if dto.ServiceName == "" {
        return nil, fmt.Errorf("service_name is required")
    }
    if dto.Price <= 0 {
        return nil, fmt.Errorf("price must be > 0")
    }
    if dto.UserID == "" {
        return nil, fmt.Errorf("user_id is required")
    }
    if dto.StartDate == "" {
        return nil, fmt.Errorf("start_date is required")
    }

    startDate, err := parseMonthYear(dto.StartDate)
    if err != nil {
        return nil, fmt.Errorf("invalid start_date: %w", err)
    }

    var endDate *time.Time
    if dto.EndDate != nil && *dto.EndDate != "" {
        ed, err := parseMonthYear(*dto.EndDate)
        if err != nil {
            return nil, fmt.Errorf("invalid end_date: %w", err)
        }
        endDate = &ed
    }

    sub := &model.Subscription{
        ServiceName: dto.ServiceName,
        Price:       dto.Price,
        UserID:      dto.UserID,
        StartDate:   startDate,
        EndDate:     endDate,
    }

    if err := s.repo.CreateSubscription(ctx, sub); err != nil {
        return nil, fmt.Errorf("Failed to create subscription: %w", err)
    }

    return sub, nil
}

func parseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, id string) (*model.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.GetSubscription(ctx, id)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, id string, dto CreateSubscriptionDTO) (*model.Subscription, error) {
	sub, err := s.repo.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	if dto.ServiceName != "" {
		sub.ServiceName = dto.ServiceName
	}
	if dto.Price > 0 {
		sub.Price = dto.Price
	}
	if dto.UserID != "" {
		sub.UserID = dto.UserID
	}
	if dto.StartDate != "" {
		startDate, err := parseMonthYear(dto.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date: %w", err)
		}
		sub.StartDate = startDate
	}
	if dto.EndDate != nil && *dto.EndDate != "" {
		endDate, err := parseMonthYear(*dto.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		sub.EndDate = &endDate
	}

	if err := s.repo.UpdateSubscription(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return s.repo.DeleteSubscription(ctx, id)
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, userID string, serviceName *string) ([]*model.Subscription, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	return s.repo.ListSubscriptions(ctx, userID, serviceName)
}

func (s *SubscriptionService) SumSubscriptions(ctx context.Context, userID string, serviceName *string, fromStr, toStr string) (int, error) {
	if userID == "" {
		return 0, fmt.Errorf("user_id is required")
	}
	from, err := parseMonthYear(fromStr)
	if err != nil {
		return 0, fmt.Errorf("invalid from: %w", err)
	}
	to, err := parseMonthYear(toStr)
	if err != nil {
		return 0, fmt.Errorf("invalid to: %w", err)
	}
	return s.repo.SumSubscriptions(ctx, userID, serviceName, from, to)
}