package postgresql

import (
	"activator/internal"
	"activator/internal/app/models"
	"database/sql"
	"log/slog"
)

type TokenRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewTokenRepo(db *sql.DB, logger *slog.Logger) *TokenRepository {
	return &TokenRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TokenRepository) Save(t models.Token) error {
	result, err := r.db.Exec(
		`INSERT INTO public.tokens 
		    (hash, user_id, expiry) 
		VALUES 
		    ($1, $2, $3);`,
		t.Hash[:], t.UserID, t.Expiry)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "token repo create")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "rows affected")

	}

	if rows != 1 {
		return internal.NewErrorf(internal.ErrorCodeUnknown, "row was not inserted")
	}

	return nil
}

func (r *TokenRepository) DeleteAll(id int32) error {
	result, err := r.db.Exec(
		`DELETE FROM 
		    public.tokens 
		WHERE 
		    user_id = $1;`,
		id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "token repo DeleteAll")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "rows affected")

	}

	if rows < 1 {
		return internal.NewErrorf(internal.ErrorCodeUnknown, "rows were not deleted")
	}

	return nil
}
