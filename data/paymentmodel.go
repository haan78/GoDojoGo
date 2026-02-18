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
