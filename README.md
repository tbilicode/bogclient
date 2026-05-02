# bogclient

Bank of Georgia API client

## CLI

```sh
Usage: bog <command> [flags]

BOG client

Flags:
  -h, --help                                     Show context-sensitive help.
  -D, --debug                                    Enable debug mode
      --o="table"                                Print output format: json|yaml|table
      --cfg="~/.config/bogclient/config.yaml"    Configuration file
      --storage="~/.config/bogclient"            flag specifies to override default location: ~/.config/bogclient. Use BOG_STORAGE environment to override
      --timeout=6                                Connection timeout

Commands:
  account statement               create statement
  account balance                 prints account balance
  account translate               translate statement to English, requires GOOGLE API KEY
  account convert transactions    convert transactions to CSV or Excel
  account convert daily           convert daily summaries to CSV or Excel
  account convert global          convert global summaries to CSV or Excel

Run "bog <command> --help" for more information on a command.
```

Example:

```sh
bin/bog account statement --from 2026-04-01 --to 2026-04-30 --summary --out /tmp/bog04.json
bin/bog account convert daily /tmp/bog04.json /tmp/bog04-daily.xlsx --format excel
bin/bog account convert global /tmp/bog04.json /tmp/bog04-global.xlsx --format excel

bin/bog account translate /tmp/bog04.json /tmp/bog04-ai.eng.json --provider ai
bin/bog account convert transactions /tmp/bog04-ai.eng.json /tmp/bog04-ai.eng.xlsx --format excel
```
