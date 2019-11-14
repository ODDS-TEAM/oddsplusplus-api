package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	// Specification holds environment variable name.
	Specification struct {
		DBHost       string
		DBName       string
		DBStatusCol  string
		DBUsersCol   string
		DBItemCol    string
		DBTypeCol    string
		DBReserveCol string
		DBSummaryCol    string
		ImgPath      string
		APIPort      string
	}
)

// Spec retrieves the value of the environment variable named by the key.
func Spec() *Specification {
	godotenv.Load()

	s := Specification{
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBUsersCol: os.Getenv("DB_USERS_COL"),
		DBStatusCol:	os.Getenv("DB_STATUS_COL"),
		DBItemCol: os.Getenv("DB_ITEM_COL"),
		DBTypeCol:	os.Getenv("DB_TPYE_COL"),
		DBReserveCol: os.Getenv("DB_RESERVE_COL"),
		DBSummaryCol: os.Getenv("DB_SUMMARY_COL"),
		APIPort: os.Getenv("API_PORT"),
	}
	return &s
}
