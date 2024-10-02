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

func defaultOptions() *Options {
	return &Options{
		Capacity: 10,
		Args:     nil,
		Logger:   slog.Default(),
	}
}

func handleOptions(options *Options) *Options {
	if options == nil {
		return defaultOptions()
	}
	if options.Capacity <= 0 {
		options.Capacity = defaultOptions().Capacity
	}
	return options
}
