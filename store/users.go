package store

import (
	"context"
	"fmt"

	"coursechain/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) PutUser(ctx context.Context, user domain.User) error {
	if err := ensureID(user.ID); err != nil {
		return err
	}
	if user.Name == "" || user.Email == "" {
		return fmt.Errorf("name and email are required")
	}
	data, err := encode(user)
	if err != nil {
		return err
	}
	return s.Update(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["User"]).Put(keyFor(user.ID), data)
	})
}

func (s *Store) GetUser(ctx context.Context, id string) (domain.User, error) {
	if err := ensureID(id); err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err := s.View(ctx, func(tx *bolt.Tx) error {
		value := tx.Bucket(bucketNames["User"]).Get(keyFor(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &user)
	})
	return user, err
}

func (s *Store) ListUsers(ctx context.Context) ([]domain.User, error) {
	users := make([]domain.User, 0)
	err := s.View(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["User"]).ForEach(func(_, value []byte) error {
			var user domain.User
			if err := decode(value, &user); err != nil {
				return err
			}
			users = append(users, user)
			return nil
		})
	})
	return users, err
}

func (s *Store) SetUserActive(ctx context.Context, id string, active bool) error {
	user, err := s.GetUser(ctx, id)
	if err != nil {
		return err
	}
	user.Active = active
	return s.PutUser(ctx, user)
}
