package wallet

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
)

type TransactionService struct {
	app.Service

	Query *TransactionQueryTask
}

func (S *TransactionService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *TransactionService) HandleTransactionQueryTask(a IWalletApp, task *TransactionQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.WalletId == 0 {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found wallet id"
		return nil
	}

	sql := bytes.NewBuffer(nil)
	args := []interface{}{}

	sql.WriteString(" WHERE walletid=?")

	args = append(args, task.WalletId)

	if task.OrderBy == "asc" {
		sql.WriteString(" ORDER BY id ASC")
	} else {
		sql.WriteString(" ORDER BY id DESC")
	}

	var pageIndex = task.PageIndex
	var pageSize = task.PageSize

	if pageIndex < 1 {
		pageIndex = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if task.Counter {
		var counter = TransactionQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize
		counter.RowCount, err = kk.DBQueryCount(db, a.GetTransactionTable(), a.GetPrefix(), sql.String(), args...)
		if err != nil {
			task.Result.Errno = ERROR_WALLET
			task.Result.Errmsg = err.Error()
			return nil
		}
		if counter.RowCount%pageSize == 0 {
			counter.PageCount = counter.RowCount / pageSize
		} else {
			counter.PageCount = counter.RowCount/pageSize + 1
		}
		task.Result.Counter = &counter
	}

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	var transactions = []Transaction{}
	var v = Transaction{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetTransactionTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_WALLET
			task.Result.Errmsg = err.Error()
			return nil
		}

		transactions = append(transactions, v)
	}

	task.Result.Transactions = transactions

	return nil
}
