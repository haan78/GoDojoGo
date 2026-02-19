package main

import (
	"GoDojoGo/Cron/data"
	"GoDojoGo/Cron/deff"
	"GoDojoGo/Cron/lib"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

type SendMailType struct {
	EmailPoolId int64  `db:"emailpool_id" json:"emailpool_id"`
	Kind        string `db:"kind" json:"kind"`
	Email       string `db:"email" json:"email"`
	Success     bool   `db:"success" json:"success"`
	Message     string `db:"mseesage" json:"message"`
}

func SendMails(db *sqlx.DB, list []data.UnsendEmailType) []SendMailType {

	anError := func(d *data.UnsendEmailType, s string) SendMailType {
		return SendMailType{
			EmailPoolId: d.EmailPoolId,
			Success:     false,
			Message:     s,
			Kind:        d.Kind,
			Email:       d.Email,
		}
	}

	var result []SendMailType

	for _, email := range list {
		var tid int = 0
		switch email.Kind {
		case "ACTIVATE":
			tid = 1
		case "INFORM":
			tid = 5
		}
		if tid > 0 {
			tx, err := db.Beginx()
			if err == nil {
				err := data.SetEmailAsSent(tx, email.EmailPoolId)
				if err == nil {
					err := lib.SendinblueTemplateEmail(deff.Settings.BREVO_API_KEY, email.Email, 1, email.GetParams())
					if err == nil {
						tx.Commit()
						result = append(result, SendMailType{
							EmailPoolId: email.EmailPoolId,
							Success:     true,
							Message:     "Mail successfuly sent",
							Kind:        email.Kind,
							Email:       email.Email,
						})
					} else {
						tx.Rollback()
						result = append(result, SendMailType{
							EmailPoolId: email.EmailPoolId,
							Success:     false,
							Message:     err.Error(),
							Kind:        email.Kind,
							Email:       email.Email,
						})
					}
				} else {
					result = append(result, SendMailType{
						EmailPoolId: email.EmailPoolId,
						Success:     false,
						Message:     err.Error(),
						Kind:        email.Kind,
						Email:       email.Email,
					})
				}
			} else {
				result = append(result, SendMailType{
					EmailPoolId: email.EmailPoolId,
					Success:     false,
					Message:     err.Error(),
					Kind:        email.Kind,
					Email:       email.Email,
				})
			}
		} else {
			result = append(result, anError(&email, "kind is unknown"))
		}
		time.Sleep(5 * time.Second)

	}
	return result
}

func main() {

	deff.LoadSettings(".env")

	MYSQL_DSN := os.Getenv("MYSQL_DSN")

	db, err := lib.DbConnect(MYSQL_DSN)
	if err == nil {
		defer db.Close()
		list, err := data.GetUnsentEmails(db)
		if err == nil {
			rList := SendMails(db, list)
			var ct, ce, cs int = 0, 0, 0
			for _, r := range rList {
				if r.Success {
					fmt.Printf("Success email = %s, kind = %s \n", r.Email, r.Kind)
					cs += 1
				} else {
					fmt.Printf("Error email = %s, kind = %s, message = %s \n", r.Email, r.Kind, r.Message)
					ce += 1
				}
				ct += 1
			}
			fmt.Printf("Total %d Success %d Error %d\n", ct, cs, ce)
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
