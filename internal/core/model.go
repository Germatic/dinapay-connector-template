package core

import (
	"time"

	contract "github.com/Germatic/dinapay-contracts/go/connectorcontract/failures"
)

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
	DestinationMode      string            `json:"destinationMode,omitempty"`
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

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type PayoutDestination struct {
	Country     string         `json:"country"`
	Currency    string         `json:"currency"`
	Beneficiary map[string]any `json:"beneficiary"`
	Rail        map[string]any `json:"rail"`
}

type CreatePayoutCommand struct {
	OperationID          string            `json:"operationId"`
	PayoutID             string            `json:"payoutId"`
	AccountID            string            `json:"accountId"`
	MerchantID           string            `json:"merchantId"`
	Provider             string            `json:"provider"`
	ProviderConnectionID string            `json:"providerConnectionId"`
	Binding              *Binding          `json:"binding,omitempty"`
	Source               Money             `json:"source"`
	Destination          PayoutDestination `json:"destination"`
	Remitter             map[string]any    `json:"remitter,omitempty"`
	Description          string            `json:"description,omitempty"`
	Metadata             map[string]any    `json:"metadata,omitempty"`
}

type CancelPayoutCommand struct {
	OperationID          string `json:"operationId"`
	PayoutID             string `json:"payoutId"`
	ProviderConnectionID string `json:"providerConnectionId"`
	Reason               string `json:"reason,omitempty"`
}

type ProviderPayment struct {
	TransactionID        string                    `json:"transactionId"`
	Provider             string                    `json:"provider"`
	ProviderConnectionID string                    `json:"providerConnectionId"`
	ProviderPaymentID    string                    `json:"providerPaymentId"`
	ProviderReference    string                    `json:"providerReference,omitempty"`
	Status               string                    `json:"status"`
	RawStatus            string                    `json:"rawStatus,omitempty"`
	ObservedAt           time.Time                 `json:"observedAt"`
	ExpiresAt            time.Time                 `json:"expiresAt,omitempty"`
	Completion           map[string]any            `json:"completion,omitempty"`
	ProviderData         map[string]any            `json:"providerData,omitempty"`
	Failure              *contract.Failure         `json:"failure,omitempty"`
	ProviderFailure      *contract.ProviderFailure `json:"providerFailure,omitempty"`
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

type ProviderPayout struct {
	PayoutID             string                    `json:"payoutId"`
	Provider             string                    `json:"provider"`
	ProviderConnectionID string                    `json:"providerConnectionId"`
	ProviderPayoutID     string                    `json:"providerPayoutId"`
	ProviderReference    string                    `json:"providerReference,omitempty"`
	Status               string                    `json:"status"`
	RawStatus            string                    `json:"rawStatus,omitempty"`
	Source               Money                     `json:"source"`
	DestinationAmount    string                    `json:"destinationAmount,omitempty"`
	DestinationCurrency  string                    `json:"destinationCurrency,omitempty"`
	ObservedAt           time.Time                 `json:"observedAt"`
	ProviderData         map[string]any            `json:"providerData,omitempty"`
	Failure              *contract.Failure         `json:"failure,omitempty"`
	ProviderFailure      *contract.ProviderFailure `json:"providerFailure,omitempty"`
}

type Capabilities struct {
	Provider        string       `json:"provider"`
	ContractVersion string       `json:"contractVersion"`
	Capabilities    []Capability `json:"capabilities"`
}

type Capability struct {
	Operation             string              `json:"operation"`
	Countries             []string            `json:"countries"`
	Currencies            []string            `json:"currencies,omitempty"`
	SourceCurrencies      []string            `json:"sourceCurrencies,omitempty"`
	DestinationCurrencies []string            `json:"destinationCurrencies,omitempty"`
	PaymentMethods        []string            `json:"paymentMethods,omitempty"`
	Rails                 []string            `json:"rails"`
	DestinationModes      []string            `json:"destinationModes,omitempty"`
	Features              []string            `json:"features,omitempty"`
	BindingRequirement    *BindingRequirement `json:"bindingRequirement,omitempty"`
}

type BindingRequirement struct {
	EntityType         string `json:"entityType"`
	ExternalEntityType string `json:"externalEntityType"`
}

type ProviderEvent struct {
	EventID              string            `json:"eventId"`
	EventType            string            `json:"eventType"`
	EventVersion         string            `json:"eventVersion"`
	Source               string            `json:"source"`
	OccurredAt           time.Time         `json:"occurredAt"`
	ObservedAt           time.Time         `json:"observedAt"`
	TransactionID        string            `json:"transactionId,omitempty"`
	PayoutID             string            `json:"payoutId,omitempty"`
	Provider             string            `json:"provider"`
	ProviderConnectionID string            `json:"providerConnectionId"`
	ProviderPaymentID    string            `json:"providerPaymentId,omitempty"`
	ProviderPayoutID     string            `json:"providerPayoutId,omitempty"`
	Data                 ProviderEventData `json:"data"`
}

type ProviderEventData struct {
	Status            string                    `json:"status"`
	RawStatus         string                    `json:"rawStatus"`
	Amount            string                    `json:"amount,omitempty"`
	Currency          string                    `json:"currency,omitempty"`
	ProviderReference string                    `json:"providerReference,omitempty"`
	ProviderData      map[string]any            `json:"providerData,omitempty"`
	Failure           *contract.Failure         `json:"failure,omitempty"`
	ProviderFailure   *contract.ProviderFailure `json:"providerFailure,omitempty"`
}

type RawWebhook struct {
	ConnectionID string
	Headers      map[string][]string
	Body         []byte
}
