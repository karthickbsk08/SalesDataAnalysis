package refreshmechanism

import (
	"salesdataanalysis/helpers"
	"time"
)

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
