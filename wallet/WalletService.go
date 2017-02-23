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

	sql.WriteString(" ORDER BY id ASC LIMIT 1")

	v := Wallet{}

	tx, err := db.Begin()

	err = func() error {

		rows, err := kk.DBQuery(tx, a.GetWalletTable(), a.GetPrefix(), sql.String(), args...)

		if err != nil {
			return err
		}

		if rows.Next() {

			scanner := kk.NewDBScaner(&v)

			err = scanner.Scan(rows)

			rows.Close()

			if err != nil {
				return err
			}

			task.Result.Wallet = &v

		} else {
			rows.Close()

			if task.Name != "" && task.Autocreate {

				v.Uid = task.Uid
				v.Name = task.Name
				v.Ctime = time.Now().Unix()

				_, err = kk.DBInsert(tx, a.GetWalletTable(), a.GetPrefix(), &v)

				if err != nil {
					return err
				}

				task.Result.Wallet = &v

			} else {
				return app.NewError(ERROR_WALLET_NOT_FOUND, "Not Found wallet")
			}
		}

		return nil
	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		tx.Rollback()
		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		} else {
			task.Result.Errno = ERROR_WALLET
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	return nil
}
