package controller

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"leadgen/dict"
	"leadgen/utils/helper"
	utilsheets "leadgen/utils/utilSheets"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var db *sql.DB
var srv *sheets.Service

func init() {
	var err error

	// Loading environment variables
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Making connection with the db
	var ConStr string = fmt.Sprintf("host=%v port=%v dbname=%v user=%v connect_timeout=10 password=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"))
	db, err = sql.Open("postgres", ConStr)
	if err != nil {
		log.Fatal(err)
	}

	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//Initializing Sheets & Drive Service
	ctx := context.Background()

	credentialFile := dict.SheetsCred["test"]
	srv, err = sheets.NewService(ctx, option.WithCredentialsFile(credentialFile))
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateLead(c *gin.Context) {
	w := c.Writer
	r := c.Request

	w.Header().Set("Content-Type", "application/json")

	//HANDLING REQUEST BODY
	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	UID := data["UID"]
	Security_token := data["SECURITY"]
	VendorID := data["VENDOR_ID"].(string)
	Comment := data["COMMENT"]

	json.NewEncoder(w).Encode(data)
	// fmt.Println(data)
	// for key, value := range data {
	// 	fmt.Println("key ", key, reflect.TypeOf(key), "value", value, reflect.TypeOf(value))
	// }

	//QUERY & PARSING DATA FROM THE DB
	query := fmt.Sprintf(`SELECT * FROM attendes WHERE "uid"=%v AND "security"='%v';`, UID, Security_token)
	fmt.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("encountered error while querying")
		log.Fatal(err)
	}
	defer rows.Close()

	var atnde_security string
	var atnde_name string
	var atnde_email string
	var atnde_phone int
	var atnde_uid int
	var values [][]interface{}
	for rows.Next() {

		if err := rows.Scan(&atnde_security, &atnde_name, &atnde_email, &atnde_phone, &atnde_uid); err != nil {
			log.Fatal(err)
		}
		fmt.Println(atnde_security, atnde_name, atnde_email, atnde_phone, atnde_uid)
		values = [][]interface{}{
			{atnde_email, atnde_name, atnde_phone, atnde_security, atnde_uid, Comment},
		}

	}

	fileid := dict.GoogleSheetID[VendorID]

	//UPLOADIND ATTENDEE DATA TO SPREADSHEET
	err = utilsheets.AppendData(srv, fileid, "Sheet1!A2", values)
	if err != nil {
		log.Fatal(err)
	}

	// spreadsheet, err := createNewSpreadsheet(srv)
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// fmt.Println(spreadsheet.SpreadsheetId, spreadsheet.SpreadsheetUrl)

	// changePerms(drv, fileid)

}

func UploadCSV(c *gin.Context) {
	r := c.Request
	w := c.Writer

	w.Header().Set("Content-Type", "application/json")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Fatal(err)
	}

	file, _, err := r.FormFile("upload-csv")
	if err != nil {
		log.Fatal()
	}

	defer file.Close()

	reader := csv.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}


	query :=os.Getenv("UPLOAD_LEADGEN_QUERY")
	fmt.Println(query)
	fmt.Println(reflect.TypeOf(query))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)

		_, err = db.Exec(query, record[0], helper.GenerateSecurityToken(record[2]), record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10])
		if err != nil {
			fmt.Println("error occuered")
			log.Fatal(err)
		}
		fmt.Println()
	}
	fmt.Println(reflect.TypeOf(file))

}
