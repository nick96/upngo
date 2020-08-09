package upngo

import "time"

type WebhooksResource struct {
	Type          string             `json:"type"`
	ID            string             `json:"id"`
	Attributes    WebhooksAttributes `json:"attributes"`
	Relationships LogsRelationship   `json:"relationships"`
	Links         SelfLinkObject     `json:"links"`
}

type WebhooksAttributes struct {
	URL         string    `json:"url"`
	Description string    `json:"description"`
	SecretKey   string    `json:"secretKey"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LogsRelatedObject struct {
	Related string `json:"related"`
}

type LogsLinkObject struct {
	Links LogsRelatedObject `json:"links"`
}

type LogsRelationship struct {
	Logs LogsLinkObject `json:"logs"`
}

type WebhooksResponse struct {
	Data  []WebhooksResource `json:"data"`
	Links LinksObject        `json:"links"`
}
