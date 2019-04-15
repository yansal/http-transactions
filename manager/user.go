package manager

import (
	"context"
	"database/sql"

	"github.com/yansal/http-transactions/model"
	"github.com/yansal/http-transactions/payload"
)

type User struct{ db *sql.DB }

func NewUser(db *sql.DB) *User { return &User{db: db} }

func (m *User) CreateUser(ctx context.Context, p *payload.User) error {
	return Transaction(ctx, m.db, func(ctx context.Context, tx Queryer) error {
		user, err := FindUser(ctx, tx, p.Email)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if user != nil {
			return nil
		}

		user, err = CreateUser(ctx, tx, &model.User{Email: p.Email})
		if err != nil {
			return err
		}

		_, err = CreateUserAction(ctx, tx, &model.UserAction{UserID: user.ID, Action: "created"})
		if err != nil {
			return err
		}
		return nil
	})
}

func FindUser(ctx context.Context, db Queryer, email string) (*model.User, error) {
	var id int64
	err := db.QueryRowContext(ctx,
		`SELECT id, email FROM users WHERE email = $1`,
		email).
		Scan(&id, &email)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, Email: email}, err
}

func CreateUser(ctx context.Context, db Queryer, user *model.User) (*model.User, error) {
	err := db.QueryRowContext(ctx,
		`INSERT INTO users(email) VALUES($1) RETURNING id, email`,
		user.Email).
		Scan(&user.ID, &user.Email)
	return user, err
}

func CreateUserAction(ctx context.Context, db Queryer, useraction *model.UserAction) (*model.UserAction, error) {
	err := db.QueryRowContext(ctx,
		`INSERT INTO user_actions(user_id, action) VALUES($1, $2) RETURNING id, user_id, action, occurred_at`,
		useraction.UserID, useraction.Action).
		Scan(&useraction.ID, &useraction.UserID, &useraction.Action, &useraction.OccurredAt)
	return useraction, err
}
