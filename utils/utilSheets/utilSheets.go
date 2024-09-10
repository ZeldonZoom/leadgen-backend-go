package utilsheets

import (
	"context"
	"leadgen/dict"
	"log"

	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var Srv *sheets.Service
var Drv *drive.Service

func init() {
	var err error

	ctx := context.Background()
	credentialFile := dict.SheetsCred["test"]
	Srv, err = sheets.NewService(ctx, option.WithCredentialsFile(credentialFile))
	if err != nil {
		log.Fatal(err)
	}

	Drv, err = drive.NewService(ctx, option.WithCredentialsFile(credentialFile))
	if err != nil {
		log.Fatal(err)
	}

}

func CreateNewSpreadsheet(s *sheets.Service) (*sheets.Spreadsheet, error) {
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

func AppendData(sheetService *sheets.Service, spreadsheedID string, rangeA1 string, values [][]interface{}) error {
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

// func changePerms(driveService *drive.Service, fileID string) error {
// 	permissions := []*drive.Permission{
// 		{
// 			Type:         "user",
// 			Role:         "writer",
// 			EmailAddress: "zeldon.zoom@gmail.com",
// 		},
// 	}

// 	for _, perm := range permissions {
// 		_, err := driveService.Permissions.Create(fileID, perm).Do()
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("success")
// 	}
// 	return nil
// }
