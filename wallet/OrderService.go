package wallet

import (
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"time"
)

type OrderService struct {
	app.Service

	Recharge *RechargeTask
	Revoke   *RevokeTask
	Transfer *TransferTask
	Get      *OrderTask
	Execute  *ExecuteTask
	Cancel   *CancelTask
}

func (S *OrderService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *OrderService) HandleRechargeTask(a IWalletApp, task *RechargeTask) error {

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

	if task.Value <= 0 {
		task.Result.Errno = ERROR_WALLET_VALUE
		task.Result.Errmsg = "Incorrect amount"
		return nil
	}

	v := Order{}

	v.Title = task.Title
	v.Action = OrderActionRecharge
	v.Ctime = time.Now().Unix()

	if task.Freeze {
		v.Action = v.Action | OrderActionFreeze
	}

	options := map[interface{}]interface{}{}

	items := []interface{}{}

	items = append(items, map[interface{}]interface{}{"walletId": task.WalletId, "value": task.Value})

	options["items"] = items
	options["payType"] = task.PayType
	options["payTradeNo"] = task.PayTradeNo

	b, err := json.Encode(options)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	v.Options = string(b)

	_, err = kk.DBInsert(db, a.GetOrderTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Order = &v

	return nil
}

func (S *OrderService) HandleRevokeTask(a IWalletApp, task *RevokeTask) error {

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

	if task.Value <= 0 {
		task.Result.Errno = ERROR_WALLET_VALUE
		task.Result.Errmsg = "Incorrect amount"
		return nil
	}

	v := Order{}

	v.Title = task.Title
	v.Action = OrderActionRevoke
	v.Ctime = time.Now().Unix()

	if task.Freeze {
		v.Action = v.Action | OrderActionFreeze
	}

	options := map[interface{}]interface{}{}

	items := []interface{}{}

	items = append(items, map[interface{}]interface{}{"walletId": task.WalletId, "value": -task.Value})

	options["items"] = items
	options["payType"] = task.PayType
	options["payTradeNo"] = task.PayTradeNo

	b, err := json.Encode(options)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	v.Options = string(b)

	_, err = kk.DBInsert(db, a.GetOrderTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Order = &v

	return nil
}

func (S *OrderService) HandleTransferTask(a IWalletApp, task *TransferTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.FwalletId == 0 {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found payment wallet id"
		return nil
	}

	if task.TwalletId == 0 {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found receipt wallet id"
		return nil
	}

	if task.Value <= 0 {
		task.Result.Errno = ERROR_WALLET_VALUE
		task.Result.Errmsg = "Incorrect amount"
		return nil
	}

	v := Order{}

	v.Title = task.Title
	v.Action = OrderActionTransfer
	v.Ctime = time.Now().Unix()

	if task.Freeze {
		v.Action = v.Action | OrderActionFreeze
	}

	options := map[interface{}]interface{}{}

	items := []interface{}{}

	items = append(items,
		map[interface{}]interface{}{"walletId": task.FwalletId, "value": -task.Value},
		map[interface{}]interface{}{"walletId": task.TwalletId, "value": task.Value})

	options["items"] = items

	b, err := json.Encode(options)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	v.Options = string(b)

	_, err = kk.DBInsert(db, a.GetOrderTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Order = &v

	return nil
}

func (S *OrderService) HandleOrderTask(a IWalletApp, task *OrderTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Id == 0 {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found order id"
		return nil
	}

	v := Order{}

	rows, err := kk.DBQuery(db, a.GetOrderTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_WALLET
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Order = &v

	} else {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND
		task.Result.Errmsg = "Not Found order"
		return nil
	}

	return nil
}

func (S *OrderService) HandleExecuteTask(a IWalletApp, task *ExecuteTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Id == 0 {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found order id"
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	err = func() error {

		rows, err := kk.DBQuery(tx, a.GetOrderTable(), a.GetPrefix(), " WHERE id=? FOR UPDATE", task.Id)

		if err != nil {
			return err
		}

		if rows.Next() {

			v := Order{}

			scanner := kk.NewDBScaner(&v)

			err = scanner.Scan(rows)

			rows.Close()

			if err != nil {
				return err
			}

			task.Result.Order = &v

			if v.Status == OrderStatusFreeze {

				rs, err := kk.DBQuery(tx, a.GetTransactionTable(), a.GetPrefix(), " WHERE orderid=? ORDER BY id ASC", v.Id)

				if err != nil {
					return err
				}

				tran := Transaction{}
				trans := []Transaction{}
				scanner = kk.NewDBScaner(&tran)

				for rs.Next() {

					err = scanner.Scan(rs)

					if err != nil {
						rs.Close()
						return err
					}

					trans = append(trans, tran)
				}

				rs.Close()

				wallet := Wallet{}
				scanner = kk.NewDBScaner(&wallet)

				for _, tran = range trans {

					rs, err = kk.DBQuery(tx, a.GetWalletTable(), a.GetPrefix(), " WHERE id=? FOR UPDATE", tran.WalletId)

					if err != nil {
						return err
					}

					if rs.Next() {

						err = scanner.Scan(rs)

						rs.Close()

						if err != nil {
							return err
						}

						if tran.Value > 0 {
							wallet.Value = wallet.Value + tran.Value
							wallet.Freeze = wallet.Freeze - tran.Value
						} else if tran.Value < 0 {
							wallet.Freeze = wallet.Freeze + tran.Value
						}

						_, err = kk.DBUpdateWithKeys(tx, a.GetWalletTable(), a.GetPrefix(), &wallet, map[string]bool{"value": true, "freeze": true})

						if err != nil {
							return err
						}

					} else {
						rs.Close()
						return app.NewError(ERROR_WALLET_NOT_FOUND, "Not Found wallet")
					}
				}

				v.Status = OrderStatusOK

				_, err = kk.DBUpdateWithKeys(tx, a.GetOrderTable(), a.GetPrefix(), &v, map[string]bool{"status": true})

				if err != nil {
					return err
				}

			} else if v.Status == OrderStatusNone {

				var options interface{} = nil

				err = json.Decode([]byte(v.Options), &options)

				if err != nil {
					return err
				}

				dynamic.Each(dynamic.Get(options, "items"), func(key interface{}, item interface{}) bool {

					walletId := dynamic.IntValue(dynamic.Get(item, "walletId"), 0)
					value := dynamic.IntValue(dynamic.Get(item, "value"), 0)

					if walletId == 0 {
						err = app.NewError(ERROR_WALLET_NOT_FOUND_ID, "Not Found wallet id")
						return false
					}

					if value == 0 {
						err = app.NewError(ERROR_WALLET_VALUE, "Incorrect amount")
						return false
					}

					wallet := Wallet{}

					scanner = kk.NewDBScaner(&wallet)

					rs, err := kk.DBQuery(tx, a.GetWalletTable(), a.GetPrefix(), " WHERE id=? FOR UPDATE", walletId)

					if err != nil {
						return false
					}

					if rs.Next() {

						err = scanner.Scan(rs)

						rs.Close()

						if err != nil {
							return false
						}

					} else {
						rs.Close()
						err = app.NewError(ERROR_WALLET_NOT_FOUND, "Not Found wallet")
						return false
					}

					if (v.Action & OrderActionFreeze) == OrderActionFreeze {

						if value > 0 {

							wallet.Freeze = wallet.Freeze + value

						} else if value < 0 {

							if wallet.Value+value < 0 {
								err = app.NewError(ERROR_WALLET_VALUE, "Incorrect amount")
								return false
							}

							wallet.Value = wallet.Value + value
							wallet.Freeze = wallet.Freeze - value

						}

					} else {

						if value > 0 {

							wallet.Value = wallet.Value + value

						} else if value < 0 {

							if wallet.Value+value < 0 {
								err = app.NewError(ERROR_WALLET_VALUE, "Incorrect amount")
								return false
							}

							wallet.Value = wallet.Value + value

						}
					}

					_, err = kk.DBUpdateWithKeys(tx, a.GetWalletTable(), a.GetPrefix(), &wallet, map[string]bool{"value": true, "freeze": true})

					if err != nil {
						return false
					}

					tran := Transaction{}

					tran.WalletId = walletId
					tran.Value = value
					tran.OrderId = v.Id
					tran.Ctime = time.Now().Unix()

					_, err = kk.DBInsert(tx, a.GetTransactionTable(), a.GetPrefix(), &tran)

					if err != nil {
						return false
					}

					return true
				})

				if err != nil {
					return err
				}

				if (v.Action & OrderActionFreeze) == OrderActionFreeze {
					v.Status = OrderStatusFreeze
					v.Action = v.Action & OrderActionMask
				} else {
					v.Status = OrderStatusOK
				}

				_, err = kk.DBUpdateWithKeys(tx, a.GetOrderTable(), a.GetPrefix(), &v, map[string]bool{"status": true, "action": true})

				if err != nil {
					return err
				}

			} else {
				return app.NewError(ERROR_WALLET_STATUS, "The current state can not be executed")
			}

		} else {
			rows.Close()
			return app.NewError(ERROR_WALLET_NOT_FOUND, "Not Found order")
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

func (S *OrderService) HandleCancelTask(a IWalletApp, task *CancelTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Id == 0 {
		task.Result.Errno = ERROR_WALLET_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found order id"
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_WALLET
		task.Result.Errmsg = err.Error()
		return nil
	}

	err = func() error {

		rows, err := kk.DBQuery(tx, a.GetOrderTable(), a.GetPrefix(), " WHERE id=? FOR UPDATE", task.Id)

		if err != nil {
			return err
		}

		if rows.Next() {

			v := Order{}

			scanner := kk.NewDBScaner(&v)

			err = scanner.Scan(rows)

			rows.Close()

			if err != nil {
				return err
			}

			task.Result.Order = &v

			if v.Status == OrderStatusFreeze {

				rs, err := kk.DBQuery(tx, a.GetTransactionTable(), a.GetPrefix(), " WHERE orderid=? ORDER BY id ASC", v.Id)

				if err != nil {
					return err
				}

				tran := Transaction{}
				trans := []Transaction{}
				scanner = kk.NewDBScaner(&tran)

				for rs.Next() {

					err = scanner.Scan(rs)

					if err != nil {
						rs.Close()
						return err
					}

					trans = append(trans, tran)
				}

				rs.Close()

				wallet := Wallet{}
				scanner = kk.NewDBScaner(&wallet)

				for _, tran = range trans {

					rs, err = kk.DBQuery(tx, a.GetWalletTable(), a.GetPrefix(), " WHERE id=? FOR UPDATE", tran.WalletId)

					if err != nil {
						return err
					}

					if rs.Next() {

						err = scanner.Scan(rs)

						rs.Close()

						if err != nil {
							return err
						}

						if tran.Value > 0 {
							wallet.Freeze = wallet.Freeze - tran.Value
						} else if tran.Value < 0 {
							wallet.Value = wallet.Value - tran.Value
							wallet.Freeze = wallet.Freeze + tran.Value
						}

						_, err = kk.DBUpdateWithKeys(tx, a.GetWalletTable(), a.GetPrefix(), &wallet, map[string]bool{"value": true, "freeze": true})

						if err != nil {
							return err
						}

					} else {
						rs.Close()
						return app.NewError(ERROR_WALLET_NOT_FOUND, "Not Found wallet")
					}
				}

				v.Status = OrderStatusCancel

				_, err = kk.DBUpdateWithKeys(tx, a.GetOrderTable(), a.GetPrefix(), &v, map[string]bool{"status": true})

				if err != nil {
					return err
				}

			} else {
				return app.NewError(ERROR_WALLET_STATUS, "The current state can not be executed")
			}

		} else {
			rows.Close()
			return app.NewError(ERROR_WALLET_NOT_FOUND, "Not Found order")
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
