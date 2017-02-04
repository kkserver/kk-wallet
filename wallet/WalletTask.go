package wallet

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WalletTaskResult struct {
	app.Result
	Wallet *Wallet `json:"wallet,omitempty"`
}

type WalletTask struct {
	app.Task
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Uid        int64  `json:"uid"`
	Autocreate bool   `json:"autocreate"`
	Result     WalletTaskResult
}

func (task *WalletTask) GetResult() interface{} {
	return &task.Result
}

func (task *WalletTask) GetInhertType() string {
	return "wallet"
}

func (task *WalletTask) GetClientName() string {
	return "Wallet.Get"
}
