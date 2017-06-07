
# DisGoIRC

Discord-IRC bridge bot.

Edit the configuration in `conf.json` before starting the bot; an example configuration is provided in `conf.json.example`.

- Strips IRC color codes, converts IRC format codes (bold, italics, underline)
- Supports multiple Discord servers bridging to one IRC server, under one Discord bot user

Requires a current version of golang installed. Start the bot by executing `go run disgoirc.go`
