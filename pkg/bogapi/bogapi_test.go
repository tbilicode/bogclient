package bogapi_test

import (
	"context"
	"testing"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tbilicode/bogclient/pkg/bogapi"
)

func Test_LoadConfig(t *testing.T) {
	cfg, err := bogapi.LoadConfig("testdata/config.yaml")
	require.NoError(t, err)
	require.Len(t, cfg.Accounts, 2)

	assert.Equal(t, "https://sandbox.businessonline.ge", cfg.ApiHost)
	assert.Equal(t, "https://sandbox.bog.ge/auth/realms/bog/protocol/openid-connect/token", cfg.AuthURL)
	assert.Equal(t, "123456", cfg.ClientID)
	assert.Equal(t, "abcdef", cfg.ClientSecret)
	assert.Equal(t, "GE12BG0000000106360001", cfg.Accounts[0].ID)
	assert.Equal(t, "Primary", cfg.Accounts[0].Name)
	assert.Equal(t, []string{"USD", "EUR", "GEL"}, cfg.Accounts[0].Currency)
	assert.Equal(t, "GE12BG0000000106360002", cfg.Accounts[1].ID)
	assert.Equal(t, "Card", cfg.Accounts[1].Name)
	assert.Equal(t, []string{"USD", "EUR", "GEL"}, cfg.Accounts[1].Currency)
}

func Test_RealAuth(t *testing.T) {
	// To test your connection, populate the config file with real credentials
	// and uncomment the following line
	t.Skip("Real auth test")

	cf, err := homedir.Expand("~/.config/bogclient/config.yaml")
	require.NoError(t, err)

	client, err := bogapi.CreateClient(cf, 6)
	require.NoError(t, err)

	ctx := context.Background()
	err = client.Authenticate(ctx)
	require.NoError(t, err)
}

func Test_MonthRange(t *testing.T) {
	tcases := []struct {
		month int
		start string
		end   string
	}{
		{1, "2021-01-01", "2021-01-31"},
		{2, "2021-02-01", "2021-02-28"},
		{3, "2021-03-01", "2021-03-31"},
		{4, "2021-04-01", "2021-04-30"},
		{5, "2021-05-01", "2021-05-31"},
		{6, "2021-06-01", "2021-06-30"},
		{7, "2021-07-01", "2021-07-31"},
		{8, "2021-08-01", "2021-08-31"},
		{9, "2021-09-01", "2021-09-30"},
		{10, "2021-10-01", "2021-10-31"},
		{11, "2021-11-01", "2021-11-30"},
		{12, "2021-12-01", "2021-12-31"},
	}

	bogapi.NowFunc = func() time.Time {
		now := time.Now()
		return time.Date(2021, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}
	defer func() {
		bogapi.NowFunc = time.Now
	}()

	for _, tc := range tcases {
		start, end := bogapi.MonthRange(tc.month)
		assert.Equal(t, tc.start, start)
		assert.Equal(t, tc.end, end)
	}
}
