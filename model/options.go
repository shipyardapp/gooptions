package model

type Options struct {
	OptionName string

	OptionPrefix string
}

func NewOptions() *Options {
	return &Options{
		OptionName:   "Option",
		OptionPrefix: "With",
	}
}
