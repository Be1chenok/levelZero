package postgres

import (
	"context"
	"fmt"

	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Order interface {
	FindOrderByUID(ctx context.Context, orderUID string) (domain.Order, error)
	FindDeliveryByOrderUID(ctx context.Context, orderUID string) (domain.Delivery, error)
	FindPaymentByOrderUID(ctx context.Context, orderUID string) (domain.Payment, error)
	FindItemsByOrderUID(ctx context.Context, orderUID string) ([]domain.Item, error)
}

type order struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) Order {
	return &order{
		db: db,
	}
}

func (o order) FindOrderByUID(ctx context.Context, orderUID string) (domain.Order, error) {
	var order domain.Order

	if err := o.db.QueryRowContext(
		ctx,
		`SELECT
		order_uid,
		track_number,
		entry,
		locale,
		internal_signature,
		customer_id,
		delivery_service,
		shardkey, o.sm_id,
		data_created,
		off_shard
		FROM orders
		WHERE order_uid=$1`,
		orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	); err != nil {
		return domain.Order{}, fmt.Errorf("failed to scan row: %w", err)
	}

	return order, nil
}

func (o order) FindDeliveryByOrderUID(ctx context.Context, orderUID string) (domain.Delivery, error) {
	var delivery domain.Delivery

	if err := o.db.QueryRowContext(
		ctx,
		`SELECT
		name,
		phone,
		zip,
		city,
		addres,
		region,
		email
		FROM delivery
		WHERE order_uid=$1`,
		orderUID).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Addres,
		&delivery.Region,
		&delivery.Email,
	); err != nil {
		return domain.Delivery{}, fmt.Errorf("failed to scan row: %w", err)
	}
	return delivery, nil
}

func (o order) FindPaymentByOrderUID(ctx context.Context, orderUID string) (domain.Payment, error) {
	var payment domain.Payment

	if err := o.db.QueryRowContext(
		ctx,
		`SELECT
		transaction,
		request_id,
		currency,
		provider,
		amount,
		payment_dt,
		bank,
		delyvery_cost,
		goods_total,
		custom_fee
		FROM payment
		WHERE order_uid=$1`,
		orderUID).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	); err != nil {
		return domain.Payment{}, fmt.Errorf("failed to scan row: %w", err)
	}

	return payment, nil
}

func (o order) FindItemsByOrderUID(ctx context.Context, orderUID string) ([]domain.Item, error) {
	var items []domain.Item

	rows, err := o.db.QueryContext(
		ctx,
		`SELECT
		chrt_id,
		track_number,
		price,
		rid,
		name,
		sale,
		size,
		total_price,
		nm_id,
		brand,
		status
		FROM items
		WHERE order_uid=$1`,
		orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterating over rows: %w", err)
	}

	return items, nil
}
