package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/armon/circbuf"
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

	args := os.Args[1:]

	if len(args) == 0 {
		return errors.New("must specify a command to run")
	}

	host, err := os.Hostname()
	if err != nil {
		return err
	}

	cmd := args[0]
	args = args[1:]

	bot, err := tgbotapi.NewBotAPI(token)

	r := runner{
		Host:   host,
		ChatID: int64(chatID),
		Bot:    bot,
		Cmd:    exec.Command(cmd, args...),
	}

	err = r.start()
	return err
}

// syncBuf is a threadsafe wrapper for *circbuf.Buffer
type syncBuf struct {
	buf *circbuf.Buffer
	sync.Mutex
}

func newSyncBuf(size int64) (*syncBuf, error) {
	buf, err := circbuf.NewBuffer(size)
	if err != nil {
		return nil, err
	}

	return &syncBuf{
		buf: buf,
	}, nil
}

func (b *syncBuf) Write(buf []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	return b.buf.Write(buf)
}

func (b *syncBuf) Bytes() []byte {
	b.Lock()
	defer b.Unlock()
	data := b.buf.Bytes()
	buf := make([]byte, len(data))
	copy(buf, data)
	return buf
}

type runner struct {
	Host   string
	ChatID int64
	Bot    *tgbotapi.BotAPI
	Cmd    *exec.Cmd

	buf *syncBuf
}

func (r *runner) start() error {
	// FIXME: make me threadsafe...
	buf, err := newSyncBuf(512)
	if err != nil {
		return err
	}
	r.buf = buf

	go func() {
		err := r.updateStatus()
		if err != nil {
			log.Println(err)
		}
	}()

	go r.handleInterrupt()

	return r.runCommand()
}

func (r *runner) runCommand() error {
	c := r.Cmd

	c.Stdin = os.Stdin

	c.Stdout = io.MultiWriter(os.Stdout, r.buf)
	c.Stderr = io.MultiWriter(os.Stderr, r.buf)

	err := c.Start()
	if err != nil {
		return err
	}

	c.Wait()

	_, err = r.sendMessage("Exit status: %d", c.ProcessState.ExitCode())
	return err
}

func (r *runner) handleInterrupt() {
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)

	for sig := range sigc {
		log.Println("SIGINT, exit now", sig)
	}
}

func (r *runner) Tag() string {
	return fmt.Sprintf("[%s] %d `%s`", r.Host, r.Cmd.Process.Pid, r.Cmd.Args[0])
}

func (r *runner) sendMessage(format string, a ...interface{}) (tgbotapi.Message, error) {
	s := fmt.Sprintf("%s\n%s", r.Tag(), fmt.Sprintf(format, a...))
	return r.Bot.Send(tgbotapi.NewMessage(r.ChatID, s))
}

func (r *runner) editMessage(msgid int, format string, a ...interface{}) (tgbotapi.Message, error) {
	s := fmt.Sprintf("%s\n%s", r.Tag(), fmt.Sprintf(format, a...))
	return r.Bot.Send(tgbotapi.NewEditMessageText(r.ChatID, msgid, s))
}

func (r *runner) updateStatus() error {
	// This loop terminates when the program terminates...
	start := time.Now()

	var m tgbotapi.Message
	for {
		var err error

		time.Sleep(2 * time.Second)
		elapsed := time.Now().Sub(start)

		txt := fmt.Sprintf("Uptime: %s", elapsed)

		tail := r.buf.Bytes()
		if len(tail) > 0 {
			txt = fmt.Sprintf("%s\n\n%s", txt, string(tail))
		}

		// create or update status text

		if m.MessageID == 0 {
			m, err = r.sendMessage("Uptime: %s", txt)

			if err != nil {
				return err
			}

			continue
		}

		_, err = r.editMessage(m.MessageID, "Uptime: %s", txt)

		if err != nil {
			return err
		}
	}
}
