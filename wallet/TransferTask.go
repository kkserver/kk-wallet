package wallet

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TransferTaskResult struct {
	app.Result
	Order *Order `json:"order,omitempty"`
}

type TransferTask struct {
	app.Task
	Name        string      `json:"name"`
	Freeze      bool        `json:"freeze"`
	FwalletId   int64       `json:"fwalletId"`
	TwalletId   int64       `json:"twalletId"`
	Value       int64       `json:"value"`
	Title       string      `json:"title"`
	NotifyUrl   string      `json:"notifyUrl"`
	Options     interface{} `json:"options"`
	AssociateId int64       `json:"associateId"`
	Result      TransferTaskResult
}

func (task *TransferTask) GetResult() interface{} {
	return &task.Result
}

func (task *TransferTask) GetInhertType() string {
	return "wallet"
}

func (task *TransferTask) GetClientName() string {
	return "Order.Transfer"
}
