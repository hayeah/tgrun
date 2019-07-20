package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	token, ok := os.LookupEnv("TG_BOT_TOKEN")
	if !ok {
		return errors.New("must set TG_BOT_TOKEN")
	}

	chatIDenv, ok := os.LookupEnv("TG_CHAT_ID")
	if !ok {
		return errors.New("must set TG_CHAT_ID")
	}

	chatID, err := strconv.Atoi(chatIDenv)
	if err != nil {
		return err
	}

	host, err := os.Hostname()
	if err != nil {
		return err
	}

	args := os.Args[1:]

	if len(args) == 0 {
		return errors.New("must specify a command to run")
	}

	cmd := args[0]
	args = args[1:]

	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin

	err = c.Start()
	if err != nil {
		return err
	}

	start := time.Now()

	c.Wait()

	elapsed := time.Now().Sub(start)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	msgtxt := fmt.Sprintf("[%s] `%s` exited with status %d\nElapsed: %s", host, cmd, c.ProcessState.ExitCode(), elapsed)

	msg := tgbotapi.NewMessage(int64(chatID), msgtxt)
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
