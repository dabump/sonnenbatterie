package sonnenbatterie

import (
	"github.com/dabump/sonnenbatterie/internal/common"
	"github.com/dabump/sonnenbatterie/internal/config"
)

type Client struct {
	httpClient common.HttpClient
	config     *config.Config
}

type Status struct {
	ApparentOutput            int     `json:"Apparent_output,omitempty"`
	BackupBuffer              string  `json:"BackupBuffer,omitempty"`
	BatteryCharging           bool    `json:"BatteryCharging,omitempty"`
	BatteryDischarging        bool    `json:"BatteryDischarging,omitempty"`
	ConsumptionAvg            float64 `json:"Consumption_Avg,omitempty"`
	ConsumptionW              float64 `json:"Consumption_W,omitempty"`
	Fac                       float64 `json:"Fac,omitempty"`
	FlowConsumptionBattery    bool    `json:"FlowConsumptionBattery,omitempty"`
	FlowConsumptionGrid       bool    `json:"FlowConsumptionGrid,omitempty"`
	FlowConsumptionProduction bool    `json:"FlowConsumptionProduction,omitempty"`
	FlowGridBattery           bool    `json:"FlowGridBattery,omitempty"`
	FlowProductionBattery     bool    `json:"FlowProductionBattery,omitempty"`
	FlowProductionGrid        bool    `json:"FlowProductionGrid,omitempty"`
	GridFeedInW               float64 `json:"GridFeedIn_W,omitempty"`
	IsSystemInstalled         int     `json:"IsSystemInstalled,omitempty"`
	OperatingMode             string  `json:"OperatingMode,omitempty"`
	PacTotalW                 float64 `json:"Pac_total_W,omitempty"`
	ProductionW               float64 `json:"Production_W,omitempty"`
	Rsoc                      float64 `json:"RSOC,omitempty"`
	RemainingCapacityWh       float64 `json:"RemainingCapacity_Wh,omitempty"`
	Sac1                      float64 `json:"Sac1,omitempty"`
	Sac2                      float64 `json:"Sac2,omitempty"`
	Sac3                      float64 `json:"Sac3,omitempty"`
	SystemStatus              string  `json:"SystemStatus,omitempty"`
	Timestamp                 string  `json:"Timestamp,omitempty"`
	Usoc                      float64 `json:"USOC,omitempty"`
	Uac                       float64 `json:"Uac,omitempty"`
	Ubat                      float64 `json:"Ubat,omitempty"`
	DischargeNotAllowed       bool    `json:"dischargeNotAllowed,omitempty"`
	GeneratorAutostart        bool    `json:"generator_autostart,omitempty"`
}
