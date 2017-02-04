package wallet

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CancelTaskResult struct {
	app.Result
	Order *Order `json:"order,omitempty"`
}

type CancelTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result CancelTaskResult
}

func (task *CancelTask) GetResult() interface{} {
	return &task.Result
}

func (task *CancelTask) GetInhertType() string {
	return "wallet"
}

func (task *CancelTask) GetClientName() string {
	return "Order.Cancel"
}
