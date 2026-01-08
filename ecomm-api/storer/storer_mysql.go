package storer

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MySQLStorer struct {
	db *sqlx.DB
}

func NewMySQLStorer(db *sqlx.DB) *MySQLStorer {
	return &MySQLStorer{db: db}
}

func (ms *MySQLStorer) CreateProduct(ctx context.Context, p *Product) (*Product, error) {
	res, err := ms.db.NamedExecContext(ctx, "INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)", p)
	if err != nil {
		return nil, fmt.Errorf("Error inserting product: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Error while getting the last product id: %w", err)
	}

	p.ID = id

	return p, nil
}

func (ms *MySQLStorer) GetProduct(ctx context.Context, id int64) (*Product, error) {
	var p Product
	err := ms.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get product: %w", err)
	}

	return &p, nil
}

func (ms *MySQLStorer) ListProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product
	err := ms.db.SelectContext(ctx, &products, "SELECT * FROM products")
	if err != nil {
		return nil, fmt.Errorf("Failed to load products: %w", err)
	}

	return products, nil
}

func (ms *MySQLStorer) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	_, err := ms.db.NamedExecContext(ctx, "UPDATE products SET name = :name, image = :image, category = :category, description = :description, rating = :rating, num_reviews = :num_reviews, price = :price, count_in_stock = :count_in_stock, updated_at = NOW() WHERE id = :id;", p)
	if err != nil {
		return nil, fmt.Errorf("Error updating product: %w", err)
	}

	return p, nil
}

func (ms *MySQLStorer) DeleteProduct(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM products WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("Error deleting product: %w", err)
	}

	return nil
}

func (ms *MySQLStorer) CreateOrder(ctx context.Context, o *Order) (*Order, error) {
	
	//start transaction
		//insert into orders
		//insert into order_items
	//commit transaction
	//rollback transaction if error

	err := ms.execTx(ctx, func(tx *sqlx.Tx) error {

		//insert into orders

		err := createOrder(ctx, tx, o)
		if err != nil {
			return fmt.Errorf("Error creating order: %v", err)
		}

		for _, oi := range o.Items {
			oi.OrderID = o.ID

			//insert into order_items

			err = createOrderItem(ctx, tx, &oi)
			if err != nil {
				return fmt.Errorf("Error creating order item: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error creating order: %v", err)
	}

	return o, nil
}

func createOrder(ctx context.Context, tx *sqlx.Tx, o *Order) error {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, user_id) VALUES (:payment_method, :tax_price, :shipping_price, :total_price, :user_id)", o)
	if err != nil {
		return fmt.Errorf("Error inserting order: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("Error getting the last order id: %v", err)
	}

	o.ID = id

	return nil
}

func createOrderItem(ctx context.Context, tx *sqlx.Tx, oi *OrderItem) error {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (:name, :quantity, :image, :price, :product_id, :order_id)", oi)
	if err != nil {
		return fmt.Errorf("Error while inserting order item: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("Error while finding last insert id: %v", err)
	}

	oi.ID = id

	return nil
}

func (ms *MySQLStorer) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := ms.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Error while starting transaction: %v", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("Error while rolling back transaction: %v", rbErr)
		}
		return fmt.Errorf("Error in transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Error while committing transaction: %v", err)
	}

	return nil
}