package token

import (
	"context"
	"math/rand/v2"
)

const (
	ISSUER_LDAP          = "ldap"
	ISSUER_FEISHU        = "feishu"
	ISSUER_PASSWORD      = "password"
	ISSUER_PRIVATE_TOKEN = "private_token"
)

var issuers = map[string]Issuer{}

func RegistryIssuer(name string, p Issuer) {
	issuers[name] = p
}

func GetIssuer(name string) Issuer {
	return issuers[name]
}

type Issuer interface {
	IssueToken(context.Context, IssueParameter) (*Token, error)
}

var (
	charlist = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// MakeBearer https://tools.ietf.org/html/rfc6750#section-2.1
// b64token    = 1*( ALPHA / DIGIT /"-" / "." / "_" / "~" / "+" / "/" ) *"="
func MakeBearer(lenth int) string {
	t := make([]byte, 0)
	for range lenth {
		rn := rand.IntN(len(charlist))
		t = append(t, charlist[rn])
	}
	return string(t)
}
