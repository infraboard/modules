package token

import "context"

var providers = map[string]Provider{}

func RegistryProvider(name string, p Provider) {
	providers[name] = p
}

func GetProvider(name string) Provider {
	return providers[name]
}

type Provider interface {
	IssueToken(context.Context, ProviderParameter) (*Token, error)
}
