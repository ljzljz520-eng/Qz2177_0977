package workflow

import (
	"context"
	"fmt"
	"strings"

	"coursechain/domain"
)

func (s *Service) RegisterUser(ctx context.Context, user domain.User) error {
	user.ID = strings.TrimSpace(user.ID)
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Role = strings.ToLower(strings.TrimSpace(user.Role))
	if user.Role == "" {
		user.Role = "student"
	}
	if user.ID == "" || user.Name == "" || user.Email == "" {
		return fmt.Errorf("id, name and email are required")
	}
	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("email is invalid")
	}
	user.Active = true
	user.CreatedAt = s.clock()
	return s.store.PutUser(ctx, user)
}

func (s *Service) DeactivateUser(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id is required")
	}
	return s.store.SetUserActive(ctx, id, false)
}

func (s *Service) EnsureActor(ctx context.Context, id string) (domain.User, error) {
	user, err := s.store.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	if !user.Active {
		return domain.User{}, fmt.Errorf("actor is inactive")
	}
	return user, nil
}

func (s *Service) Course() string {
	return s.course
}
