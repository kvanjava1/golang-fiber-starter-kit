package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"golang-starter-kit/pkg/cache"
)

const (
	CacheTagUser = "user"
	CacheKeyUserAll = "user:all"
)

type Service struct {
	repo  Repository
	cache cache.Cache
}

func NewService(repo Repository, c cache.Cache) *Service {
	return &Service{
		repo:  repo,
		cache: c,
	}
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	if s.cache != nil {
		var cached []User
		err := s.cache.RememberWithTags(ctx, CacheKeyUserAll, 5*time.Minute, []string{CacheTagUser}, &cached, func() (any, error) {
			return s.repo.FindAll(ctx)
		})
		if err == nil {
			return cached, nil
		}
	}
	return s.repo.FindAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	key := fmt.Sprintf("user:%s", id)
	if s.cache != nil {
		var item User
		err := s.cache.RememberWithTags(ctx, key, 5*time.Minute, []string{CacheTagUser}, &item, func() (any, error) {
			return s.repo.FindByID(ctx, id)
		})
		if err == nil {
			return &item, nil
		}
	}
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	item := &User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Status:    "active",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.InvalidateTag(ctx, CacheTagUser)
	}

	return item, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	item.Name = req.Name
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.InvalidateTag(ctx, CacheTagUser)
	}

	return item, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.cache != nil {
		_ = s.cache.InvalidateTag(ctx, CacheTagUser)
	}

	return nil
}
