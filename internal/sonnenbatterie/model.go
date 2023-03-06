package sonnenbatterie

import (
	"github.com/dabump/sonnenbatterie/internal/config"
)

type Client struct {
	httpClient HttpClient
	config     *config.Config
}
type Status struct {
	ApparentOutput            int     `json:"Apparent_output,omitempty"`
	BackupBuffer              string  `json:"BackupBuffer,omitempty"`
	BatteryCharging           bool    `json:"BatteryCharging,omitempty"`
	BatteryDischarging        bool    `json:"BatteryDischarging,omitempty"`
	ConsumptionAvg            int     `json:"Consumption_Avg,omitempty"`
	ConsumptionW              int     `json:"Consumption_W,omitempty"`
	Fac                       float64 `json:"Fac,omitempty"`
	FlowConsumptionBattery    bool    `json:"FlowConsumptionBattery,omitempty"`
	FlowConsumptionGrid       bool    `json:"FlowConsumptionGrid,omitempty"`
	FlowConsumptionProduction bool    `json:"FlowConsumptionProduction,omitempty"`
	FlowGridBattery           bool    `json:"FlowGridBattery,omitempty"`
	FlowProductionBattery     bool    `json:"FlowProductionBattery,omitempty"`
	FlowProductionGrid        bool    `json:"FlowProductionGrid,omitempty"`
	GridFeedInW               int     `json:"GridFeedIn_W,omitempty"`
	IsSystemInstalled         int     `json:"IsSystemInstalled,omitempty"`
	OperatingMode             string  `json:"OperatingMode,omitempty"`
	PacTotalW                 int     `json:"Pac_total_W,omitempty"`
	ProductionW               int     `json:"Production_W,omitempty"`
	Rsoc                      int     `json:"RSOC,omitempty"`
	RemainingCapacityWh       int     `json:"RemainingCapacity_Wh,omitempty"`
	Sac1                      int     `json:"Sac1,omitempty"`
	Sac2                      int     `json:"Sac2,omitempty"`
	Sac3                      int     `json:"Sac3,omitempty"`
	SystemStatus              string  `json:"SystemStatus,omitempty"`
	Timestamp                 string  `json:"Timestamp,omitempty"`
	Usoc                      int     `json:"USOC,omitempty"`
	Uac                       int     `json:"Uac,omitempty"`
	Ubat                      int     `json:"Ubat,omitempty"`
	DischargeNotAllowed       bool    `json:"dischargeNotAllowed,omitempty"`
	GeneratorAutostart        bool    `json:"generator_autostart,omitempty"`
}
