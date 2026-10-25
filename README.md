# Alert Nearby Bot

A Telegram bot that watches public alert channels and forwards only the
messages mentioning your location (or any configured keyword) to your private
chat.

It uses a Telegram **user client** to read the channels and a **bot** to
deliver the filtered alerts.

## How it works

1. `watcher` logs in as a Telegram user and polls the channels for new posts.
2. `matcher` keeps only posts that contain one of your keywords.
3. `notifier` sends those posts to your chat(s) through a bot.

## Project structure

```
alert-nearby-bot/
├── cmd/bot/main.go          # entry point, wires everything together
├── internal/
│   ├── config/              # loads and validates environment variables
│   ├── watcher/             # reads the channels (gotd user client)
│   ├── matcher/             # keyword filtering (+ unit test)
│   └── notifier/            # sends alerts (bot)
├── go.mod
└── README.md
```

## Configuration

Create a `.env` file in the project root (next to `go.mod`) with:

| Variable | Required | Description |
|----------|----------|-------------|
| `API_ID` | yes | Your api_id (a number) from https://my.telegram.org |
| `API_HASH` | yes | Your api_hash (a string) from https://my.telegram.org |
| `PHONE` | yes | Phone number for the user login, e.g. `+380XXXXXXXXX` |
| `BOT_TOKEN` | yes | Bot token from @BotFather |
| `MY_CHAT_IDS` | yes | Comma-separated chat IDs (from @userinfobot) to send alerts to |
| `CHANNELS` | yes | Comma-separated public channel usernames to watch (without `@`) |
| `KEYWORDS` | yes | Comma-separated keywords, case-insensitive |
| `POLL_INTERVAL_SECONDS` | no | How often to check the channels, in seconds (default: `10`) |
| `SESSION_FILE` | no | Session file path (default: `session.json`) |

Example `.env`:

```
API_ID=1234567
API_HASH=your_api_hash_here
PHONE=+380XXXXXXXXX
BOT_TOKEN=123456:ABC-DEF1234567890
MY_CHAT_IDS=123456789,987654321
CHANNELS=AlertChannel1,another_alert_channel
KEYWORDS=Київська,Київ
POLL_INTERVAL_SECONDS=10
SESSION_FILE=session.json
```

`MY_CHAT_IDS`, `CHANNELS`, and `KEYWORDS` each accept one or more
comma-separated values.

## Run

```
go mod tidy
go run ./cmd/bot
```

On the first run, Telegram sends a login code — enter it once. The session is
saved to `session.json` and reused on later runs, so you won't be asked again.

Before running, press **Start** on your bot in Telegram once, otherwise it
cannot send you messages.

## Test

```
go test ./...
```

## Running 24/7

The bot only runs while its machine is on, so for continuous use deploy it to an
always-on server (e.g. a small Linux VPS). Build a binary with
`GOOS=linux GOARCH=amd64 go build -o alertbot ./cmd/bot`, copy it to the server
together with `session.json` (generated locally first), set the environment
variables there, and run it under a process manager such as `systemd` so it
restarts automatically.