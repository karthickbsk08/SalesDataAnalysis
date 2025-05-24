package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/helpers"
)

func ProvideCustomerAnalysis(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.StartFunc()

	if r.Method == http.MethodGet {

		var lCustomerStatRec CustomerDetails
		lCustomerStatRec.ErrMsg = ""
		lCustomerStatRec.Status = "S"

		lErr := CustomerAnalysis(lDebug, lCustomerStatRec)
		if lErr != nil {
			lDebug.Log(helpers.Elog, lErr.Error())
			lCustomerStatRec.ErrMsg = lErr.Error()
			lCustomerStatRec.Status = "E"
		}

		lErr = json.NewEncoder(w).Encode(lCustomerStatRec)
		if lErr != nil {
			lDebug.Log(helpers.Elog, lErr.Error())
			lCustomerStatRec.ErrMsg = lErr.Error()
			lCustomerStatRec.Status = "E"
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func CustomerAnalysis(pDebug *helpers.HelperStruct, pCustomerStatRec CustomerDetails) error {
	pDebug.StartFunc()

	lResult := dbconnection.GRMPostgres.Table("orders o").Joins("join products p on o.product_id =p.product_id").Select("count(distinct o.customer_id) total_no_of_customers", "count(order_id) total_orders_cnt", "sum((o.quantity_sold*p.unit_price) -o.discount +o.shipping_cost)/count(o.order_id) average_order_value").Scan(&pCustomerStatRec)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error.Error())
		return helpers.ErrReturn(lResult.Error)
	}
	fmt.Println(pCustomerStatRec)
	pDebug.ExitFunc()
	return nil
}
