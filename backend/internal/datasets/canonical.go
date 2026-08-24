package datasets

import (
	"time"
)

type CanonicalTransaction struct {
	TransactionID            string    `json:"transaction_id"`
	CustomerID               string    `json:"customer_id"`
	MerchantID               string    `json:"merchant_id"`
	Timestamp                time.Time `json:"timestamp"`
	Amount                   float64   `json:"amount"`
	Currency                 string    `json:"currency"`
	TransactionType          string    `json:"transaction_type"`
	MerchantCategory         string    `json:"merchant_category"`
	Location                 string    `json:"location"`
	Country                  string    `json:"country"`
	DeviceID                 string    `json:"device_id"`
	DeviceType               string    `json:"device_type"`
	IPAddress                string    `json:"ip_address"`
	PaymentMethod            string    `json:"payment_method"`
	AccountAgeDays           int       `json:"account_age_days"`
	TransactionVelocity      int       `json:"transaction_velocity"`
	PreviousTransactionCount int       `json:"previous_transaction_count"`
	PreviousFraudCount       int       `json:"previous_fraud_count"`
	FraudLabel               *bool     `json:"fraud_label"`
	SourceDataset            string    `json:"source_dataset"`
	RawRecordHash            string    `json:"raw_record_hash"`
	Synthetic                bool      `json:"synthetic"`
}

type DatasetMetadata struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Source           string  `json:"source"`
	Status           string  `json:"status"`
	RowCount         int64   `json:"row_count"`
	ColumnCount      int     `json:"column_count"`
	QualityScore     float64 `json:"quality_score"`
	HasFraudLabels   bool    `json:"has_fraud_labels"`
	LastScannedAt    string  `json:"last_scanned_at"`
}

type DatasetAdapter interface {
	GetMetadata() DatasetMetadata
	ReadStream(ch chan<- CanonicalTransaction, maxRows int) error
}
