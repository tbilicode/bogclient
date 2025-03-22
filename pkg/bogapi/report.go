package bogapi

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/xuri/excelize/v2"
)

/*
Transaction represents a single transaction record,
it should produce a CSV line with the following fields:

Date,
Doc N,
Loro Account,
Debit,
Credit,
Rate,
Debit Amount in Gel,
Credit Amount in Gel,
Entry Comment,
Operation Type,
Operation ID,
Ref,
Sender Name,
Sender Number Taxpayer,
Sender Account N,
Sender Bank Code,
Sender Bank Name,
Recipient Name,
Recipient Number Taxpayer,
Recipient Account N,
Recipient Bank Code,
Recipient Bank Name,
Nomination,
Additional Info,
Amount,
Amount in Gel,
Turnover Debit,
Turnover Credit,
Turnover Debit  in Gel,
Turnover Credit in Gel,
Balance at end of day,
Balance at end of day in Gel,
Balance
*/
type Transaction struct {
	Date                    string  `json:"Date" csv:"Date"`
	DocumentNumber          string  `json:"DocumentNumber" csv:"Doc N"`
	OperationID             uint64  `json:"OperationID" csv:"Operation ID"`
	OperationType           string  `json:"OperationType" csv:"Operation Type"`
	Account                 string  `json:"Account" csv:"Account"`
	Currency                string  `json:"Currency" csv:"Currency"`
	LoroAccount             string  `json:"LoroAccount" csv:"Loro Account"`
	Debit                   float64 `json:"Debit" csv:"Debit"`
	Credit                  float64 `json:"Credit" csv:"Credit"`
	Rate                    float64 `json:"Rate" csv:"Rate"`
	DebitAmountInGel        float64 `json:"DebitAmountInGel" csv:"Debit Amount in Gel"`
	CreditAmountInGel       float64 `json:"CreditAmountInGel" csv:"Credit Amount in Gel"`
	EntryComment            string  `json:"EntryComment" csv:"Entry Comment"`
	Ref                     string  `json:"Ref" csv:"Ref"`
	SenderName              string  `json:"SenderName" csv:"Sender Name"`
	SenderNumberTaxpayer    string  `json:"SenderNumberTaxpayer" csv:"Sender Number Taxpayer"`
	SenderAccountN          string  `json:"SenderAccountN" csv:"Sender Account N"`
	SenderBankCode          string  `json:"SenderBankCode" csv:"Sender Bank Code"`
	SenderBankName          string  `json:"SenderBankName" csv:"Sender Bank Name"`
	RecipientName           string  `json:"RecipientName" csv:"Recipient Name"`
	RecipientNumberTaxpayer string  `json:"RecipientNumberTaxpayer" csv:"Recipient Number Taxpayer"`
	RecipientAccountN       string  `json:"RecipientAccountN" csv:"Recipient Account N"`
	RecipientBankCode       string  `json:"RecipientBankCode" csv:"Recipient Bank Code"`
	RecipientBankName       string  `json:"RecipientBankName" csv:"Recipient Bank Name"`
	Nomination              string  `json:"Nomination" csv:"Nomination"`
	AdditionalInfo          string  `json:"AdditionalInfo" csv:"Additional Info"`
	Amount                  float64 `json:"Amount" csv:"Amount"`
	AmountInGel             float64 `json:"AmountInGel" csv:"Amount in Gel"`
	TurnoverDebit           float64 `json:"TurnoverDebit" csv:"Turnover Debit"`
	TurnoverCredit          float64 `json:"TurnoverCredit" csv:"Turnover Credit"`
	TurnoverDebitInGel      float64 `json:"TurnoverDebitInGel" csv:"Turnover Debit in Gel"`
	TurnoverCreditInGel     float64 `json:"TurnoverCreditInGel" csv:"Turnover Credit in Gel"`
	BalanceAtEndOfDay       float64 `json:"BalanceAtEndOfDay" csv:"Balance at end of day"`
	BalanceAtEndOfDayInGel  float64 `json:"BalanceAtEndOfDayInGel" csv:"Balance at end of day in Gel"`
	Balance                 float64 `json:"Balance" csv:"Balance"`

	Recort Record `json:"-"`
}

type TransactionSlice []Transaction

func (t TransactionSlice) Len() int {
	return len(t)
}

func (t TransactionSlice) Less(i, j int) bool {
	return t[i].Date < t[j].Date
}

func (t TransactionSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TransactionSlice) ToCSV(w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write CSV header
	header := []string{
		"Date", "Doc N", "Operation ID", "Operation Type", "Account", "Currency", "Loro Account", "Debit", "Credit", "Rate",
		"Debit Amount in Gel", "Credit Amount in Gel", "Entry Comment",
		"Ref", "Sender Name", "Sender Number Taxpayer", "Sender Account N",
		"Sender Bank Code", "Sender Bank Name", "Recipient Name", "Recipient Number Taxpayer",
		"Recipient Account N", "Recipient Bank Code", "Recipient Bank Name", "Nomination",
		"Additional Info", "Amount", "Amount in Gel", "Turnover Debit", "Turnover Credit",
		"Turnover Debit in Gel", "Turnover Credit in Gel", "Balance at end of day",
		"Balance at end of day in Gel", "Balance",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write CSV rows
	for _, transaction := range t {
		row := []string{
			transaction.Date,
			transaction.DocumentNumber,
			formatUInt(transaction.OperationID),
			transaction.OperationType,
			transaction.Account,
			transaction.Currency,
			transaction.LoroAccount,
			formatFloat(transaction.Debit),
			formatFloat(transaction.Credit),
			formatFloat(transaction.Rate),
			formatFloat(transaction.DebitAmountInGel),
			formatFloat(transaction.CreditAmountInGel),
			transaction.EntryComment,
			transaction.Ref,
			transaction.SenderName,
			transaction.SenderNumberTaxpayer,
			transaction.SenderAccountN,
			transaction.SenderBankCode,
			transaction.SenderBankName,
			transaction.RecipientName,
			transaction.RecipientNumberTaxpayer,
			transaction.RecipientAccountN,
			transaction.RecipientBankCode,
			transaction.RecipientBankName,
			transaction.Nomination,
			transaction.AdditionalInfo,
			formatFloat(transaction.Amount),
			formatFloat(transaction.AmountInGel),
			formatFloat(transaction.TurnoverDebit),
			formatFloat(transaction.TurnoverCredit),
			formatFloat(transaction.TurnoverDebitInGel),
			formatFloat(transaction.TurnoverCreditInGel),
			formatFloat(transaction.BalanceAtEndOfDay),
			formatFloat(transaction.BalanceAtEndOfDayInGel),
			formatFloat(transaction.Balance),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func formatUInt(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func (t TransactionSlice) Dedup() TransactionSlice {
	transactionMap := make(map[string]Transaction)
	var transactions TransactionSlice

	for _, transaction := range t {
		key := transaction.DocumentNumber + transaction.LoroAccount + transaction.Date
		if _, exists := transactionMap[key]; !exists {
			transactionMap[key] = transaction
			transactions = append(transactions, transaction)
		}
	}

	return transactions
}

func Report(r *AccountStatements) TransactionSlice {
	var transactions TransactionSlice

	for _, accountStatement := range r.Combined {
		for _, record := range accountStatement.Records {
			transaction := Transaction{
				Date:                    record.EntryDate.String(),
				DocumentNumber:          record.EntryDocumentNumber,
				Account:                 accountStatement.Account,
				Currency:                accountStatement.Currency,
				LoroAccount:             record.EntryAccountNumber,
				Debit:                   record.EntryAmountDebit,
				Credit:                  record.EntryAmountCredit,
				Rate:                    record.DocumentRate,
				DebitAmountInGel:        record.EntryAmountDebitBase,
				CreditAmountInGel:       record.EntryAmountCreditBase,
				EntryComment:            record.EntryComment,
				OperationType:           record.DocumentProductGroup,
				OperationID:             uint64(record.EntryId),
				Ref:                     record.EntryDocumentNumber,
				SenderName:              record.SenderDetails.Name,
				SenderNumberTaxpayer:    record.SenderDetails.Inn,
				SenderAccountN:          record.SenderDetails.AccountNumber,
				SenderBankCode:          record.SenderDetails.BankCode,
				SenderBankName:          record.SenderDetails.BankName,
				RecipientName:           record.BeneficiaryDetails.Name,
				RecipientNumberTaxpayer: record.BeneficiaryDetails.Inn,
				RecipientAccountN:       record.BeneficiaryDetails.AccountNumber,
				RecipientBankCode:       record.BeneficiaryDetails.BankCode,
				RecipientBankName:       record.BeneficiaryDetails.BankName,
				Nomination:              record.DocumentNomination,
				AdditionalInfo:          record.DocumentInformation,
				Amount:                  record.EntryAmount,
				AmountInGel:             record.EntryAmountBase,
				TurnoverDebit:           record.EntryAmountDebit,
				TurnoverCredit:          record.EntryAmountCredit,
				TurnoverDebitInGel:      record.EntryAmountDebitBase,
				TurnoverCreditInGel:     record.EntryAmountCreditBase,
				BalanceAtEndOfDay:       record.EntryAmount,
				BalanceAtEndOfDayInGel:  record.EntryAmountBase,
				Balance:                 record.EntryAmount,
			}

			transactions = append(transactions, transaction)
		}
	}

	sort.Slice(transactions, func(i, j int) bool {
		if transactions[i].Date == transactions[j].Date {
			return transactions[i].OperationID < transactions[j].OperationID
		}
		return transactions[i].Date < transactions[j].Date
	})

	return transactions
}

func (t TransactionSlice) ToExcel(w io.Writer) error {
	f := excelize.NewFile()
	sheet := "Statement of Accounts"
	_ = f.SetSheetName(f.GetSheetName(0), sheet)

	// Write Excel header
	header := []string{
		"Date", "Doc N", "Operation ID", "Operation Type", "Account", "Currency", "Loro Account", "Debit", "Credit", "Rate",
		"Debit Amount in Gel", "Credit Amount in Gel", "Entry Comment",
		"Ref", "Sender Name", "Sender Number Taxpayer", "Sender Account N",
		"Sender Bank Code", "Sender Bank Name", "Recipient Name", "Recipient Number Taxpayer",
		"Recipient Account N", "Recipient Bank Code", "Recipient Bank Name", "Nomination",
		"Additional Info", "Amount", "Amount in Gel", "Turnover Debit", "Turnover Credit",
		"Turnover Debit in Gel", "Turnover Credit in Gel", "Balance at end of day",
		"Balance at end of day in Gel", "Balance",
	}
	for i, h := range header {
		col, _ := excelize.ColumnNumberToName(i + 1)
		_ = f.SetCellValue(sheet, col+"1", h)
	}

	// Write Excel rows
	for i, transaction := range t {
		row := []any{
			transaction.Date,
			transaction.DocumentNumber,
			formatUInt(transaction.OperationID),
			transaction.OperationType,
			transaction.Account,
			transaction.Currency,
			transaction.LoroAccount,
			transaction.Debit,
			transaction.Credit,
			transaction.Rate,
			transaction.DebitAmountInGel,
			transaction.CreditAmountInGel,
			transaction.EntryComment,
			transaction.Ref,
			transaction.SenderName,
			transaction.SenderNumberTaxpayer,
			transaction.SenderAccountN,
			transaction.SenderBankCode,
			transaction.SenderBankName,
			transaction.RecipientName,
			transaction.RecipientNumberTaxpayer,
			transaction.RecipientAccountN,
			transaction.RecipientBankCode,
			transaction.RecipientBankName,
			transaction.Nomination,
			transaction.AdditionalInfo,
			transaction.Amount,
			transaction.AmountInGel,
			transaction.TurnoverDebit,
			transaction.TurnoverCredit,
			transaction.TurnoverDebitInGel,
			transaction.TurnoverCreditInGel,
			transaction.BalanceAtEndOfDay,
			transaction.BalanceAtEndOfDayInGel,
			transaction.Balance,
		}
		for j, value := range row {
			col, _ := excelize.ColumnNumberToName(j + 1)
			cell := fmt.Sprintf("%s%d", col, i+2)
			_ = f.SetCellValue(sheet, cell, value)
		}
	}

	if err := f.Write(w); err != nil {
		return err
	}

	return nil
}
