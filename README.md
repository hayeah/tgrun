Notify you in Telegram chat when a command line program is done running.

Install:

```
go get github.com/hayeah/tgrun
```

Configure the bot credentials:

```
export TG_BOT_TOKEN=
export TG_CHAT_ID=
```

(See: [how to get bot token and chat id](telegram-api-send-message-personal-notification-bot/)).

Then run your (possibly long-running)

```
tgrun myprogram arg1 arg2 arg3
```

On exit, a notification will be sent:

```
[mbp15.local] `myprogram` exited with status 0
Elapsed: 5.537169ms
```
