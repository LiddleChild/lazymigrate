package log

import (
	"encoding/json"
	"io"
	"log/slog"
	"time"
)

type LogDispatcher struct {
	writer  io.Writer
	decoder *json.Decoder
	ch      chan Message
}

func NewLogDispatcher() *LogDispatcher {
	var (
		pr, pw  = io.Pipe()
		decoder = json.NewDecoder(pr)
		ch      = make(chan Message, 128)
	)

	handler := &LogDispatcher{
		writer:  pw,
		decoder: decoder,
		ch:      ch,
	}

	return handler
}

func (d *LogDispatcher) Handle(level slog.Level) slog.Handler {
	opts := slog.HandlerOptions{
		Level: level,
	}

	go d.handle()

	return slog.NewJSONHandler(d.writer, &opts)
}

func (d *LogDispatcher) Pull() <-chan Message {
	return d.ch
}

func (d *LogDispatcher) handle() {
	for {
		if !d.decoder.More() {
			continue
		}

		var msg struct {
			LogAttribute

			Time  time.Time `json:"time"`
			Level LogLevel  `json:"level"`
			Msg   string    `json:"msg"`
		}

		_ = d.decoder.Decode(&msg)

		d.ch <- Message{
			Time:      msg.Time,
			Level:     msg.Level,
			Message:   msg.Msg,
			Secondary: msg.Secondary,
		}
	}
}
