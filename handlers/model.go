package handlers

type CustomerDetails struct {
	TotalNoofCustomers int     `json:"Total_no_of_Customers" gorm:"column:total_no_of_customers"`
	TotalOrders        int     `json:"Total_orders_cnt" gorm:"column:total_orders_cnt"`
	AverageOrderValue  float64 `json:"average_order_value" gorm:"column:average_order_value"`
	Status             string  `json:"status" gorm:"column:status"`
	ErrMsg             string  `json:"errMsg" gorm:"column:errMsg"`
}

type MonthTrend struct {
	MonthYear  string  `gorm:"column:month_year"`
	MonthTrend float64 `gorm:"column:month_trend"`
}

type YearTrend struct {
	Yearwise  string  `gorm:"column:yearwise"`
	YearTrend float64 `gorm:"column:year_trend"`
}

type YearQuarterTrend struct {
	YearQuarter string  `gorm:"column:year_quater"`
	Total       float64 `gorm:"column:year_quarter"`
}

type TotalRevenue struct {
	TotalRevenue float64 `gorm:"column:totalrevenue"`
}

type ProductRevenue struct {
	ProductName string  `gorm:"column:name"`
	Revenue     float64 `gorm:"column:prod_by_totalrevenue"`
}

type CategoryRevenue struct {
	CategoryName string  `gorm:"column:name"`
	Revenue      float64 `gorm:"column:cat_by_totalrevenue"`
}

type RegionRevenue struct {
	RegionName string  `gorm:"column:name"`
	Revenue    float64 `gorm:"column:region_by_totalrevenue"`
}

type RevenueReport struct {
	MonthlyTrends   []MonthTrend       `json:"monthly_trend"`
	YearlyTrends    []YearTrend        `json:"yearly_trend"`
	QuarterlyTrends []YearQuarterTrend `json:"quarterly_trend"`
	TotalRevenue    TotalRevenue       `json:"total_revenue"`
	ProductWise     []ProductRevenue   `json:"product_wise"`
	CategoryWise    []CategoryRevenue  `json:"category_wise"`
	RegionWise      []RegionRevenue    `json:"region_wise"`
	Status          string             `json:"status"`
	ErrMsg          string             `json:"errMsg"`
}
