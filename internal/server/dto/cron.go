package dto

type CronTask struct {
	ID       string `json:"id"`
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
}
type GetAppCronResponse struct {
	Tasks []CronTask `json:"tasks"`
}

type RunAppCronRequest struct {
	Name string `json:"name" validate:"appName"`
	ID   string `json:"id" validate:"required,alphanum"`
}
