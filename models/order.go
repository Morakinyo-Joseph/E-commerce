package models

import (
	"time"
)

// Order represents a purchase transaction
type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null" json:"user_id"`
	User        User        `gorm:"foreignKey:UserID" json:"-"`
	Status      string      `gorm:"type:varchar(20);default:'Pending'" json:"status"` // 'Pending', 'Completed', 'Cancelled'
	TotalAmount float64     `gorm:"not null" json:"total_amount"`
	CreatedAt   time.Time   `json:"created_at"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}

// OrderItem represents a single product within an order
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null" json:"order_id"`
	Order     Order   `gorm:"foreignKey:OrderID" json:"-"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"-"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"` // Price of the product at the time of purchase
}
