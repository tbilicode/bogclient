package account

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/tbilicode/bogclient/internal/cli"
	"github.com/tbilicode/bogclient/pkg/bogapi"
	"github.com/tbilicode/bogclient/pkg/translate"
)

type Cmd struct {
	Statement StatementCmd `cmd:"" help:"create statement"`
	Balance   BalanceCmd   `cmd:"" help:"prints account balance"`
	Translate TranslateCmd `cmd:"" help:"translate statement to English, requires GOOGLE API KEY"`
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

// StatementCmd create statement
type StatementCmd struct {
	Account  string `help:"Filter by account, empty for all"`
	Currency string `help:"Filter by currency, empty for all"`
	Month    int    `help:"month to summarize, in 1-12 format"`
	From     string `help:"start date"`
	To       string `help:"end date"`
	Summary  bool   `help:"add summary"`
	Out      string `help:"output file, if not provided prints to stdout"`
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

	if cmd.Out != "" {
		return ctx.WriteFile(cmd.Out, res)
	}

	return ctx.Print(res)
}

type TranslateCmd struct {
	In  string `kong:"arg" help:"input file" required:""`
	Out string `kong:"arg" help:"output file" required:""`
}

func (cmd *TranslateCmd) Run(ctx *cli.Cli) error {
	tr := translate.NewTranslator()

	dict := filepath.Join(ctx.Storage, "dict.json")
	dict, _ = homedir.Expand(dict)

	err := tr.LoadDictionary(dict, true)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(cmd.In)
	if err != nil {
		return err
	}

	doc := new(bogapi.AccountStatements)
	err = json.Unmarshal(data, doc)
	if err != nil {
		return err
	}

	texts, err := tr.Extract(doc)
	if err != nil {
		return err
	}

	err = tr.Translate(ctx.Context(), texts)
	if err != nil {
		return err
	}

	_ = tr.SaveDictionary(dict)

	replaced, err := tr.Update(ctx.Context(), doc)
	if err != nil {
		return err
	}
	ctx.Print(replaced)

	data, err = json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(cmd.Out, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
