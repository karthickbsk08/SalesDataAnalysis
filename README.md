# Ecommerce SalesDataAnalysis
A Assessment for Lumel organization 


This project processes large-scale sales data from CSV into a normalized  postgresSQL database. It supports both scheduled and on-demand data refreshes with validation, deduplication, and logging. A RESTful API provides endpoints to trigger data updates and calculate key revenue metrics by date, product, category, or region.and customer analysis Designed for scalability and reliable business insights.



## 🛠 Tech Stack

- **Go (Golang)** `v1.23.0` – Backend APIs and CSV processing  
- **GORM** – ORM for seamless database interactions  
- **PostGres (SQL)** – Relational database for storing normalized sales data  
- **JSON API** – For communication between frontend and backend
## 📦 External Packages

- **TOML** – For configuration file parsing  
- **GORM** – ORM for efficient database interaction  
- **Mux** – HTTP router for building RESTful APIs
- **godotenv** – loads environment variables from a .env file.
- **helpers** – Efficient Log file mechanism - Own package.


## 📥 Install Dependencies

```bash
go mod tidy 
```
# 🌐 SalesDataAnalysis API Documentation

Below is a list of available API endpoints.

---

## 🔗 Base URL

## http://localhost:19998/

> Use this base URL to access all SalesDataAnalysis-related endpoints during local development.

---
### 📘 API Endpoints

| Route           | Method | Request Body                                                                                      | Description                                                                                                  |
|----------------|--------|---------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `/customeranalysis` | POST    | ```{"start_date": "2024-01-01","end_date": "2024-12-31"}```                                                                                       | Used to fetch customer analysis on a datarange from the sales data.        |
| `/revenuereport` | POST   | ```{"start_date": "2024-01-01","end_date": "2024-12-31","indicator": "MONTHLYTREND,YEARLYTREND,QUARTERLYTREND,TOTALREV,PRODUCTWISE,CATEGORYWISE,REGIONWISE"}```                                | Retrieves sales revenue data: total revenue, revenue by <ins> product</ins> ,  <ins> by category</ins> , and <ins> by region </ins>  within the date range. |
| `/refreshDataOnDemand` | POST   | No_body                                | on-demand data refreshes with validation, deduplication, and logging. |



### ✅ Success Response Format for `/customeranalysis`

```json
{"Total_no_of_Customers":3,"Total_orders_cnt":5,"average_order_value":933.312,"status":"S","errMsg":""}
```
### ❌ Failure Response Format (Validation or User Error)

```json
{
    "Total_no_of_Customers": 0,
    "Total_orders_cnt": 0,
    "average_order_value": 0,
    "status": "E",
    "errMsg": "handlers/customeranalysis.go @@ salesdataanalysis/handlers.CustomerAnalysis @@ ln 54 @@ ERROR: date/time field value out of range: \"2024-1-0\" (SQLSTATE 22008)"
}
```

### ✅ Success Response Format `/revenuereport` 

`Request Body`
```json
{
    "start_date": "2023-12-15",
    "end_date": "2024-05-18",
    "indicator": "MONTHLYTREND,TOTALREV"
}
```

### total_revenueby `parameters `

`MONTHLYTREND,TOTALREV` in Request Body indicator



```json
{
    "monthly_trend": [
        {
            "MonthYear": "2024-05",
            "MonthTrend": 2488.1
        },
        {
            "MonthYear": "2024-04",
            "MonthTrend": 309.49
        },
        {
            "MonthYear": "2024-03",
            "MonthTrend": 188
        },
        {
            "MonthYear": "2024-02",
            "MonthTrend": 148.98
        },
        {
            "MonthYear": "2024-01",
            "MonthTrend": 1314
        },
        {
            "MonthYear": "2023-12",
            "MonthTrend": 334
        }
    ],
    "yearly_trend": null,
    "quarterly_trend": null,
    "total_revenue": {
        "TotalRevenue": 4782.57
    },
    "product_wise": null,
    "category_wise": null,
    "region_wise": null,
    "status": "S",
    "errMsg": ""
}
```


`MONTHLYTREND,YEARLYTREND,QUARTERLYTREND,TOTALREV,PRODUCTWISE,CATEGORYWISE,REGIONWISE ` in Request Body indicator

```json
{
    "monthly_trend": [
        {
            "MonthYear": "2024-05",
            "MonthTrend": 2488.1
        },
        {
            "MonthYear": "2024-04",
            "MonthTrend": 309.49
        },
        {
            "MonthYear": "2024-03",
            "MonthTrend": 188
        },
        {
            "MonthYear": "2024-02",
            "MonthTrend": 148.98
        },
        {
            "MonthYear": "2024-01",
            "MonthTrend": 1314
        }
    ],
    "yearly_trend": [
        {
            "Yearwise": "2024",
            "YearTrend": 4448.57
        }
    ],
    "quarterly_trend": [
        {
            "YearQuarter": "2024-Q2",
            "Total": 2797.59
        },
        {
            "YearQuarter": "2024-Q1",
            "Total": 1650.98
        }
    ],
    "total_revenue": {
        "TotalRevenue": 4448.57
    },
    "product_wise": [
        {
            "ProductName": "UltraBoost Running Shoes",
            "Revenue": 188
        },
        {
            "ProductName": "Sony WH-1000XM5 Headphones",
            "Revenue": 309.49
        },
        {
            "ProductName": "iPhone 15 Pro",
            "Revenue": 3802.1
        },
        {
            "ProductName": "Levi's 501 Jeans",
            "Revenue": 148.98
        }
    ],
    "category_wise": [
        {
            "CategoryName": "Electronics",
            "Revenue": 4111.59
        },
        {
            "CategoryName": "Shoes",
            "Revenue": 188
        },
        {
            "CategoryName": "Clothing",
            "Revenue": 148.98
        }
    ],
    "region_wise": [
        {
            "RegionName": "Asia",
            "Revenue": 2637.08
        },
        {
            "RegionName": "South America",
            "Revenue": 188
        },
        {
            "RegionName": "Europe",
            "Revenue": 1314
        },
        {
            "RegionName": "North America",
            "Revenue": 309.49
        }
    ],
    "status": "S",
    "errMsg": ""
}
```
### ❌ Failure Response Format (Validation or User Error)

```json
{
    "monthly_trend": null,
    "yearly_trend": null,
    "quarterly_trend": null,
    "total_revenue": {
        "TotalRevenue": 0
    },
    "product_wise": null,
    "category_wise": null,
    "region_wise": null,
    "status": "E",
    "errMsg": "  Invalid input: no valid indicator found"
}
```

### Program In-Out API  Logging

```json
{
    "RequestId": "a8d7996b76b2409ead8e767663915963",
    "RespBody": "{\"monthly_trend\":null,\"yearly_trend\":null,\"quarterly_trend\":null,\"total_revenue\":{\"TotalRevenue\":0},\"product_wise\":null,\"category_wise\":null,\"region_wise\":null,\"status\":\"E\",\"errMsg\":\"  Invalid input: no valid indicator found\"}\n",
    "ResponseStatus": 200,
    "ReqDateTime": "0001-01-01T00:00:00Z",
    "RealIP": "",
    "ForwardedIP": "[::1]:38316",
    "Method": "POST",
    "Path": "/revenuereport?",
    "Host": "localhost:19998",
    "RemoteAddr": "[::1]:38316",
    "Header": " Cache-Control-no-cache Accept-Encoding-gzip, deflate, br Connection-keep-alive Content-Length-98 Content-Type-text/plain User-Agent-PostmanRuntime/7.44.0 Accept-*/* Postman-Token-fb795341-2e3d-48c5-b6f3-02eec7faefd4",
    "Endpoint": "/revenuereport",
    "RespDateTime": "2025-05-25T10:18:40.766423598+05:30",
    "ReqBody": "{\n    \"start_date\": \"2024-01-01\",\n    \"end_date\": \"2024-12-31\",\n    \"indicator\": \"MONTHLYTRENDs\"\n}",
    "RequestUnixTime": 1748148520,
    "ResponseUnixTime": 1748148520,
    "PDebug": {
        "Sid": "a8d7996b76b2409ead8e767663915963",
        "Reference": ""
    }
}
```

###  🔁 Refresh Mechanism Implemented

```json
{
		"request_id" : "a5e19891-e4c0-4c31-9fbe-25bc4fcf05f9",
		"status" : "S",
		"total_records_affected" : 6,
		"error_message" : "",
		"created_time" : "2025-05-25T04:13:20.964Z",
		"updated_time" : "2025-05-25T04:13:21.055Z",
		"created_by" : "ondemandapi",
		"updated_by" : "ondemandapi",
		"refresh_type" : "Autobot",
		"duration_seconds" : 0
	}
```

# Automate the Cloning the Repo and Running the program through .bash file 

## FileName : startprogram.sh 
### Relative_FilePath : startprogram.sh

## Steps to Run the Above bash file

- Download the `startprogram.sh` file to any directory and follow the steps below.

```bash
chmod +x ./startprogram.sh

./startprogram.sh

```


- Validate TOML configuration to avoid parse errors.
- Ensure database schema is always synchronized with app (especially log table columns).
- API supports multiple indicator queries for flexible reporting.
- Logging infrastructure designed to help identify root causes quickly.


### 📝 Notes

- All endpoints return JSON.
- please change the toml file according to you requirement
