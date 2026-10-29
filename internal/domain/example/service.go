package example

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"golang-starter-kit/pkg/cache"
	"golang-starter-kit/pkg/mail"
)

const (
	CacheTagCruds = "cruds"
	CacheKeyAll   = "cruds:all"
)

type Service struct {
	repo  Repository
	cache cache.Cache
	mail  *mail.MailClient
}

func NewService(repo Repository, c cache.Cache, m *mail.MailClient) *Service {
	return &Service{
		repo:  repo,
		cache: c,
		mail:  m,
	}
}

func (s *Service) List(ctx context.Context) ([]Crud, error) {
	if s.cache != nil {
		var cached []Crud
		err := s.cache.RememberWithTags(ctx, CacheKeyAll, 5*time.Minute, []string{CacheTagCruds}, &cached, func() (any, error) {
			return s.repo.FindAll(ctx)
		})
		if err == nil {
			return cached, nil
		}
	}

	return s.repo.FindAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Crud, error) {
	cacheKey := fmt.Sprintf("cruds:%s", id)
	if s.cache != nil {
		var item Crud
		err := s.cache.RememberWithTags(ctx, cacheKey, 5*time.Minute, []string{CacheTagCruds}, &item, func() (any, error) {
			return s.repo.FindByID(ctx, id)
		})
		if err == nil {
			return &item, nil
		}
	}

	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateCrudRequest) (*Crud, error) {
	item := &Crud{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      "active",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.InvalidateTag(ctx, CacheTagCruds)
	}

	return item, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateCrudRequest) (*Crud, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	item.Title = req.Title
	item.Description = req.Description
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.InvalidateTag(ctx, CacheTagCruds)
	}

	return item, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.cache != nil {
		_ = s.cache.InvalidateTag(ctx, CacheTagCruds)
	}

	return nil
}

func (s *Service) SendExampleEmail(ctx context.Context, email, name string) error {
	if s.mail == nil {
		return fmt.Errorf("mail client is not configured")
	}

	welcomeMail := mail.NewWelcomeMail(name)
	return s.mail.To(email).Queue(ctx, welcomeMail)
}
