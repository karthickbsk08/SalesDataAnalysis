package refreshmechanism

import (
	"salesdataanalysis/helpers"
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

type StagingSale struct {
	OrderID         int       `gorm:"column:order_id"`
	ProductID       string    `gorm:"column:product_id;size:20;not null"`
	CustomerID      string    `gorm:"column:customer_id;size:20;not null"`
	ProductName     string    `gorm:"column:product_name;size:255;not null"`
	Category        string    `gorm:"column:category;size:100;not null"`
	Region          string    `gorm:"column:region;size:100;not null"`
	DateOfSale      time.Time `gorm:"column:date_of_sale;not null"`
	QuantitySold    int       `gorm:"column:quantity_sold;not null"`
	UnitPrice       float64   `gorm:"column:unit_price;type:numeric(10,2);not null"`
	Discount        float64   `gorm:"column:discount;type:numeric(5,2);not null"`
	ShippingCost    float64   `gorm:"column:shipping_cost;type:numeric(10,2);not null"`
	PaymentMethod   string    `gorm:"column:payment_method;size:50;not null"`
	CustomerName    string    `gorm:"column:customer_name;size:255;not null"`
	CustomerEmail   string    `gorm:"column:customer_email;size:255;not null"`
	CustomerAddress string    `gorm:"column:customer_address;type:text;not null"`
}

type RefreshLogActivity struct {
	RequestID            string                `gorm:"column:request_id;primaryKey" json:"request_id"`
	Status               string                `gorm:"column:status;type:varchar(20);not null" json:"status"`
	TotalRecordsAffected int64                 `gorm:"column:total_records_affected" json:"total_records_affected"`
	ErrorMessage         string                `gorm:"column:error_message;type:text" json:"error_message"`
	CreatedAt            time.Time             `gorm:"column:created_time;autoCreateTime" json:"created_time"`
	UpdatedAt            time.Time             `gorm:"column:updated_time;autoUpdateTime" json:"updated_time"`
	CreatedBy            string                `gorm:"column:created_by;type:varchar(200)" json:"created_by"`
	UpdatedBy            string                `gorm:"column:updated_by;type:varchar(200)" json:"updated_by"`
	RefreshType          string                `gorm:"column:refresh_type;type:varchar(50)" json:"refresh_type"`
	DurationSeconds      int                   `gorm:"column:duration_seconds" json:"duration_seconds"`
	pDebug               *helpers.HelperStruct `gorm:"-"`
}
