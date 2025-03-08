package cli

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/effective-security/porto/xhttp/correlation"
	"github.com/effective-security/x/ctl"
	"github.com/effective-security/x/values"
	"github.com/effective-security/xlog"
	"github.com/mitchellh/go-homedir"
	"github.com/tbilicode/bogclient/pkg/bogapi"
	"github.com/tbilicode/bogclient/pkg/print"
)

var logger = xlog.NewPackageLogger("github.com/tbilicode/bogclient/internal", "cli")

var (
	// DefaultStoragePath specifies default storage path
	DefaultStoragePath = "~/.config/bogclient"
)

// Cli provides CLI context to run commands
type Cli struct {
	Version ctl.VersionFlag `name:"version" help:"Print version information and quit" hidden:""`
	Debug   bool            `short:"D" help:"Enable debug mode"`
	O       string          `help:"Print output format: json|yaml|table" default:"table"`
	Cfg     string          `help:"Configuration file" default:"~/.config/bogclient/config.yaml"`
	Storage string          `help:"flag specifies to override default location: ~/.config/bogclient. Use BOG_STORAGE environment to override" default:"~/.config/bogclient"`
	Timeout int             `help:"Connection timeout"  default:"6"`

	TimeFormat string `name:"time" help:"Print time format: utc|local|ago" hidden:"" default:"utc"`

	// Output is the destination for all output from the command, typically set to os.Stdout
	output io.Writer
	// ErrOutput is the destination for errors.
	// If not set, errors will be written to os.StdError
	errOutput io.Writer

	client bogapi.Client
	ctx    context.Context
}

// Context for requests
func (c *Cli) Context() context.Context {
	if c.ctx == nil {
		c.ctx = correlation.WithMetaFromContext(context.Background())
		logger.ContextKV(c.ctx, xlog.DEBUG, "context", "created")
	}
	return c.ctx
}

// IsJSON returns true if the output format us JSON
func (c *Cli) IsJSON() bool {
	return c.O == "json"
}

// Writer returns a writer for control output
func (c *Cli) Writer() io.Writer {
	if c.output != nil {
		return c.output
	}
	return os.Stdout
}

// WithWriter allows to specify a custom writer
func (c *Cli) WithWriter(out io.Writer) *Cli {
	c.output = out
	return c
}

// ErrWriter returns a writer for control output
func (c *Cli) ErrWriter() io.Writer {
	if c.errOutput != nil {
		return c.errOutput
	}
	return os.Stderr
}

// WithErrWriter allows to specify a custom error writer
func (c *Cli) WithErrWriter(out io.Writer) *Cli {
	c.errOutput = out
	return c
}

// AfterApply hook loads config
func (c *Cli) AfterApply(app *kong.Kong, vars kong.Vars) error {
	xlog.SetFormatter(xlog.NewPrettyFormatter(c.ErrWriter()))
	if c.Debug {
		xlog.SetGlobalLogLevel(xlog.DEBUG)
	} else {
		xlog.SetGlobalLogLevel(xlog.ERROR)
	}

	//print.DefaultTimeFormat = c.TimeFormat
	return nil
}

// Client returns client
func (c *Cli) Client() (bogapi.Client, error) {
	if c.client == nil {
		// expand Storage in order of priorities: flag, Env, config, default
		storage := values.StringsCoalesce(
			c.Storage,
			os.Getenv("BOG_STORAGE"),
			DefaultStoragePath,
		)

		c.Storage, _ = homedir.Expand(storage)

		cfgpath := values.StringsCoalesce(
			c.Cfg,
			filepath.Join(c.Storage, "config.yaml"),
		)

		cfg, _ := homedir.Expand(cfgpath)

		client, err := bogapi.CreateClient(cfg, c.Timeout)
		if err != nil {
			return nil, err
		}
		c.client = client
	}
	return c.client, nil
}

// Print response to out
func (c *Cli) Print(value any) error {
	return print.Object(c.Writer(), c.O, value)
}
