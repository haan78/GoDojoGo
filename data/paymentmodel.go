package data

import lib "GoDojoGo/lib"

type PatmentModelType struct {
	Name   string  `db:"model_name" json:"model_name"`
	Charge float64 `db:"charge" json:"charge"`
}

func PeymentGetAll() ([]PatmentModelType, error) {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		list, err := lib.GenericQuery[PatmentModelType](db, `SELECT model_name, charge FROM paymentmodel`)
		if err == nil {
			return list, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func PaymentSetAll(plist []PatmentModelType) error {
	db, err := lib.DbConnect()
	if err == nil {
		defer db.Close()
		tx, err := db.Beginx()
		if err == nil {
			for _, p := range plist {
				_, err := tx.Exec("UPDATE paymentmodel SET charge = ? WHERE model_name = ?", p.Charge, p.Name)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
			tx.Commit()
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
