package rmqc

import "log/slog"

type Options struct {
	AutoAck   bool
	Exclusive bool
	Capacity  int
	NoLocal   bool
	NoWait    bool
	Consumer  string
	Args      map[string]interface{}
	Logger    *slog.Logger
	Debug     bool
}

func handleOptions(options *Options) *Options {
	if options == nil {
		options = &Options{}
	}
	if options.Capacity <= 0 {
		options.Capacity = 10
	}
	if options.Logger == nil {
		options.Logger = slog.Default()
	}
	return options
}
