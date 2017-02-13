package wallet

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RechargeTaskResult struct {
	app.Result
	Order *Order `json:"order,omitempty"`
}

type RechargeTask struct {
	app.Task
	Name        string      `json:"name"`
	Freeze      bool        `json:"freeze"`
	WalletId    int64       `json:"walletId"`
	Value       int64       `json:"value"`
	Title       string      `json:"title"`
	NotifyUrl   string      `json:"notifyUrl"`
	Options     interface{} `json:"options"`
	AssociateId int64       `json:"associateId"`
	Result      RechargeTaskResult
}

func (task *RechargeTask) GetResult() interface{} {
	return &task.Result
}

func (task *RechargeTask) GetInhertType() string {
	return "wallet"
}

func (task *RechargeTask) GetClientName() string {
	return "Order.Recharge"
}
