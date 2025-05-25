package refreshmechanism

import (
	"encoding/csv"
	"fmt"
	"os"
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/helpers"
	"strconv"
	"time"

	"gorm.io/gorm/clause"
)

/*
func FetchPrevStagingDetails(pDebug *helpers.HelperStruct) (map[string]StagingSale, error) {
	pDebug.StartFunc()

	var lPrevStagingRecords []StagingSale
	DupCheckMap := make(map[string]StagingSale, 0)
	lResult := dbconnection.GRMPostgres.Table("staging_sales").Scan(&lPrevStagingRecords)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error)
		return DupCheckMap, helpers.ErrReturn(lResult.Error)
	}

	if len(lPrevStagingRecords) > 0 && len(DupCheckMap) == 0 {
		for i := 0; i < len(lPrevStagingRecords); i++ {
			key := fmt.Sprintf("%d_%s", lPrevStagingRecords[i].OrderID, lPrevStagingRecords[i].ProductID)
			DupCheckMap[key] = lPrevStagingRecords[i]
		}
	}

	pDebug.ExitFunc()
	return DupCheckMap, nil
}

func (a StagingSale) Equals(b StagingSale) bool {
	return a.OrderID == b.OrderID &&
		a.ProductID == b.ProductID &&
		a.CustomerID == b.CustomerID &&
		a.ProductName == b.ProductName &&
		a.Category == b.Category &&
		a.Region == b.Region &&
		a.DateOfSale.Equal(b.DateOfSale) && // for time.Time
		a.QuantitySold == b.QuantitySold &&
		a.UnitPrice == b.UnitPrice &&
		a.Discount == b.Discount &&
		a.ShippingCost == b.ShippingCost &&
		a.PaymentMethod == b.PaymentMethod &&
		a.CustomerName == b.CustomerName &&
		a.CustomerEmail == b.CustomerEmail &&
		a.CustomerAddress == b.CustomerAddress
} */

func LoadCSVInStagingTable(pDebug *helpers.HelperStruct) (int64, error) {
	pDebug.StartFunc()

	var lTotalAffectedRecordCnt int64
	// var lPrevStagingMap map[string]StagingSale
	var lNewRecArr []StagingSale
	// var lChangesAvail bool
	// var minorChangesAvail bool

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

	// lPrevStagingMap, lErr = FetchPrevStagingDetails(pDebug)
	// if lErr != nil {
	// 	pDebug.Log(helpers.Elog, fmt.Sprintf("Failed to read CSV: %s", lErr))
	// 	return lTotalAffectedRecordCnt, helpers.ErrReturn(lErr)
	// }

	// // check existing prev record length
	// if len(lPrevStagingMap) == 0 && len(records) > 0 {
	// 	lChangesAvail = true
	// } else {
	// 	lChangesAvail = true
	// }

	// if lChangesAvail {
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
		dateOfSale, _ := time.Parse("02-01-2006", row[6])

		// map to check duplicate
		// key := fmt.Sprintf("%d_%s", orderID, row[1])

		entrynew := StagingSale{
			OrderID:         orderID,
			ProductID:       row[1],
			CustomerID:      row[2],
			ProductName:     row[3],
			Category:        row[4],
			Region:          row[5],
			DateOfSale:      dateOfSale,
			QuantitySold:    quantity,
			UnitPrice:       unitPrice,
			Discount:        discount,
			ShippingCost:    shipping,
			PaymentMethod:   row[11],
			CustomerName:    row[12],
			CustomerEmail:   row[13],
			CustomerAddress: row[14],
		}

		// if val, ok := lPrevStagingMap[key]; ok {
		// 	if !entrynew.Equals(val) {
		// 		minorChangesAvail = true
		// 	}
		// } else {
		// 	minorChangesAvail = true
		// }
		lNewRecArr = append(lNewRecArr, entrynew)
	}

	// if minorChangesAvail {
	lresult := dbconnection.GRMPostgres.Table("staging_sales").Clauses(
		clause.OnConflict{
			UpdateAll: true,
			Columns:   []clause.Column{{Name: "order_id"}}}).CreateInBatches(&lNewRecArr, 1000)
	if lresult.Error != nil {
		pDebug.Log(helpers.Elog, fmt.Sprintf("Insert error: %s", lresult.Error))
		return lTotalAffectedRecordCnt, helpers.ErrReturn(lresult.Error)
	}
	lTotalAffectedRecordCnt = lresult.RowsAffected
	// }

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
	lResult := dbconnection.GRMPostgres.Table("staging_sales s").Select("DISTINCT s.category AS name").Scan(&lCategoriesData)
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

/* entry := map[string]interface{}{
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
} */
