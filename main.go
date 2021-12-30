package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/durianpay/dpay-common/db"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type IDs struct {
	merchantID string
	paymentID  string
}

type items struct {
	ID               int    `db:"id"`
	DsSettlementFees int    `db:"ds_settlement_fee"`
	DsPaymentRefID   string `db:"ds_payment_ref_id"`
	Difference       int    `db:"difference"`
}

const (
	getItemsToBeUpdated = `Select sd.id, sd.ds_settlement_fee, sd.ds_payment_ref_id, (sd.ds_settlement_fee-385000) AS "difference" FROM settlement_details AS sd 
		JOIN payment AS p ON p.id = sd.payment_id WHERE p.provider_id = 'XENDIT' AND p.payment_details_type = 'va_details'
		AND sd.merchant_id = $1 AND sd.payment_id = $2 ;`

	update = `Update settlement_details 
		Set ds_settlement_fee = 385000 
		where id = $1
		RETURNING id;`

	insertIntoLogsTable = `Insert Into Logs
	(payment_id, merchant_id, payment_ds_ref_id, ds_settlement_fee_charged, updated_ds_settlement_fee, "difference")
	Values ($1, $2, $3, $4, 385000, $5)
	Returning id`
)

func InitDB() (err error) {
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Println("error", err)
		return
	}

	dbName := os.Getenv("DB_NAME")

	err = db.Init(&db.Config{
		Driver: "postgres",
		URL:    fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbName),
	})
	if err != nil {
		fmt.Println("error", err)
		return
	}

	return
}

func GetItemsToBeUpdated(ctx context.Context, db *sqlx.DB, merchantID string, paymentID string) (dbItems []items, err error) {
	err = db.SelectContext(ctx, &dbItems, getItemsToBeUpdated, merchantID, paymentID)
	if err != nil {
		fmt.Printf("error getting items from db %s %s\n", "error", err.Error())
		return
	}
	return
}

func UpdateDsFees(ctx context.Context, db *sqlx.DB, dbItem []items, merchantID string, paymentID string) (err error) {
	if len(dbItem) == 0 {
		fmt.Println("no rows to updated for merchant id:", merchantID, "and payment id:", paymentID)
		return
	}

	for _, item := range dbItem {
		var ID1 int
		var ID2 int

		difference := int(item.DsSettlementFees/100) - 3850

		err = db.GetContext(ctx, &ID1, insertIntoLogsTable, paymentID, merchantID, item.DsPaymentRefID, item.DsSettlementFees, difference)
		if err != nil {
			fmt.Println("error inserting into logs table", item.ID, "error", err.Error())
			continue
		}

		err = db.GetContext(ctx, &ID2, update, item.ID)
		if err != nil {
			fmt.Println("error updating ds_settlement_fees with id", item.ID, "error", err.Error())
			continue
		}
	}
	fmt.Println("updated for merchant id:", merchantID, "and payement id:", paymentID)
	return
}

func main() {

	err := InitDB()
	if err != nil {
		fmt.Println("error initializing db")
		return
	}

	appDB := db.Get()
	ctx := context.TODO()

	data, err := GetData()
	if err != nil {
		fmt.Println("error reading/getting data")
		return
	}

	for i := 0; i < len(data); i++ {
		merchantID := data[i].merchantID
		paymentID := data[i].paymentID

		dbItems, err := GetItemsToBeUpdated(ctx, appDB, merchantID, paymentID)
		if err != nil {
			fmt.Println("error in getting rows from the tables")
			return
		}

		err = UpdateDsFees(ctx, appDB, dbItems, merchantID, paymentID)
		if err != nil {
			fmt.Println("error in updating fees", err.Error())
			return
		}
	}
}
