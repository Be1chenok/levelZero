package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Be1chenok/levelZero/internal/domain"
)

type Order interface {
	AddOrder(ctx context.Context, order domain.Order) error
	FindAllOrders(ctx context.Context) ([]domain.Order, error)
	FindOrderByUID(ctx context.Context, orderUID string) (domain.Order, error)
	FindDeliveryByOrderUID(ctx context.Context, orderUID string) (domain.Delivery, error)
	FindPaymentByOrderUID(ctx context.Context, orderUID string) (domain.Payment, error)
	FindItemsByOrderUID(ctx context.Context, orderUID string) ([]domain.Item, error)
}

type order struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) Order {
	return &order{
		db: db,
	}
}

func wrapCommitError(err, e error) error {
	return fmt.Errorf("commit tx: %w:%w", err, e)
}

func wrapBeginError(err error) error {
	return fmt.Errorf("begin tx: %w", err)
}

func wrapRollbackError(err, e error) error {
	return fmt.Errorf("rollback tx: %w:%w", err, e)
}

func (o order) AddOrder(ctx context.Context, order domain.Order) (err error) {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to open transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = wrapRollbackError(err, e)

				return
			}

			return
		}

		if e := tx.Commit(); e != nil {
			err = wrapCommitError(err, e)
		}
	}()

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO orders (
		uid,
		track_number,
		entry,
		locale,
		internal_signature,
		customer_id,
		delivery_service,
		shardkey,
		sm_id,
		date_created,
		oof_shard
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.UID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	); err != nil {
		return fmt.Errorf("failed to insert data into orders table: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO deliveries (
		order_uid,
		name,
		phone,
		zip,
		city,
		address,
		region,
		email
		) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.UID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	); err != nil {
		return fmt.Errorf("failed to insert data into deliveries table: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO payments (
		order_uid,
		transaction,
		request_id,
		currency,
		provider,
		amount,
		payment_dt,
		bank,
		delivery_cost,
		goods_total,
		custom_fee
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.UID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	); err != nil {
		return fmt.Errorf("failed to insert data into payments table: %w", err)
	}

	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO items (
		order_uid,
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
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL query: %w", err)
	}
	defer stmt.Close()

	for _, item := range order.Items {
		if _, err := stmt.ExecContext(
			ctx,
			order.UID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		); err != nil {
			return fmt.Errorf("failed to insert data into items table: %w", err)
		}
	}

	return nil
}

func (o order) FindAllOrders(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	rows, err := o.db.QueryContext(
		ctx,
		`SELECT
		uid,
		track_number,
		entry,
		locale,
		internal_signature,
		customer_id,
		delivery_service,
		shardkey,
		sm_id,
		date_created,
		oof_shard
		FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(
			&order.UID,
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
			return nil, domain.ErrNothingFound
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterating over rows: %w", err)
	}

	for idx := range orders {
		delivery, err := o.FindDeliveryByOrderUID(ctx, orders[idx].UID)
		if err != nil {
			return nil, fmt.Errorf("failed to find deliveries by orderUID: %w", err)
		}

		payment, err := o.FindPaymentByOrderUID(ctx, orders[idx].UID)
		if err != nil {
			return nil, fmt.Errorf("failed to find payments by orderUID: %w", err)
		}

		items, err := o.FindItemsByOrderUID(ctx, orders[idx].UID)
		if err != nil {
			return nil, fmt.Errorf("failed to find items by orderUID: %w", err)
		}

		orders[idx].Delivery = delivery
		orders[idx].Payment = payment
		orders[idx].Items = items
	}

	return orders, nil
}

func (o order) FindOrderByUID(ctx context.Context, orderUID string) (domain.Order, error) {
	var order domain.Order

	if err := o.db.QueryRowContext(
		ctx,
		`SELECT
		uid,
		track_number,
		entry,
		locale,
		internal_signature,
		customer_id,
		delivery_service,
		shardkey,
		sm_id,
		date_created,
		oof_shard
		FROM orders
		WHERE uid=$1`,
		orderUID).Scan(
		&order.UID,
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
		return domain.Order{}, domain.ErrNothingFound
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
		address,
		region,
		email
		FROM deliveries
		WHERE order_uid=$1`,
		orderUID).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	); err != nil {
		return domain.Delivery{}, domain.ErrNothingFound
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
		delivery_cost,
		goods_total,
		custom_fee
		FROM payments
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
		return domain.Payment{}, domain.ErrNothingFound
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
		WHERE order_uid=$1
		ORDER BY id ASC`,
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
			return nil, domain.ErrNothingFound
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterating over rows: %w", err)
	}

	return items, nil
}
