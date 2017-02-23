package wallet

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TransactionQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
	RowCount  int `json:"rowCount"`
}

type TransactionQueryTaskResult struct {
	app.Result
	Counter      *TransactionQueryCounter `json:"counter,omitempty"`
	Transactions []Transaction            `json:"transactions,omitempty"`
}

type TransactionQueryTask struct {
	app.Task
	WalletId  int64  `json:"walletId"`
	OrderId   int64  `json:"orderId"`
	Status    string `json:"status"`
	OrderBy   string `json:"orderBy"` // desc, asc
	PageIndex int    `json:"p"`
	PageSize  int    `json:"size"`
	Counter   bool   `json:"counter"`
	Result    TransactionQueryTaskResult
}

func (task *TransactionQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *TransactionQueryTask) GetInhertType() string {
	return "wallet"
}

func (task *TransactionQueryTask) GetClientName() string {
	return "Transaction.Query"
}
