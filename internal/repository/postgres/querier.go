// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package postgres

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteOrder(ctx context.Context, id int64) error
	DeleteOrderItem(ctx context.Context, id int64) error
	DeletePayment(ctx context.Context, id int64) error
	DeleteProduct(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetOrder(ctx context.Context, id int64) (Order, error)
	GetOrderItem(ctx context.Context, id int64) (OrderItem, error)
	GetPayment(ctx context.Context, id int64) (Payment, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	GetUser(ctx context.Context, id int64) (User, error)
	ListOrderItems(ctx context.Context) ([]OrderItem, error)
	ListOrderItemsByOrder(ctx context.Context, orderID int64) ([]OrderItem, error)
	ListOrderItemsByProduct(ctx context.Context, productID int64) ([]OrderItem, error)
	ListOrders(ctx context.Context) ([]Order, error)
	ListPayments(ctx context.Context) ([]Payment, error)
	ListProducts(ctx context.Context) ([]Product, error)
	ListUsers(ctx context.Context) ([]User, error)
	SearchOrdersByStatus(ctx context.Context, status OrderStatus) ([]Order, error)
	SearchOrdersByUser(ctx context.Context, userID int64) ([]Order, error)
	SearchPaymentsByOrder(ctx context.Context, orderID int64) ([]Payment, error)
	SearchPaymentsByStatus(ctx context.Context, status PaymentStatus) ([]Payment, error)
	SearchPaymentsByUser(ctx context.Context, userID int64) ([]Payment, error)
	SearchProductsByCategory(ctx context.Context, category string) ([]Product, error)
	SearchProductsByName(ctx context.Context, dollar_1 sql.NullString) ([]Product, error)
	SearchUsersByEmail(ctx context.Context, email string) ([]User, error)
	SearchUsersByName(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) (OrderItem, error)
	UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (Payment, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) (Product, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
