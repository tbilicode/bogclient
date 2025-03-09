package bogapi_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tbilicode/bogclient/pkg/bogapi"
)

func TestStatementIDResponse(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/statementid_response.json")
	require.NoError(t, err)

	var res bogapi.StatementResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	assert.Equal(t, 10120, res.ID)
	assert.Equal(t, 2, res.Count)
	assert.Equal(t, 2, len(res.Records))
}

func TestAccountBalance(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/account_balance.json")
	require.NoError(t, err)

	var res bogapi.AccountBalance
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	assert.Equal(t, 23083.33, res.AvailableBalance)
	assert.Equal(t, 23083.33, res.CurrentBalance)
}

func TestStatementSummary(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/statement_summary.json")
	require.NoError(t, err)

	var res bogapi.StatementSummary
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	assert.Equal(t, "sample string 1", res.GlobalSummary.AccountNumber)
	assert.Equal(t, 2, len(res.DailySummaries))
}

func TestAccountStatements(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/statement_feb.json")
	require.NoError(t, err)

	var res bogapi.AccountStatements
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	assert.Equal(t, 6, len(res.Combined))
}

func TestStatementResponse(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/raw_statement.json")
	require.NoError(t, err)

	var res bogapi.StatementResponse

	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()

	err = d.Decode(&res)
	require.NoError(t, err)
}
