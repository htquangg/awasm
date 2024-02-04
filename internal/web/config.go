package web

type Config struct {
	ShowStartBanner bool `json:"show_start_banner" mapstructure:"show_start_banner" yaml:"show_start_banner"`

	Addr string `json:"addr" mapstructure:"addr" yaml:"addr"`
}
