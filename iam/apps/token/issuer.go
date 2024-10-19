package token

import "context"

var issuer = map[string]Issuer{}

func RegistryIssuer(name string, p Issuer) {
	issuer[name] = p
}

func GetIssue(name string) Issuer {
	return issuer[name]
}

type Issuer interface {
	IssueToken(context.Context, IssueParameter) (*Token, error)
}
