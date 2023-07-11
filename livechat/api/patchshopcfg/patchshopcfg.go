// Package patchshopcfg defines the request and response data models for the API endpoint PATCH /shops/{shop_id}/config.
package patchshopcfg

// Request is the request data model for the API endpoint PATCH /shops/{shop_id}/config.
type Request struct {
	TemplatePayload string `json:"templatePayload"`
}

// Response is the response data model for the API endpoint PATCH /shops/{shop_id}/config.
type Response struct {
	TemplateID string `json:"templateID"`
}
