package handlers

import (
	"encoding/json"
	"net/http"
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/helpers"
)

func ProvideCustomerAnalysis(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.StartFunc()

	if r.Method == http.MethodPost {
		var lBody ConditonBody

		var lCustomerStatRec CustomerDetails

		lErr := json.NewDecoder(r.Body).Decode(&lBody)
		if lErr != nil {
			lDebug.Log(helpers.Elog, lErr.Error())
			lCustomerStatRec.ErrMsg = lErr.Error()
			lCustomerStatRec.Status = "E"
		} else {
			lErr = CustomerAnalysis(lDebug, &lCustomerStatRec, lBody.StartDate, lBody.EndDate)
			if lErr != nil {
				lDebug.Log(helpers.Elog, lErr.Error())
				lCustomerStatRec.ErrMsg = lErr.Error()
				lCustomerStatRec.Status = "E"
			} else {
				lCustomerStatRec.Status = "S"
				lCustomerStatRec.ErrMsg = ""
			}

		}
		lErr = json.NewEncoder(w).Encode(lCustomerStatRec)
		if lErr != nil {
			lDebug.Log(helpers.Elog, lErr.Error())
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func CustomerAnalysis(pDebug *helpers.HelperStruct, pCustomerStatRec *CustomerDetails, startDate, endDate string) error {
	pDebug.StartFunc()

	lResult := dbconnection.GRMPostgres.Table("orders o").Joins("join products p on o.product_id =p.product_id").Where(`o.date_of_sale between ? and ?`, startDate, endDate).Select("count(distinct o.customer_id) total_no_of_customers", "count(order_id) total_orders_cnt", "sum((o.quantity_sold*p.unit_price) -o.discount +o.shipping_cost)/count(o.order_id) average_order_value").Scan(&pCustomerStatRec)
	if lResult.Error != nil {
		pDebug.Log(helpers.Elog, lResult.Error.Error())
		return helpers.ErrReturn(lResult.Error)
	}
	pDebug.ExitFunc()
	return nil
}
