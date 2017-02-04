package wallet

import (
	"bytes"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"time"
)

type WalletService struct {
	app.Service

	Get *WalletTask
}

func (S *WalletService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *WalletService) HandleWalletTask(a IWalletApp, task *WalletTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Id == 0 && task.Name == "" {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found wallet id"
		return nil
	}

	sql := bytes.NewBuffer(nil)
	args := []interface{}{}

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	}

	if task.Name != "" {
		sql.WriteString(" AND name=?")
		args = append(args, task.Name)
	}

	if task.Uid != 0 {
		sql.WriteString(" AND uid=?")
		args = append(args, task.Uid)
	}

	rows, err := kk.DBQuery(db, a.GetWalletTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	v := Wallet{}

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_WALLET
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Wallet = &v

	} else if task.Autocreate {

		v.Uid = task.Uid
		v.Name = task.Name
		v.Ctime = time.Now().Unix()

		_, err = kk.DBInsert(db, a.GetWalletTable(), a.GetPrefix(), &v)

		if err != nil {
			task.Result.Errno = ERROR_WALLET
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Wallet = &v

	} else {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND
		task.Result.Errmsg = "Not Found wallet"
		return nil
	}

	return nil
}
