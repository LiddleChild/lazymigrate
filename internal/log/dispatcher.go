package log

import (
	"encoding/json"
	"io"
	"log/slog"
	"time"
)

type LogDispatcher struct {
	handler *slog.JSONHandler
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
		handler: slog.NewJSONHandler(pw, nil),
		decoder: decoder,
		ch:      ch,
	}

	return handler
}

func (d *LogDispatcher) Handler() slog.Handler {
	go d.handle()

	return d.handler
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
			Time  time.Time `json:"time"`
			Level LogLevel  `json:"level"`
			Msg   string    `json:"msg"`
			Attrs struct {
				Secondary bool `json:"secondary"`
			} `json:"attributes"`
		}

		_ = d.decoder.Decode(&msg)

		d.ch <- Message{
			Time:      msg.Time,
			Level:     msg.Level,
			Message:   msg.Msg,
			Secondary: msg.Attrs.Secondary,
		}
	}
}
