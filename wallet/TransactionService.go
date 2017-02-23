package wallet

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"strings"
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

	if task.Status == "" {
		sql.WriteString(fmt.Sprintf("FROM %s%s as t", a.GetPrefix(), a.GetTransactionTable().Name))
	} else {
		sql.WriteString(fmt.Sprintf("FROM %s%s as t LEFT JOIN %s%s as o ON t.orderid=o.id", a.GetPrefix(), a.GetTransactionTable().Name, a.GetPrefix(), a.GetOrderTable().Name))
	}

	sql.WriteString(" WHERE t.walletid=?")

	args = append(args, task.WalletId)

	if task.OrderId != 0 {
		sql.WriteString(" AND t.orderid=?")
		args = append(args, task.OrderId)
	}

	if task.Status != "" {

		sql.WriteString(" AND o.status IN (")

		for i, s := range strings.Split(task.Status, ",") {
			if i != 0 {
				sql.WriteString(",")
			}
			sql.WriteString("?")
			args = append(args, s)
		}

		sql.WriteString(")")

	}
	if task.OrderBy == "asc" {
		sql.WriteString(" ORDER BY t.id ASC")
	} else {
		sql.WriteString(" ORDER BY t.id DESC")
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

		rs, err := db.Query("SELECT COUNT(*) "+sql.String(), args...)

		if err != nil {
			task.Result.Errno = ERROR_WALLET
			task.Result.Errmsg = err.Error()
			return nil
		}

		defer rs.Close()

		if rs.Next() {

			err = rs.Scan(&counter.RowCount)

			if err != nil {
				task.Result.Errno = ERROR_WALLET
				task.Result.Errmsg = err.Error()
				return nil
			}
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

	rows, err := db.Query("SELECT t.* "+sql.String(), args...)

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
