
# DisGoIRC

Discord-IRC bridge bot.

- Strips IRC color codes, converts IRC format codes (bold, italics, underline)
- Supports multiple Discord servers bridging to one IRC server, under one Discord bot user

## Running the bot

Requires a current version of Go installed (1.8+ recommended; 1.7 or below *may* work, but are untested)

Create a configuration in `conf.json` before starting the bot; an example configuration is provided in `conf.json.example`.  
You may specify an alternate configuration file if desired with `-config=<file>`.

Start the bot: `go run disgoirc.go`  
Start the bot in debug mode: `go run disgoirc.go -debug`

## Contributing

Pull requests are appreciated.  
Please make sure to lint your code; CI will fail any commits which do not pass `make lint`.
