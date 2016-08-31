package ldap

import "testing"

func TestLDAPAuthProvider_Authenticate(t *testing.T) {
	// This is publicly available test ldap server available here
	// more info available here http://www.forumsys.com/en/tutorials/integration-how-to/ldap/online-ldap-test-server/
	ldapAuthProvider := LDAPAuthProvider{
		Host:         "ldap.forumsys.com",
		Port:         389,
		BindDN:       "cn=read-only-admin,dc=example,dc=com",
		BindPassword: "password",
		UserFilter:   "(uid=%s)",
		BaseDN:       "dc=example,dc=com",
	}

	ldapAuthProvider.Connect()
	defer ldapAuthProvider.Close()

	authenticated, _, err := ldapAuthProvider.Authenticate("gauss", "password")
	if err != nil || !authenticated {
		t.Errorf("LDAP auth provider failed")
	}
}
