package token

import (
	"context"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

const (
	ISSUER_LDAP          = "ldap"
	ISSUER_PASSWORD      = "password"
	ISSUER_PRIVATE_TOKEN = "private_token"
)

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

// MakeBearer https://tools.ietf.org/html/rfc6750#section-2.1
// b64token    = 1*( ALPHA / DIGIT /"-" / "." / "_" / "~" / "+" / "/" ) *"="
func MakeBearer(lenth int) string {
	charlist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	t := make([]string, lenth)
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano() + int64(lenth) + rand.Int63n(10000))))
	for i := 0; i < lenth; i++ {
		rn := r.Intn(len(charlist))
		w := charlist[rn : rn+1]
		t = append(t, w)
	}

	token := strings.Join(t, "")
	return token
}
