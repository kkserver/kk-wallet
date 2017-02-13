package wallet

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RevokeTaskResult struct {
	app.Result
	Order *Order `json:"order,omitempty"`
}

type RevokeTask struct {
	app.Task
	Name        string      `json:"name"`
	Freeze      bool        `json:"freeze"`
	WalletId    int64       `json:"walletId"`
	Value       int64       `json:"value"`
	Title       string      `json:"title"`
	NotifyUrl   string      `json:"notifyUrl"`
	Options     interface{} `json:"options"`
	AssociateId int64       `json:"associateId"`
	Result      RevokeTaskResult
}

func (task *RevokeTask) GetResult() interface{} {
	return &task.Result
}

func (task *RevokeTask) GetInhertType() string {
	return "wallet"
}

func (task *RevokeTask) GetClientName() string {
	return "Order.Revoke"
}
