package deff

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const VERSION = "1.0.0+1"

type SettingsType struct {
	PORT               int
	MYSQL_DSN          string
	JWT_SECRET         string
	BREVO_API_KEY      string
	EMAIL_ACTIVATE_URL string
}

var Settings SettingsType

func LoadSettings(path string) {

	_, err := os.Stat(path)
	if err == nil {
		err := godotenv.Load(path)
		if err != nil {
			panic(err)
		}
	}

	Settings.JWT_SECRET = os.Getenv("JWT_SECRET")
	Settings.MYSQL_DSN = os.Getenv("MYSQL_DSN")
	Settings.BREVO_API_KEY = os.Getenv("BREVO_API_KEY")
	if os.Getenv("PORT") != "" {
		num, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil {
			Settings.PORT = num
		} else {
			panic(err)
		}
	} else {
		Settings.PORT = 1323
	}
}
