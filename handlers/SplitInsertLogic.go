package handlers

import (
	"encoding/csv"
	"fmt"
	"os"
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/helpers"
	"strconv"
	"time"
)

func LoadCSVInStagingTable(pDebug *helpers.HelperStruct) (int64, error) {
	pDebug.StartFunc()

	var lTotalAffectedRecordCnt int64

	filePath := "/home/it-lap9/Videos/SalesDataAnalysis/dbfilecollection/salesdata.csv"
	file, lErr := os.Open(filePath)
	if lErr != nil {
		pDebug.Log(helpers.Elog, fmt.Sprintf("Failed to open file: %s", lErr))
		return lTotalAffectedRecordCnt, helpers.ErrReturn(lErr)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, lErr := reader.ReadAll()
	if lErr != nil {
		pDebug.Log(helpers.Elog, fmt.Sprintf("Failed to read CSV: %s", lErr))
		return lTotalAffectedRecordCnt, helpers.ErrReturn(lErr)
	}

	// Skip header
	for i, row := range records {
		if i == 0 {
			continue
		}

		orderID, _ := strconv.Atoi(row[0])
		quantity, _ := strconv.Atoi(row[7])
		unitPrice, _ := strconv.ParseFloat(row[8], 64)
		discount, _ := strconv.ParseFloat(row[9], 64)
		shipping, _ := strconv.ParseFloat(row[10], 64)
		dateOfSale, _ := time.Parse("2006-01-02", row[6])

		entry := map[string]interface{}{
			"order_id":         orderID,
			"product_id":       row[1],
			"customer_id":      row[2],
			"product_name":     row[3],
			"category":         row[4],
			"region":           row[5],
			"date_of_sale":     dateOfSale,
			"quantity_sold":    quantity,
			"unit_price":       unitPrice,
			"discount":         discount,
			"shipping_cost":    shipping,
			"payment_method":   row[11],
			"customer_name":    row[12],
			"customer_email":   row[13],
			"customer_address": row[14],
		}

		lresult := dbconnection.GRMPostgres.Table("staging_sales").Create(&entry)
		if lresult.Error != nil {
			pDebug.Log(helpers.Elog, fmt.Sprintf("Insert error: %s", lresult.Error))
			return lTotalAffectedRecordCnt, helpers.ErrReturn(lresult.Error)
		}
		lTotalAffectedRecordCnt = lresult.RowsAffected
	}

	pDebug.ExitFunc()
	return lTotalAffectedRecordCnt, nil
}

func InsertIntoAllTablesUniquely(pDebug *helpers.HelperStruct) error {
	pDebug.StartFunc()
	var lCategoriesData []Category
	var lRegionsData []Region
	var lProductData []Product
	var lCustomerData []Customer
	var lOrderData []Order

	//categories data
	lResult := dbconnection.GRMPostgres.Table("staging_sales s").Distinct("category name").Scan(&lCategoriesData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	lResult = dbconnection.GRMPostgres.Table("categories").Create(&lCategoriesData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	//Regions data
	lResult = dbconnection.GRMPostgres.Table("staging_sales").Distinct("region name").Scan(&lRegionsData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	lResult = dbconnection.GRMPostgres.Table("regions").Create(&lRegionsData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	//products data
	lResult = dbconnection.GRMPostgres.Table("staging_sales s").Joins("LEFT JOIN categories c ON c.name=s.category").Distinct(" s.product_id product_id ,s.product_name name,c.category_id category_id,s.unit_price unit_price").Scan(&lProductData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	lResult = dbconnection.GRMPostgres.Table("products").Create(&lProductData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	//customers data customer_id,"name",email,address
	lResult = dbconnection.GRMPostgres.Table("staging_sales s").Distinct("s.customer_id customer_id,s.customer_name name,s.customer_email email,s.customer_address address").Scan(&lCustomerData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	lResult = dbconnection.GRMPostgres.Table("customers").Create(&lCustomerData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	//customers data customer_id,"name",email,address
	lResult = dbconnection.GRMPostgres.Table("staging_sales ss").Joins(`
	left join products p on p.product_id =ss.product_id 
	left join customers c on c.customer_id =ss.customer_id 
	left join regions r on r."name"=ss.region  
	left join categories c2 on c2.name=ss.category`).Select(`ss.order_id order_id,p.product_id product_id,c.customer_id customer_id,r.region_id region_id,
ss.date_of_sale date_of_sale,ss.quantity_sold quantity_sold,ss.discount discount,ss.shipping_cost shipping_cost,ss.payment_method payment_method`).Scan(&lOrderData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	lResult = dbconnection.GRMPostgres.Table("orders").Create(&lOrderData)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return helpers.ErrReturn(lResult.Error)
	}

	pDebug.ExitFunc()
	return nil

}
