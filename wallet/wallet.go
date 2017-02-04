package wallet

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/client"
	"github.com/kkserver/kk-lib/kk/app/remote"
	"github.com/kkserver/kk-lib/kk/json"
	Value "github.com/kkserver/kk-lib/kk/value"
	"math/rand"
	"reflect"
	"time"
)

const OrderStatusNone = 0
const OrderStatusFreeze = 100
const OrderStatusOK = 200
const OrderStatusCancel = 300

const OrderActionNone = 0
const OrderActionRecharge = 1
const OrderActionRevoke = 2
const OrderActionTransfer = 3
const OrderActionFreeze = 1 << 15

type Wallet struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Uid    int64  `json:"uid"`
	Value  int64  `json:"value"`
	Freeze int64  `json:"freeze"`
	Ctime  int64  `json:"ctime"`
}

type Order struct {
	Id      int64  `json:"id"`
	Uid     int64  `json:"uid"`
	Action  int    `json:"action"`
	Title   string `json:"title"`
	Options string `json:"options"`
	Status  int    `json:"status"`
	Ctime   int64  `json:"ctime"`
}

type Transaction struct {
	Id       int64 `json:"id"`
	WalletId int64 `json:"walletId"`
	OrderId  int64 `json:"orderId"`
	Value    int64 `json:"value"`
	Freeze   int64 `json:"freeze"`
	Ctime    int64 `json:"ctime"`
}

type IWalletApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetWalletTable() *kk.DBTable
	GetOrderTable() *kk.DBTable
	GetTransactionTable() *kk.DBTable
}

type WalletApp struct {
	app.App
	DB *app.DBConfig

	Remote *remote.Service
}

func (C *WalletApp) GetDB() (*sql.DB, error) {
	return C.DB.Get(C)
}

func (C *WalletApp) GetPrefix() string {
	return C.DB.Prefix
}

func (C *WalletApp) GetWalletTable() *kk.DBTable {
	return nil
}

func (C *WalletApp) GetOrderTable() *kk.DBTable {
	return nil
}

func (C *WalletApp) GetTransactionTable() *kk.DBTable {
	return nil
}
