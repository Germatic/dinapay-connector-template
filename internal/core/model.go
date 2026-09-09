package core

import "time"

type Binding struct {
	BindingID          string `json:"bindingId"`
	EntityType         string `json:"entityType"`
	EntityID           string `json:"entityId"`
	ExternalEntityType string `json:"externalEntityType"`
	ExternalEntityID   string `json:"externalEntityId"`
}

type CreatePaymentCommand struct {
	OperationID          string            `json:"operationId"`
	TransactionID        string            `json:"transactionId"`
	Provider             string            `json:"provider"`
	ProviderConnectionID string            `json:"providerConnectionId"`
	Binding              *Binding          `json:"binding,omitempty"`
	Amount               string            `json:"amount"`
	Currency             string            `json:"currency"`
	PaymentMethod        string            `json:"paymentMethod"`
	Rail                 string            `json:"rail,omitempty"`
	Description          string            `json:"description,omitempty"`
	ExpiresAt            time.Time         `json:"expiresAt,omitempty"`
	ReturnURLs           map[string]string `json:"returnUrls,omitempty"`
	Customer             map[string]any    `json:"customer"`
	Metadata             map[string]any    `json:"metadata,omitempty"`
}

type CancelPaymentCommand struct {
	OperationID          string `json:"operationId"`
	TransactionID        string `json:"transactionId"`
	ProviderConnectionID string `json:"providerConnectionId"`
	Reason               string `json:"reason,omitempty"`
}

type CreateRefundCommand struct {
	OperationID          string         `json:"operationId"`
	RefundID             string         `json:"refundId"`
	TransactionID        string         `json:"transactionId"`
	ProviderConnectionID string         `json:"providerConnectionId"`
	Amount               string         `json:"amount"`
	Currency             string         `json:"currency"`
	Reason               string         `json:"reason,omitempty"`
	Metadata             map[string]any `json:"metadata,omitempty"`
}

type ProviderPayment struct {
	TransactionID        string         `json:"transactionId"`
	Provider             string         `json:"provider"`
	ProviderConnectionID string         `json:"providerConnectionId"`
	ProviderPaymentID    string         `json:"providerPaymentId"`
	ProviderReference    string         `json:"providerReference,omitempty"`
	Status               string         `json:"status"`
	RawStatus            string         `json:"rawStatus,omitempty"`
	ObservedAt           time.Time      `json:"observedAt"`
	ExpiresAt            time.Time      `json:"expiresAt,omitempty"`
	Completion           map[string]any `json:"completion,omitempty"`
	ProviderData         map[string]any `json:"providerData,omitempty"`
}

type ProviderRefund struct {
	RefundID             string         `json:"refundId"`
	TransactionID        string         `json:"transactionId"`
	Provider             string         `json:"provider"`
	ProviderConnectionID string         `json:"providerConnectionId"`
	ProviderRefundID     string         `json:"providerRefundId"`
	Status               string         `json:"status"`
	RawStatus            string         `json:"rawStatus,omitempty"`
	Amount               string         `json:"amount"`
	Currency             string         `json:"currency"`
	ObservedAt           time.Time      `json:"observedAt"`
	ProviderData         map[string]any `json:"providerData,omitempty"`
}

type Capabilities struct {
	Provider        string       `json:"provider"`
	ContractVersion string       `json:"contractVersion"`
	Capabilities    []Capability `json:"capabilities"`
}

type Capability struct {
	Operation          string              `json:"operation"`
	Countries          []string            `json:"countries"`
	Currencies         []string            `json:"currencies"`
	PaymentMethods     []string            `json:"paymentMethods"`
	Rails              []string            `json:"rails"`
	DestinationModes   []string            `json:"destinationModes,omitempty"`
	Features           []string            `json:"features,omitempty"`
	BindingRequirement *BindingRequirement `json:"bindingRequirement,omitempty"`
}

type BindingRequirement struct {
	EntityType         string `json:"entityType"`
	ExternalEntityType string `json:"externalEntityType"`
}

type ProviderEvent struct {
	EventID              string         `json:"eventId"`
	EventType            string         `json:"eventType"`
	EventVersion         string         `json:"eventVersion"`
	Source               string         `json:"source"`
	OccurredAt           time.Time      `json:"occurredAt"`
	ObservedAt           time.Time      `json:"observedAt"`
	TransactionID        string         `json:"transactionId"`
	Provider             string         `json:"provider"`
	ProviderConnectionID string         `json:"providerConnectionId"`
	ProviderPaymentID    string         `json:"providerPaymentId"`
	Data                 map[string]any `json:"data"`
}

type RawWebhook struct {
	ConnectionID string
	Headers      map[string][]string
	Body         []byte
}
