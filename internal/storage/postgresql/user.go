package postgresql

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"activator/internal"
	"activator/internal/app/models"
	m "activator/internal/rest/models"
)

type UserRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserRepo(db *sql.DB, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) Create(p m.CreateParams) (models.User, error) {
	var id int32

	if err := r.db.QueryRow(
		`INSERT INTO public.users 
		    (email, user_name, nick_name) 
		VALUES 
		    ($1, $2, $3) 
		RETURNING id;`,
		p.Email, p.Name, p.NickName,
	).Scan(&id); err != nil {
		return models.User{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo create")
	}

	r.logger.Debug("record created", "id", id)

	return models.User{
		ID:        id,
		Email:     p.Email,
		Name:      p.Name,
		NickName:  p.NickName,
		Activated: false,
	}, nil
}

func (r *UserRepository) Delete(id int32) error {
	result, err := r.db.Exec("DELETE FROM public.users WHERE id = $1;", id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo delete")
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo delete")
	}
	if deleted != 1 {
		return internal.NewErrorf(internal.ErrorCodeNotFound, "resourse with id %d not found", id)
	}

	r.logger.Debug("record deleted", "id", id)
	return nil
}

func (r *UserRepository) FindByToken(tokenPlaintext string) (int32, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
	SELECT users.id
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.expiry > $2;`

	var id int32
	if err := r.db.QueryRow(query, tokenHash[:], time.Now()).Scan(&id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, internal.NewErrorf(internal.ErrorCodeNotFound, "resourse with token %s not found", tokenPlaintext)
		default:
			return 0, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "find by token")
		}
	}

	return id, nil
}

func (r *UserRepository) Activate(id int32) error {
	result, err := r.db.Exec("UPDATE public.users SET activated = TRUE WHERE id = $1;", id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo activate")
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "repo activate")
	}
	if updated != 1 {
		return internal.NewErrorf(internal.ErrorCodeNotFound, "resourse with id %d not found", id)
	}

	r.logger.Debug("record activated", "id", id)
	return nil
}

func (r *UserRepository) FindById(id int32) (models.User, error) {
	var user models.User

	query := `
	SELECT 
	    id, email, user_name, nick_name, activated
	FROM 
	    users
	WHERE 
	    id = $1;`

	if err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.NickName, &user.Activated,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, internal.NewErrorf(internal.ErrorCodeNotFound, "resourse with id %d not found", id)
		default:
			return models.User{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "find user by id")
		}
	}

	return user, nil
}
