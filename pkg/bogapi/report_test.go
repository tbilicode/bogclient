package bogapi_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tbilicode/bogclient/pkg/bogapi"
)

func TestReport(t *testing.T) {
	data, err := os.ReadFile("testdata/statement_feb.json")
	require.NoError(t, err)

	var res bogapi.AccountStatements
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)

	transactions := bogapi.Report(&res)

	// Create a new CSV file
	file, err := os.Create("testdata/statement_feb.csv")
	require.NoError(t, err)
	defer file.Close()

	require.NoError(t, transactions.ToCSV(file))
}
