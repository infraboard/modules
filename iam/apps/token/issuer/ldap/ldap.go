package ldap

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

// OWASP recommends to escape some special characters.
// https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/LDAP_Injection_Prevention_Cheat_Sheet.md
const specialLDAPRunes = ",#+<>;\"="

// UserProvider LDAP provider
type UserProvider interface {
	CheckConnect() error
	CheckUserPassword(username string, password string) (bool, error)
	GetDetails(username string) (*UserProfile, error)
	UpdatePassword(username string, newPassword string) error
}

func NewLdapProvider(conf Config) *LdapProvider {
	return &LdapProvider{
		conf: conf,
		log:  log.Sub("ldap"),
	}
}

type LdapProvider struct {
	conf Config
	log  *zerolog.Logger
}

// UserProfile todo
type UserProfile struct {
	DN          string
	Emails      []string
	Username    string
	DisplayName string
	Groups      []string
}

func (p *LdapProvider) dialTLS(network, addr string, config *tls.Config) (Connection, error) {
	conn, err := ldap.DialTLS(network, addr, config)
	if err != nil {
		return nil, err
	}

	return NewLDAPConnectionImpl(conn), nil
}

func (p *LdapProvider) dial(network, addr string) (Connection, error) {
	conn, err := ldap.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return NewLDAPConnectionImpl(conn), nil
}

func (p *LdapProvider) connect(userDN string, password string) (Connection, error) {
	var conn Connection

	url, err := url.Parse(p.conf.Url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL to LDAP: %s", url)
	}

	if url.Scheme == "ldaps" {
		p.log.Debug().Msg("LDAP client starts a TLS session")
		tlsConn, err := p.dialTLS("tcp", url.Host, &tls.Config{
			InsecureSkipVerify: p.conf.SkipVerify,
		})
		if err != nil {
			return nil, err
		}

		conn = tlsConn
	} else {
		p.log.Debug().Msg("LDAP client starts a session over raw TCP")
		rawConn, err := p.dial("tcp", url.Host)
		if err != nil {
			return nil, err
		}
		conn = rawConn
	}

	if err := conn.Bind(userDN, password); err != nil {
		return nil, err
	}

	return conn, nil
}

// CheckConnect todo
func (p *LdapProvider) CheckConnect() error {
	adminClient, err := p.connect(p.conf.BindDn, p.conf.BindPassword)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	return nil
}

// CheckUserPassword checks if provided password matches for the given user.
func (p *LdapProvider) CheckUserPassword(inputUsername string, password string) (*UserProfile, error) {
	adminClient, err := p.connect(p.conf.BindDn, p.conf.BindPassword)
	if err != nil {
		return nil, err
	}
	defer adminClient.Close()

	profile, err := p.getUserProfile(adminClient, inputUsername)
	if err != nil {
		return nil, err
	}

	conn, err := p.connect(profile.DN, password)
	if err != nil {
		return nil, fmt.Errorf("authentication of user %s failed. Cause: %s", inputUsername, err)
	}
	defer conn.Close()

	return profile, nil
}

func (p *LdapProvider) ldapEscape(inputUsername string) string {
	inputUsername = ldap.EscapeFilter(inputUsername)
	for _, c := range specialLDAPRunes {
		inputUsername = strings.ReplaceAll(inputUsername, string(c), fmt.Sprintf("\\%c", c))
	}

	return inputUsername
}

func (p *LdapProvider) resolveUserFilter(userFilter string, inputUsername string) string {
	inputUsername = p.ldapEscape(inputUsername)

	// We temporarily keep placeholder {0} for backward compatibility.
	userFilter = strings.ReplaceAll(userFilter, "{0}", inputUsername)

	// The {username} placeholder is equivalent to {0}, it's the new way, a named placeholder.
	userFilter = strings.ReplaceAll(userFilter, "{input}", inputUsername)

	// {username_attribute} and {mail_attribute} are replaced by the content of the attribute defined
	// in configuration.
	userFilter = strings.ReplaceAll(userFilter, "{username_attribute}", p.conf.UserNameAttribute)
	userFilter = strings.ReplaceAll(userFilter, "{mail_attribute}", p.conf.MailAttribute)
	return userFilter
}

func (p *LdapProvider) getUserProfile(conn Connection, inputUsername string) (*UserProfile, error) {
	userFilter := p.resolveUserFilter(p.conf.UserFilter, inputUsername)
	p.log.Debug().Msgf("Computed user filter is %s", userFilter)

	baseDN := p.conf.BaseDn

	attributes := []string{"dn",
		p.conf.MailAttribute,
		p.conf.UserNameAttribute,
		p.conf.DisplayNameAttribute,
	}

	// Search for the given username.
	searchRequest := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		1, 0, false, userFilter, attributes, nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot find user DN of user %s. Cause: %s", inputUsername, err)
	}

	if len(sr.Entries) == 0 {
		return nil, exception.NewNotFound("user not found")
	}

	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("multiple users %s found", inputUsername)
	}

	userProfile := UserProfile{
		DN: sr.Entries[0].DN,
	}

	for _, attr := range sr.Entries[0].Attributes {
		if attr.Name == p.conf.MailAttribute {
			userProfile.Emails = attr.Values
		}

		if attr.Name == p.conf.UserNameAttribute {
			if len(attr.Values) != 1 {
				return nil, fmt.Errorf("user %s cannot have multiple value for attribute %s",
					inputUsername, p.conf.UserNameAttribute)
			}

			userProfile.Username = attr.Values[0]
		}
		if attr.Name == p.conf.DisplayNameAttribute {
			userProfile.DisplayName = attr.Values[0]
		}
	}

	if userProfile.DN == "" {
		return nil, fmt.Errorf("no DN has been found for user %s", inputUsername)
	}

	return &userProfile, nil
}

func (p *LdapProvider) resolveGroupFilter(inputUsername string, profile *UserProfile) (string, error) { //nolint:unparam
	inputUsername = p.ldapEscape(inputUsername)

	// We temporarily keep placeholder {0} for backward compatibility.
	groupFilter := strings.ReplaceAll(p.conf.GroupFilter, "{0}", inputUsername)
	groupFilter = strings.ReplaceAll(groupFilter, "{input}", inputUsername)

	if profile != nil {
		// We temporarily keep placeholder {1} for backward compatibility.
		groupFilter = strings.ReplaceAll(groupFilter, "{1}", ldap.EscapeFilter(profile.Username))
		groupFilter = strings.ReplaceAll(groupFilter, "{username}", ldap.EscapeFilter(profile.Username))
		groupFilter = strings.ReplaceAll(groupFilter, "{dn}", ldap.EscapeFilter(profile.DN))
	}

	return groupFilter, nil
}

// GetDetails retrieve the groups a user belongs to.
func (p *LdapProvider) GetDetails(inputUsername string) (*UserProfile, error) {
	conn, err := p.connect(p.conf.BindDn, p.conf.BindPassword)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	profile, err := p.getUserProfile(conn, inputUsername)
	if err != nil {
		return nil, err
	}

	GroupFilter, err := p.resolveGroupFilter(inputUsername, profile)
	if err != nil {
		return nil, fmt.Errorf("unable to create group filter for user %s. Cause: %s", inputUsername, err)
	}

	p.log.Debug().Msgf("Computed groups filter is %s", GroupFilter)

	groupBaseDN := p.conf.BaseDn

	// Search for the given username.
	searchGroupRequest := ldap.NewSearchRequest(
		groupBaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, GroupFilter, []string{p.conf.GroupNameAttribute}, nil,
	)

	sr, err := conn.Search(searchGroupRequest)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve groups of user %s. Cause: %s", inputUsername, err)
	}

	for _, res := range sr.Entries {
		if len(res.Attributes) == 0 {
			p.log.Warn().Msgf("No groups retrieved from LDAP for user %s", inputUsername)
			break
		}
		// Append all values of the document. Normally there should be only one per document.
		profile.Groups = append(profile.Groups, res.Attributes[0].Values...)
	}

	return profile, nil
}

// UpdatePassword update the password of the given user.
func (p *LdapProvider) UpdatePassword(inputUsername string, newPassword string) error {
	client, err := p.connect(p.conf.BindDn, p.conf.BindPassword)

	if err != nil {
		return fmt.Errorf("unable to update password. Cause: %s", err)
	}

	profile, err := p.getUserProfile(client, inputUsername)

	if err != nil {
		return fmt.Errorf("unable to update password. Cause: %s", err)
	}

	modifyRequest := ldap.NewModifyRequest(profile.DN, nil)

	modifyRequest.Replace("userPassword", []string{newPassword})

	err = client.Modify(modifyRequest)

	if err != nil {
		return fmt.Errorf("unable to update password. Cause: %s", err)
	}

	return nil
}
