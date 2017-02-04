package wallet

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ExecuteTaskResult struct {
	app.Result
	Order *Order `json:"order,omitempty"`
}

type ExecuteTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result ExecuteTaskResult
}

func (task *ExecuteTask) GetResult() interface{} {
	return &task.Result
}

func (task *ExecuteTask) GetInhertType() string {
	return "wallet"
}

func (task *ExecuteTask) GetClientName() string {
	return "Order.Execute"
}
