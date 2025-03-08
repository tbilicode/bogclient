package account

import (
	"errors"

	"github.com/tbilicode/bogclient/internal/cli"
	"github.com/tbilicode/bogclient/pkg/bogapi"
)

type Cmd struct {
	Statement StatementCmd `cmd:"" help:"create statement"`
	Balance   BalanceCmd   `cmd:"" help:"prints account balance"`
}

// StatementCmd create statement
type StatementCmd struct {
	Account  string `help:"Filter by account, empty for all"`
	Currency string `help:"Filter by currency, empty for all"`
	Month    int    `help:"month to summarize, in 1-12 format"`
	From     string `help:"start date"`
	To       string `help:"end date"`
	Summary  bool   `help:"add summary"`
}

func (cmd *StatementCmd) Run(ctx *cli.Cli) error {
	if cmd.Month == 0 && (cmd.From == "" || cmd.To == "") {
		return errors.New("either month or start and end dates must be provided")
	}
	if cmd.Month != 0 && (cmd.From != "" || cmd.To != "") {
		return errors.New("either month or start and end dates must be provided")
	}

	req := &bogapi.StatementRequest{
		StartDate: cmd.From,
		EndDate:   cmd.To,
		Account:   cmd.Account,
		Currency:  cmd.Currency,
		Summary:   cmd.Summary,
	}

	if cmd.Month != 0 {
		req.StartDate, req.EndDate = bogapi.MonthRange(cmd.Month)
	}

	client, err := ctx.Client()
	if err != nil {
		return err
	}

	res, err := client.AllStatements(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Print(res)
}

// BalanceCmd prints account balance
type BalanceCmd struct {
}

func (cmd *BalanceCmd) Run(ctx *cli.Cli) error {
	client, err := ctx.Client()
	if err != nil {
		return err
	}

	res, err := client.AllBalances(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.Print(res)
}
