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
  account statement    create statement
  account balance      prints account balance
  account translate    translate statement to English, requires GOOGLE API KEY
  account convert      convert statement to CSV or Excel

Run "bog <command> --help" for more information on a command.
```
