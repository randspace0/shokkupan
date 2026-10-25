package dto

type NotificationsRequest struct {
	App string `query:"app" validate:"omitempty,appName"`
}

type GetNotificationsResponse struct {
	Installed  bool   `json:"installed"`
	Configured bool   `json:"configured"`
	WebhookURL string `json:"webhook_url"` // masked
}

type SetNotificationsRequest struct {
	App        string `query:"app" validate:"omitempty,appName"`
	WebhookURL string `json:"webhook_url" validate:"required"`
}
