package bogapi

import (
	"time"

	"github.com/pkg/errors"
)

type Time time.Time

func (ct *Time) UnmarshalJSON(b []byte) error {
	// Trim quotes from JSON string
	str := string(b)
	if str == "null" {
		return nil
	}
	if len(str) > 2 && str[0] == '"' {
		str = str[1 : len(str)-1]
	}

	// Define possible time formats
	formats := []string{
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		tim, err := time.Parse(format, str)
		if err == nil {
			*ct = Time(tim)
			return nil
		}
	}
	return errors.Errorf("unable to parse time: %s", str)
}

type GlobalSummary struct {
	AccountNumber   string  `json:"AccountNumber"`
	Currency        string  `json:"Currency"`
	StartDate       Time    `json:"StartDate"`
	EndDate         Time    `json:"EndDate"`
	PeriodStartDate Time    `json:"PeriodStartDate"`
	PeriodEndDate   Time    `json:"PeriodEndDate"`
	InAmount        float64 `json:"InAmount"`
	InAmountBase    float64 `json:"InAmountBase"`
	InRate          float64 `json:"InRate"`
	OutAmount       float64 `json:"OutAmount"`
	OutAmountBase   float64 `json:"OutAmountBase"`
	OutRate         float64 `json:"OutRate"`
	CreditSum       float64 `json:"CreditSum"`
	DebitSum        float64 `json:"DebitSum"`
}

type DailySummary struct {
	Balance     float64 `json:"Balance"`
	BalanceBase float64 `json:"BalanceBase"`
	CreditSum   float64 `json:"CreditSum"`
	DebitSum    float64 `json:"DebitSum"`
	Rate        float64 `json:"Rate"`
	EntryCount  int     `json:"EntryCount"`
	Date        Time    `json:"Date"`
}

type StatementSummary struct {
	GlobalSummary  GlobalSummary  `json:"GlobalSummary"`
	DailySummaries []DailySummary `json:"DailySummaries"`
}

type AccountBalance struct {
	AvailableBalance float64 `json:"AvailableBalance"`
	CurrentBalance   float64 `json:"CurrentBalance"`
}

type SenderDetails struct {
	Name          string `json:"Name"`
	Inn           string `json:"Inn"`
	AccountNumber string `json:"AccountNumber"`
	BankCode      string `json:"BankCode"`
	BankName      string `json:"BankName"`
}

type BeneficiaryDetails struct {
	Name          string `json:"Name"`
	Inn           string `json:"Inn"`
	AccountNumber string `json:"AccountNumber"`
	BankCode      string `json:"BankCode"`
	BankName      string `json:"BankName"`
}

type Record struct {
	EntryDate                          Time               `json:"EntryDate"`
	EntryDocumentNumber                string             `json:"EntryDocumentNumber"`
	EntryAccountNumber                 string             `json:"EntryAccountNumber"`
	EntryAmountDebit                   float64            `json:"EntryAmountDebit"`
	EntryAmountDebitBase               float64            `json:"EntryAmountDebitBase"`
	EntryAmountCredit                  float64            `json:"EntryAmountCredit"`
	EntryAmountCreditBase              float64            `json:"EntryAmountCreditBase"`
	EntryAmountBase                    float64            `json:"EntryAmountBase"`
	EntryAmount                        float64            `json:"EntryAmount"`
	EntryComment                       string             `json:"EntryComment"`
	EntryDepartment                    string             `json:"EntryDepartment"`
	EntryAccountPoint                  string             `json:"EntryAccountPoint"`
	DocumentProductGroup               string             `json:"DocumentProductGroup"`
	DocumentValueDate                  Time               `json:"DocumentValueDate"`
	SenderDetails                      SenderDetails      `json:"SenderDetails"`
	BeneficiaryDetails                 BeneficiaryDetails `json:"BeneficiaryDetails"`
	DocumentTreasuryCode               string             `json:"DocumentTreasuryCode"`
	DocumentNomination                 string             `json:"DocumentNomination"`
	DocumentInformation                string             `json:"DocumentInformation"`
	DocumentSourceAmount               float64            `json:"DocumentSourceAmount"`
	DocumentSourceCurrency             string             `json:"DocumentSourceCurrency"`
	DocumentDestinationAmount          float64            `json:"DocumentDestinationAmount"`
	DocumentDestinationCurrency        string             `json:"DocumentDestinationCurrency"`
	DocumentReceiveDate                Time               `json:"DocumentReceiveDate"`
	DocumentBranch                     string             `json:"DocumentBranch"`
	DocumentDepartment                 string             `json:"DocumentDepartment"`
	DocumentActualDate                 Time               `json:"DocumentActualDate"`
	DocumentExpiryDate                 Time               `json:"DocumentExpiryDate"`
	DocumentRateLimit                  float64            `json:"DocumentRateLimit"`
	DocumentRate                       float64            `json:"DocumentRate"`
	DocumentRegistrationRate           float64            `json:"DocumentRegistrationRate"`
	DocumentSenderInstitution          string             `json:"DocumentSenderInstitution"`
	DocumentIntermediaryInstitution    string             `json:"DocumentIntermediaryInstitution"`
	DocumentBeneficiaryInstitution     string             `json:"DocumentBeneficiaryInstitution"`
	DocumentPayee                      string             `json:"DocumentPayee"`
	DocumentCorrespondentAccountNumber string             `json:"DocumentCorrespondentAccountNumber"`
	DocumentCorrespondentBankCode      string             `json:"DocumentCorrespondentBankCode"`
	DocumentCorrespondentBankName      string             `json:"DocumentCorrespondentBankName"`
	DocumentKey                        float64            `json:"DocumentKey"`
	EntryId                            float64            `json:"EntryId"`
	DocumentPayerName                  string             `json:"DocumentPayerName"`
	DocumentPayerInn                   string             `json:"DocumentPayerInn"`
	DocComment                         string             `json:"DocComment"`
}

type StatementResponse struct {
	ID      int      `json:"Id"`
	Count   int      `json:"Count"`
	Records []Record `json:"Records"`
}

type AccountStatement struct {
	Account   string `json:"Account"`
	Currency  string `json:"Currency"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`

	StatementID int               `json:"StatementID"`
	Records     []Record          `json:"Records"`
	Summary     *StatementSummary `json:"Summary"`
}

type AccountStatements struct {
	Combined []*AccountStatement `json:"Combined"`
}
