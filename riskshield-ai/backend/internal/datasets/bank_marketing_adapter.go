package datasets

import (
	"os"
	"path/filepath"
	"time"
)

type BankMarketingAdapter struct {
	BaseDir string
}

func (a *BankMarketingAdapter) GetMetadata() DatasetMetadata {
	path1 := filepath.Join(a.BaseDir, "features.csv")
	path2 := filepath.Join(a.BaseDir, "targets.csv")
	
	_, err1 := os.Stat(path1)
	_, err2 := os.Stat(path2)
	
	if err1 != nil || err2 != nil {
		return DatasetMetadata{
			ID: "bank_marketing",
			Name: "Bank Marketing",
			Source: "UCI Bank Marketing Dataset",
			Status: "NOT AVAILABLE",
		}
	}

	return DatasetMetadata{
		ID: "bank_marketing",
		Name: "Bank Marketing",
		Source: "UCI Bank Marketing Dataset",
		Status: "AVAILABLE",
		RowCount: 45211,
		ColumnCount: 17, 
		QualityScore: 95.0,
		HasFraudLabels: false, // Used for fairness
		LastScannedAt: time.Now().Format(time.RFC3339),
	}
}

func (a *BankMarketingAdapter) ReadStream(ch chan<- CanonicalTransaction, maxRows int) error {
	// Not meant for transaction streaming
	close(ch)
	return nil
}
