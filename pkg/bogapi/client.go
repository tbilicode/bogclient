package bogapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/effective-security/porto/pkg/retriable"
	"github.com/effective-security/porto/xhttp/header"
	"github.com/effective-security/x/values"
	"github.com/effective-security/xlog"
	"github.com/pkg/errors"
	"github.com/tbilicode/bogclient/internal/version"
	"golang.org/x/net/context"
)

var logger = xlog.NewPackageLogger("github.com/tbilicode/bogclient/pkg", "bogapi")

// NowFunc is a function that returns the current time
var NowFunc = time.Now

type Client interface {
	Accounts() []Account
	Authenticate(ctx context.Context) error
	Statement(ctx context.Context, req *StatementRequest) (*StatementResponse, error)
	AllStatements(ctx context.Context, req *StatementRequest) (*AccountStatements, error)
	StatementSummary(ctx context.Context, account, currency string, id int) (*StatementSummary, error)
	Balance(ctx context.Context, account, currency string) (*AccountBalance, error)
	AllBalances(ctx context.Context) (map[string]*AccountBalance, error)
}

type client struct {
	cfg        *Config
	httpClient *retriable.Client
	auth       *AuthResponse
}

func CreateClient(file string, timeoutSec int) (Client, error) {
	cfg, err := LoadConfig(file)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load config")
	}

	server := values.StringsCoalesce(cfg.ApiHost, os.Getenv("BOG_SERVER"))
	client, err := retriable.Default(server)
	if err != nil {
		return nil, err
	}

	if timeoutSec > 0 {
		client.WithTimeout(time.Second * time.Duration(timeoutSec))
	}

	client = client.WithUserAgent("tbilicode-bogclient " + version.Current().String())

	return NewClient(cfg, client), nil
}

func NewClient(cfg *Config, httpClient *retriable.Client) Client {
	c := &client{
		cfg:        cfg,
		httpClient: httpClient,
	}
	return c
}

func (c *client) HTTPClient() *retriable.Client {
	return c.httpClient
}

func (c *client) Accounts() []Account {
	return c.cfg.Accounts
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`

	Expires time.Time `json:"-"`
}

func (c *client) Authenticate(ctx context.Context) error {
	if c.auth != nil && time.Now().Before(c.auth.Expires) {
		return nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.cfg.ClientID)
	data.Set("client_secret", c.cfg.ClientSecret)

	req, err := http.NewRequest(http.MethodPost, c.cfg.AuthURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set(header.ContentType, "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.cfg.ClientID, c.cfg.ClientSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ContextKV(ctx, xlog.ERROR,
			"err", err.Error(),
		)
		return errors.WithMessage(err, "failed to authenticate")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("authentication failed: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.WithMessage(err, "failed to read response body")
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return errors.WithMessage(err, "failed to parse response body")
	}

	authResp.Expires = time.Now().Add(time.Second * time.Duration(authResp.ExpiresIn))
	c.auth = &authResp

	c.httpClient.AddHeader(header.Authorization, authResp.TokenType+" "+authResp.AccessToken)
	c.httpClient.AddHeader(header.ContentType, "application/json")
	c.httpClient.AddHeader("Accept-Language", "en")
	return nil
}

func (c *client) StatementSummary(ctx context.Context, account, currency string, id int) (*StatementSummary, error) {
	err := c.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/statement/summary/%s/%s/%d", account, currency, id)
	var summary StatementSummary
	hdr, status, err := c.httpClient.Get(ctx, path, &summary)
	if err != nil {
		logger.ContextKV(ctx, xlog.ERROR,
			"account", account,
			"currency", currency,
			"status", status,
			"header", hdr,
			"err", err.Error(),
		)
		return nil, errors.WithMessagef(err, "failed to get statement summary: %s %s", account, currency)
	}

	return &summary, err
}

type StatementRequest struct {
	Account  string
	Currency string
	// StartDate specifies the start date for the statement period
	StartDate string
	// EndDate specifies the end date for the statement period
	EndDate string
	Summary bool
}

func (c *client) Statement(ctx context.Context, req *StatementRequest) (*StatementResponse, error) {
	err := c.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/statement/%s/%s/%s/%s", req.Account, req.Currency, req.StartDate, req.EndDate)
	var res StatementResponse
	hdr, status, err := c.httpClient.Get(ctx, path, &res)
	if err != nil {
		logger.ContextKV(ctx, xlog.ERROR,
			"account", req.Account,
			"currency", req.Currency,
			"status", status,
			"header", hdr,
			"err", err.Error(),
		)
		return nil, errors.WithMessagef(err, "failed to create statement: %s %s - [%s,%s]",
			req.Account, req.Currency, req.StartDate, req.EndDate)
	}
	return &res, err
}

func (c *client) AllStatements(ctx context.Context, req *StatementRequest) (*AccountStatements, error) {
	res := AccountStatements{}
	for _, acc := range c.cfg.Accounts {
		if req.Account != "" && req.Account != acc.ID {
			continue
		}
		for _, currency := range acc.Currency {
			if req.Currency != "" && req.Currency != currency {
				continue
			}

			st, err := c.Statement(ctx, &StatementRequest{
				Account:   acc.ID,
				Currency:  currency,
				StartDate: req.StartDate,
				EndDate:   req.EndDate,
			})
			if err != nil {
				return nil, err
			}

			ast := &AccountStatement{
				Account:     acc.ID,
				Currency:    currency,
				StartDate:   req.StartDate,
				EndDate:     req.EndDate,
				StatementID: st.ID,
				Records:     st.Records,
			}

			if req.Summary {
				sum, err := c.StatementSummary(ctx, acc.ID, currency, st.ID)
				if err != nil {
					return nil, err
				}
				ast.Summary = sum
			}
			res.Combined = append(res.Combined, ast)
		}
	}
	return &res, nil
}

func (c *client) Balance(ctx context.Context, account, currency string) (*AccountBalance, error) {
	err := c.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/accounts/%s/%s", account, currency)
	var balance AccountBalance
	hdr, status, err := c.httpClient.Get(ctx, path, &balance)
	if err != nil {
		logger.ContextKV(ctx, xlog.ERROR,
			"account", account,
			"currency", currency,
			"status", status,
			"header", hdr,
			"err", err.Error(),
		)
		return nil, errors.WithMessagef(err, "failed to get balance: %s %s", account, currency)
	}
	return &balance, err
}

func (c *client) AllBalances(ctx context.Context) (map[string]*AccountBalance, error) {
	summary := make(map[string]*AccountBalance)
	for _, acc := range c.cfg.Accounts {
		for _, currency := range acc.Currency {
			sum, err := c.Balance(ctx, acc.ID, currency)
			if err != nil {
				return nil, err
			}
			summary[acc.ID+" "+currency] = sum
		}
	}
	return summary, nil
}

func MonthRange(month int) (string, string) {
	now := NowFunc()
	start := time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	return start.Format("2006-01-02"), end.Format("2006-01-02")
}
