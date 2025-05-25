package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"salesdataanalysis/constants"
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/helpers"
	"strings"
)

type ConditonBody struct {
	StartDate                 string `json:"start_date"`
	EndDate                   string `json:"end_date"`
	ConsolidatedDataIndicator string `json:"indicator"`
}

func RevenueReportAPI(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.StartFunc()

	if r.Method == http.MethodPost {
		var lBody ConditonBody
		var lRespRec RevenueReport
		lRespRec.ErrMsg = ""
		lRespRec.Status = "S"

		lErr := json.NewDecoder(r.Body).Decode(&lBody)
		if lErr != nil {
			lDebug.Log(helpers.Elog, lErr.Error())
			lRespRec.ErrMsg = helpers.ErrPrint(lErr)
			lRespRec.Status = "E"
		} else {
			lErr = GenerateRevenueReport(lDebug, lBody, &lRespRec)
			if lErr != nil {
				lDebug.Log(helpers.Elog, lErr.Error())
				lRespRec.ErrMsg = helpers.ErrPrint(lErr)
				lRespRec.Status = "E"
			}
		}
		lErr = json.NewEncoder(w).Encode(lRespRec)
		if lErr != nil {
			lDebug.Log(helpers.Elog, lErr.Error())
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	lDebug.ExitFunc()

}

func GenerateRevenueReport(pDebug *helpers.HelperStruct, pConditionRec ConditonBody, pReport *RevenueReport) error {
	pDebug.StartFunc()

	pConditionRec.ConsolidatedDataIndicator = strings.TrimSpace(pConditionRec.ConsolidatedDataIndicator)
	if pConditionRec.ConsolidatedDataIndicator == "" {
		return helpers.ErrReturn(fmt.Errorf(" No Indicator Provided"))
	}
	valid := []string{"MONTHLYTREND", "YEARLYTREND", "QUARTERLYTREND", "TOTALREV", "PRODUCTWISE", "CATEGORYWISE", "REGIONWISE"}

	indicators := strings.Split(pConditionRec.ConsolidatedDataIndicator, ",")
	foundValid := false
	for _, ind := range indicators {
		ind = strings.TrimSpace(ind)
		for _, v := range valid {
			if ind == v {
				foundValid = true
				break
			}
		}
		if foundValid {
			break
		}
	}

	if !foundValid {
		return helpers.ErrReturn(fmt.Errorf(" Invalid input: no valid indicator found"))

	}
	for i := 0; i < len(indicators); i++ {
		if indicators[i] == constants.MONTHLY_TREND_IND {
			// Monthly
			if err := dbconnection.GRMPostgres.Raw(`
		SELECT TO_CHAR(o.date_of_sale, 'YYYY-MM') AS month_year,
		        ROUND(SUM((o.quantity_sold * p.unit_price) * (1 - o.discount) + o.shipping_cost), 2) AS month_trend
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
		where o.date_of_sale between ? and ?
		GROUP BY TO_CHAR(o.date_of_sale, 'YYYY-MM')
		ORDER BY TO_CHAR(o.date_of_sale, 'YYYY-MM') DESC
	`, pConditionRec.StartDate, pConditionRec.EndDate).Scan(&pReport.MonthlyTrends).Error; err != nil {
				return helpers.ErrReturn(err)
			}
		}

		if indicators[i] == constants.YEARLY_TREND_IND {
			// Yearly
			if err := dbconnection.GRMPostgres.Raw(`
		SELECT TO_CHAR(o.date_of_sale, 'YYYY') AS yearwise,
		       ROUND(SUM((o.quantity_sold * p.unit_price) * (1 - o.discount) + o.shipping_cost), 2) AS year_trend
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
		where o.date_of_sale between ? and ?
		GROUP BY TO_CHAR(o.date_of_sale, 'YYYY')
		ORDER BY TO_CHAR(o.date_of_sale, 'YYYY') DESC
	`, pConditionRec.StartDate, pConditionRec.EndDate).Scan(&pReport.YearlyTrends).Error; err != nil {
				return helpers.ErrReturn(err)
			}
		}

		if indicators[i] == constants.QUAT_TREND_IND {
			// Quarterly
			if err := dbconnection.GRMPostgres.Raw(`
		SELECT EXTRACT(YEAR FROM o.date_of_sale) || '-' || 'Q' || LPAD(EXTRACT(QUARTER FROM o.date_of_sale)::TEXT, 1, '0') AS year_quater,
		       ROUND(SUM((o.quantity_sold * p.unit_price) * (1 - o.discount) + o.shipping_cost), 2) AS year_quarter
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
		where o.date_of_sale between ? and ?
		GROUP BY EXTRACT(YEAR FROM o.date_of_sale), EXTRACT(QUARTER FROM o.date_of_sale)
		ORDER BY EXTRACT(YEAR FROM o.date_of_sale) DESC, EXTRACT(QUARTER FROM o.date_of_sale) DESC
	`, pConditionRec.StartDate, pConditionRec.EndDate).Scan(&pReport.QuarterlyTrends).Error; err != nil {
				return helpers.ErrReturn(err)
			}
		}

		if indicators[i] == constants.TOTAL_REVENUE_IND {
			// Total revenue
			if err := dbconnection.GRMPostgres.Raw(`
		SELECT ROUND(SUM((o.quantity_sold * p.unit_price) * (1 - o.discount) + o.shipping_cost), 2) AS totalRevenue
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
		WHERE o.date_of_sale BETWEEN ? AND ?
	`, pConditionRec.StartDate, pConditionRec.EndDate).Scan(&pReport.TotalRevenue).Error; err != nil {
				return helpers.ErrReturn(err)
			}
		}

		if indicators[i] == constants.PRODUCT_WISE_IND {
			// Product-wise
			if err := dbconnection.GRMPostgres.Raw(`
		SELECT p.name, ROUND(SUM((o.quantity_sold * p.unit_price) * (1 - o.discount) + o.shipping_cost), 2) AS prod_by_totalrevenue
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
		WHERE o.date_of_sale BETWEEN ? AND ?
		GROUP BY p.product_id
	`, pConditionRec.StartDate, pConditionRec.EndDate).Scan(&pReport.ProductWise).Error; err != nil {
				return helpers.ErrReturn(err)
			}

		}

		if indicators[i] == constants.CATEGORY_WISE_IND {
			// Category-wise
			if err := dbconnection.GRMPostgres.Raw(`
		SELECT c.name, ROUND(SUM((o.quantity_sold * p.unit_price) * (1 - o.discount) + o.shipping_cost), 2) AS cat_by_totalrevenue
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
		JOIN categories c ON p.category_id = c.category_id
		where o.date_of_sale between ? and ?
		GROUP BY c.category_id
		ORDER BY c.category_id
	`, pConditionRec.StartDate, pConditionRec.EndDate).Scan(&pReport.CategoryWise).Error; err != nil {
				return helpers.ErrReturn(err)
			}
		}

		if indicators[i] == constants.REGION_IND {
			// Region-wise
			if err := dbconnection.GRMPostgres.Raw(`
		SELECT r.name, ROUND(SUM((o.quantity_sold * p.unit_price) * (1 - o.discount) + o.shipping_cost), 2) AS region_by_totalrevenue
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
		JOIN regions r ON o.region_id = r.region_id
		where o.date_of_sale between ? and ?
		GROUP BY r.region_id
		ORDER BY r.region_id
	`, pConditionRec.StartDate, pConditionRec.EndDate).Scan(&pReport.RegionWise).Error; err != nil {
				return helpers.ErrReturn(err)
			}
		}

	}

	pDebug.ExitFunc()
	return nil
}
