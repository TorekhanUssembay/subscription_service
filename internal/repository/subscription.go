package repository

import (
    "context"
    "fmt"
    "time"

    "github.com/TorekhanUssembay/subscription_service/internal/model"

    "github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
    db *pgxpool.Pool
}

func NewSubscriptionRepo(db *pgxpool.Pool) *SubscriptionRepo {
    return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
    query := `
        INSERT INTO subscriptions (
            service_name, price, user_id, start_date, end_date
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
    return r.db.QueryRow(
        ctx,
        query,
        sub.ServiceName,
        sub.Price,
        sub.UserID,
        sub.StartDate,
        sub.EndDate,
    ).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
}

func (r *SubscriptionRepo) GetSubscription(ctx context.Context, id string) (*model.Subscription, error) {
    query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        WHERE id = $1
    `

    sub := &model.Subscription{}
    err := r.db.QueryRow(ctx, query, id).Scan(
        &sub.ID,
        &sub.ServiceName,
        &sub.Price,
        &sub.UserID,
        &sub.StartDate,
        &sub.EndDate,
        &sub.CreatedAt,
        &sub.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("Failed to get subscription: %w", err)
    }

    return sub, nil
}

func (r *SubscriptionRepo) UpdateSubscription(ctx context.Context, sub *model.Subscription) error {
    query := `
        UPDATE subscriptions
        SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = NOW()
        WHERE id = $6
    `
    _, err := r.db.Exec(ctx, query,
        sub.ServiceName,
        sub.Price,
        sub.UserID,
        sub.StartDate,
        sub.EndDate,
        sub.ID,
    )
    if err != nil {
        return fmt.Errorf("Failed to update subscription: %w", err)
    }
    return nil
}

func (r *SubscriptionRepo) DeleteSubscription(ctx context.Context, id string) error {
    _, err := r.db.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
    if err != nil {
        return fmt.Errorf("Failed to delete subscription: %w", err)
    }
    return nil
}

func (r *SubscriptionRepo) ListSubscriptions(ctx context.Context, userID string, serviceName *string) ([]*model.Subscription, error) {
    query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at FROM subscriptions WHERE user_id = $1`
    args := []interface{}{userID}

    if serviceName != nil {
        query += " AND service_name = $2"
        args = append(args, *serviceName)
    }

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("Failed to list subscriptions: %w", err)
    }
    defer rows.Close()

    var subs []*model.Subscription
    for rows.Next() {
        sub := &model.Subscription{}
        if err := rows.Scan(
            &sub.ID,
            &sub.ServiceName,
            &sub.Price,
            &sub.UserID,
            &sub.StartDate,
            &sub.EndDate,
            &sub.CreatedAt,
            &sub.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("Failed to scan subscription: %w", err)
        }
        subs = append(subs, sub)
    }
    return subs, nil
}

func (r *SubscriptionRepo) SumSubscriptions(ctx context.Context, userID string, serviceName *string, from, to time.Time) (int, error) {
    query := `SELECT COALESCE(SUM(price),0) FROM subscriptions WHERE user_id = $1 AND start_date >= $2 AND start_date <= $3`
    args := []interface{}{userID, from, to}

    if serviceName != nil {
        query += " AND service_name = $4"
        args = append(args, *serviceName)
    }

    var sum int
    err := r.db.QueryRow(ctx, query, args...).Scan(&sum)
    if err != nil {
        return 0, fmt.Errorf("Failed to sum subscriptions: %w", err)
    }
    return sum, nil
}