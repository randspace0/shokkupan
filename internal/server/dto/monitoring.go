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
	Apps   []MonitoringApp   `json:"apps"`
	Events []MonitoringEvent `json:"events"`
}
