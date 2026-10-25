package dto

type MonitoringApp struct {
	Name         string `json:"name"`
	IsDeployed   bool   `json:"is_deployed"`
	IsRunning    bool   `json:"is_running"`
	NumProcesses int    `json:"num_processes"`
	Restore      bool   `json:"restore"`
}

type MonitoringEvent struct {
	AppName string `json:"app_name"`
	Running bool   `json:"running"`
	At      string `json:"at"`
}

type GetMonitoringResponse struct {
	StatsAvailable bool              `json:"stats_available"`
	Apps           []MonitoringApp   `json:"apps"`
	Events         []MonitoringEvent `json:"events"`
}

type GetMonitoringStatsRequest struct {
	App   string `query:"app" validate:"appName"`
	Hours int    `query:"hours" validate:"omitempty,min=1,max=168"`
}

type MonitoringStatsPoint struct {
	At       string  `json:"at"`
	CPU      float64 `json:"cpu"`
	MemBytes int64   `json:"mem_bytes"`
	MemLimit int64   `json:"mem_limit"`
}

type GetMonitoringStatsResponse struct {
	Points []MonitoringStatsPoint `json:"points"`
}
