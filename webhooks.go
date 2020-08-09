package upngo

import "time"

type WebhookEventType string

const (
	WebhookEventTypeTransactionCreated WebhookEventType = "TRANSACTION_CREATED"
	WebhookEventTypeTransactionSettled WebhookEventType = "TRANSACTION_SETTLED"
	WebhookEventTypeTransactionDeleted WebhookEventType = "TRANSACTION_DELETED"
	WebhookEventTypePing               WebhookEventType = "PING"
)

type WebhookResource struct {
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
	Data  []WebhookResource `json:"data"`
	Links LinksObject       `json:"links"`
}

type WebhookResponse struct {
	Data  WebhookResource `json:"data"`
	Links LinksObject     `json:"links"`
}

type WebhookInputResourceAttributes struct {
	// Max length 300 chars
	URL string `json:"url"`
	// Max length 64 chars
	Description string `json:"description"`
}

type WebhookInputResource struct {
	Attributes WebhookInputResourceAttributes `json:"attributes"`
}

type RegisterWebhookRequest struct {
	Data WebhookInputResource `json:"data"`
}

type WebhookEventResourceAttributes struct {
	EventType WebhookEventType `json:"eventType"`
	CreatedAt time.Time        `json:"createdAt"`
}

type WebhookEventResourceRelationships struct {
	Webhook struct {
		Data struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
		Links struct {
			Related string `json:"related"`
		} `json:"links"`
	} `json:"webhook"`
	Transaction struct {
		Data struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
		Links RelatedLinksObject `json:"links"`
	} `json:"transaction"`
}

type WebhookEventResource struct {
	Type          string                            `json:"type"`
	ID            string                            `json:"id"`
	Attributes    WebhookEventResourceAttributes    `json:"attributes"`
	Relationships WebhookEventResourceRelationships `json:"relationships"`
}

type WebhookPingResponse struct {
	Data WebhookEventResource `json:"data"`
}
