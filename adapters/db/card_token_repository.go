package db

import (
	"context"
	"database/sql"
	cardtoken "payment-layer-card-api/entities/card_token"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type CardTokenRepository struct {
	db *sql.DB
}

func NewCardTokenRepository(db *sql.DB) *CardTokenRepository {
	return &CardTokenRepository{
		db: db,
	}
}

func (r *CardTokenRepository) Insert(ctx context.Context, cardToken *cardtoken.CardToken) *errors.ErrorOutput {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO card_tokens (id, card_id, card_token, gateway) VALUES ($1, $2, $3, $4) RETURNING id`,
		cardToken.ID,
		cardToken.CardID,
		cardToken.CardToken,
		cardToken.Gateway,
	).Scan(&cardToken.ID)

	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}

func (r *CardTokenRepository) GetByCardID(ctx context.Context, cardID string) ([]*cardtoken.CardToken, *errors.ErrorOutput) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, card_id, card_token, gateway, created_at, updated_at FROM card_tokens WHERE card_id = $1`,
		cardID,
	)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.NewError(errors.NotFoundError, err)
		}
		return nil, errors.NewError(errors.InternalServerError, err)
	}

	defer rows.Close()

	var cardTokens []*cardtoken.CardToken

	for rows.Next() {
		cardToken := &cardtoken.CardToken{}

		if err := rows.Scan(&cardToken.ID,
			&cardToken.CardID,
			&cardToken.CardToken,
			&cardToken.Gateway,
			&cardToken.CreatedAt,
			&cardToken.UpdatedAt); err != nil {
			return nil, errors.NewError(errors.InternalServerError, err)
		}

		cardTokens = append(cardTokens, cardToken)
	}

	return cardTokens, nil
}

func (r *CardTokenRepository) GetByCardIDAndGateway(ctx context.Context, cardID string, gateway string) (*cardtoken.CardToken, *errors.ErrorOutput) {
	cardToken := &cardtoken.CardToken{}

	err := r.db.QueryRowContext(ctx, `
         SELECT id,
                card_id,
                card_token,
                gateway,
                created_at,
                updated_at
           FROM card_tokens
          WHERE card_id = $1
            AND gateway = $2`,
		cardID,
		gateway).Scan(&cardToken.ID,
		&cardToken.CardID,
		&cardToken.CardToken,
		&cardToken.Gateway,
		&cardToken.CreatedAt,
		&cardToken.UpdatedAt)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.NewError(errors.NotFoundError, err)
		}

		return nil, errors.NewError(errors.InternalServerError, err)
	}

	return cardToken, nil
}

func (r *CardTokenRepository) DeleteByCardToken(ctx context.Context, cardToken string) *errors.ErrorOutput {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM card_tokens WHERE card_token = $1`,
		cardToken,
	)

	if err != nil {
		return errors.NewError(errors.InternalServerError, err)
	}

	return nil
}
