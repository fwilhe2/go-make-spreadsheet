package main

import (
	"encoding/json"
	"fmt"
	spreadsheetlib "github.com/fwilhe2/go-spreadsheetlib"
	"os"
)

type Root struct {
	BankData struct {
		AccountHolder struct {
			AccountDetails struct {
				Accounts json.RawMessage `json:"accounts"` // Just grab the accounts part
			} `json:"accountDetails"`
		} `json:"accountHolder"`
	} `json:"bankData"`
}

type Transaction struct {
	ID     string `json:"id"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	dat, err := os.ReadFile("input.json")
	check(err)
	jsonData := string(dat)

	// First, unmarshal just the path down to accounts
	var root Root
	if err := json.Unmarshal([]byte(jsonData), &root); err != nil {
		panic(err)
	}

	// Unmarshal accounts array
	var accounts []map[string]json.RawMessage
	if err := json.Unmarshal(root.BankData.AccountHolder.AccountDetails.Accounts, &accounts); err != nil {
		panic(err)
	}

	inputCells := [][]spreadsheetlib.Cell{
		{spreadsheetlib.MakeCell("Id", "string"), spreadsheetlib.MakeCell("Value", "string")},
	}

	// Loop through accounts and parse transactions
	for i, acc := range accounts {
		var transactionsWrapper struct {
			List []struct {
				TransactionDetails Transaction `json:"transactionDetails"`
			} `json:"list"`
		}
		if err := json.Unmarshal(acc["transactions"], &transactionsWrapper); err != nil {
			panic(err)
		}

		fmt.Printf("Account %d Transactions:\n", i+1)
		for _, tx := range transactionsWrapper.List {
			fmt.Printf("  ID: %s, Amount: %s %s\n", tx.TransactionDetails.ID, tx.TransactionDetails.Amount.Value, tx.TransactionDetails.Amount.Currency)
			inputCells = append(inputCells, []spreadsheetlib.Cell{spreadsheetlib.MakeCell(tx.TransactionDetails.ID, "string"), spreadsheetlib.MakeCell(tx.TransactionDetails.Amount.Value, "currency")})
		}
	}

	spreadsheet := spreadsheetlib.MakeSpreadsheet(inputCells)

	os.Mkdir("output", 0777)

	buff := spreadsheetlib.MakeOds(spreadsheet)

	archive, err := os.Create(fmt.Sprintf("output/%s.%s", "spreadsheetlib", "ods"))
	if err != nil {
		panic(err)
	}

	archive.Write(buff.Bytes())

	archive.Close()

}
