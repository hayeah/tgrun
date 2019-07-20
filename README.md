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

# Behaviour

Suppose that we are running the [`now`](now/main.go) program, which prints the current time every second.

```
# go build -o now.bin now/main.go
tgrun ./now.bin
```

The tail of stdout and stderr will continuously update in a dedicated message bubble:

```
[mbp15.local pid=9017] `./now.bin`
Uptime: Uptime: 18.181806172s

5:56.938398 +0800 HKT m=+9.021042450
2019-07-20 20:05:57.938569 +0800 HKT m=+10.021222764
2019-07-20 20:05:58.9437 +0800 HKT m=+11.026363178
2019-07-20 20:05:59.945458 +0800 HKT m=+12.028130709
2019-07-20 20:06:00.947535 +0800 HKT m=+13.030216833
2019-07-20 20:06:01.947877 +0800 HKT m=+14.030567394
2019-07-20 20:06:02.949774 +0800 HKT m=+15.032473857
2019-07-20 20:06:03.950218 +0800 HKT m=+16.032927214
2019-07-20 20:06:04.955319 +0800 HKT m=+17.038037926
2019-07-20 20:06:05.957598 +0800 HKT m=+18.040326255
```

On exit, a notification will be sent:

```
[mbp15.local pid=9017] `./now.bin`
Exit status: -1
```
