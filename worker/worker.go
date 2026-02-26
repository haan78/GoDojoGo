package worker

import (
	"GoDojoGo/data"
	"GoDojoGo/deff"
	"GoDojoGo/lib"
	"GoDojoGo/libbrevo"
	"context"
	"fmt"
	"time"
)

func doWork(ctx context.Context) ([]string, error) {
	db, err := lib.DbConnectX(ctx)
	if err == nil {
		defer db.Close()
		list, err := data.GetEmailList(db, "PENDING", 5)
		if err == nil {
			var errarr []string
			for _, e := range list {
				msg := fmt.Sprintf("Guid %s, Email %s, Kind %s, Status ", e.Guid, e.Email, e.Kind)
				tx, err := db.Beginx()
				if err == nil {
					err := data.SetEmail(tx, e.Guid, "SENT")
					if err == nil {
						err := libbrevo.SendinblueTemplateEmail(deff.Settings.BREVO_API_KEY, e.Email, e.GetTempId(), e.GetParams())
						if err == nil {
							tx.Commit()
							msg += "OK!"
						} else {
							tx.Rollback()
							msg += "Error: " + err.Error()
						}
						time.Sleep(5 * time.Second)
					} else {
						msg += "Error: " + err.Error()
					}
				} else {
					msg += "Error: " + err.Error()
				}
				errarr = append(errarr, msg)
			}
			return errarr, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func Run(ctx context.Context, sec int) {
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			list, err := doWork(ctx)
			if err == nil {
				for _, r := range list {
					fmt.Println("email sending report")
					fmt.Printf("%s", r)
				}
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}
