package controller

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"leadgen/dict"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var db *sql.DB
var srv *sheets.Service
var drv *drive.Service

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

	drv, err = drive.NewService(ctx, option.WithCredentialsFile(credentialFile))
	if err != nil {
		log.Fatal(err)
	}

}

func GenerateLead(w http.ResponseWriter, r *http.Request) {
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

	err = appendData(srv, fileid, "Sheet1!A2", values)
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

func UploadCSV() {

}

func generateSecurityToken(input string) string {
	// Convert input to bytes
	temp := []byte(input)

	// Generate MD5 hash
	hash := md5.New()
	hash.Write(temp)
	tokenBytes := hash.Sum(nil)

	// Encode hash using Base64
	tokenBase64 := base64.StdEncoding.EncodeToString(tokenBytes)

	// Remove the last two characters
	if len(tokenBase64) > 2 {
		tokenBase64 = tokenBase64[:len(tokenBase64)-2]
	}

	// Remove non-alphanumeric characters
	re := regexp.MustCompile("[^A-Za-z0-9]+")
	token := re.ReplaceAllString(tokenBase64, "")

	return token
}

func appendData(sheetService *sheets.Service, spreadsheedID string, rangeA1 string, values [][]interface{}) error {
	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}

	_, err := sheetService.Spreadsheets.Values.Append(spreadsheedID, rangeA1, valueRange).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		return err
	}

	return nil
}

func createNewSpreadsheet(s *sheets.Service) (*sheets.Spreadsheet, error) {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: "test",
		},
	}

	call := s.Spreadsheets.Create(spreadsheet)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func changePerms(driveService *drive.Service, fileID string) error {
	permissions := []*drive.Permission{
		{
			Type:         "user",
			Role:         "writer",
			EmailAddress: "zeldon.zoom@gmail.com",
		},
	}

	for _, perm := range permissions {
		_, err := driveService.Permissions.Create(fileID, perm).Do()
		if err != nil {
			return err
		}
		fmt.Println("success")
	}
	return nil
}
