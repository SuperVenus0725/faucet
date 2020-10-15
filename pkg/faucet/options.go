package faucet

var defaultOptions = Options{
	AppCli:       "gaiacli",
	KeyName:      "faucet",
	Denom:        "atom",
	CreditAmount: 10000000,
	MaxCredit:    100000000,
}

type Options struct {
	AppCli          string
	KeyringPassword string
	KeyName         string
	KeyMnemonic     string
	Denom           string
	CreditAmount    uint64
	MaxCredit       uint64
}

type Option func(*Options)

func CliName(s string) Option {
	return func(opts *Options) {
		opts.AppCli = s
	}
}

func KeyringPassword(s string) Option {
	return func(opts *Options) {
		opts.KeyringPassword = s
	}
}

func KeyName(s string) Option {
	return func(opts *Options) {
		opts.KeyName = s
	}
}

func WithMnemonic(s string) Option {
	return func(opts *Options) {
		opts.KeyMnemonic = s
	}
}

func Denom(s string) Option {
	return func(opts *Options) {
		opts.Denom = s
	}
}

func CreditAmount(v uint64) Option {
	return func(opts *Options) {
		opts.CreditAmount = v
	}
}

func MaxCredit(v uint64) Option {
	return func(opts *Options) {
		opts.MaxCredit = v
	}
}
