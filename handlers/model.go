package handlers

import (
	"time"
)

// Category represents the categories table
type Category struct {
	CategoryID uint   `gorm:"primaryKey;column:category_id"`
	Name       string `gorm:"column:name;type:varchar(100);unique;not null"`
}

// Region represents the regions table
type Region struct {
	RegionID uint   `gorm:"primaryKey;column:region_id"`
	Name     string `gorm:"column:name;type:varchar(100);unique;not null"`
}

// Product represents the products table
type Product struct {
	ProductID  string  `gorm:"primaryKey;column:product_id;type:varchar(20)"`
	Name       string  `gorm:"column:name;type:varchar(255);not null"`
	CategoryID uint    `gorm:"column:category_id;not null"`
	UnitPrice  float64 `gorm:"column:unit_price;type:numeric(10,2);not null"`

	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID"`
}

// Customer represents the customers table
type Customer struct {
	CustomerID string `gorm:"primaryKey;column:customer_id;type:varchar(20)"`
	Name       string `gorm:"column:name;type:varchar(255);not null"`
	Email      string `gorm:"column:email;type:varchar(255);unique;not null"`
	Address    string `gorm:"column:address;type:text;not null"`
}

// Order represents the orders table
type Order struct {
	OrderID       int       `gorm:"primaryKey;column:order_id"`
	ProductID     string    `gorm:"column:product_id;type:varchar(20);not null"`
	CustomerID    string    `gorm:"column:customer_id;type:varchar(20);not null"`
	RegionID      uint      `gorm:"column:region_id;not null"`
	DateOfSale    time.Time `gorm:"column:date_of_sale;type:date;not null"`
	QuantitySold  int       `gorm:"column:quantity_sold;not null"`
	Discount      float64   `gorm:"column:discount;type:numeric(5,2);default:0.00"`
	ShippingCost  float64   `gorm:"column:shipping_cost;type:numeric(10,2);default:0.00"`
	PaymentMethod string    `gorm:"column:payment_method;type:varchar(50);not null"`

	Product  Product  `gorm:"foreignKey:ProductID;references:ProductID"`
	Customer Customer `gorm:"foreignKey:CustomerID;references:CustomerID"`
	Region   Region   `gorm:"foreignKey:RegionID;references:RegionID"`
}
type CustomerDetails struct {
	TotalNoofCustomers int     `json:"Total_no_of_Customers" gorm:"column:total_no_of_customers"`
	TotalOrders        int     `json:"Total_orders_cnt" gorm:"column:total_orders_cnt"`
	AverageOrderValue  float64 `json:"average_order_value" gorm:"column:average_order_value"`
	Status             string  `json:"status" gorm:"column:status"`
	ErrMsg             string  `json:"errMsg" gorm:"column:errMsg"`
}
