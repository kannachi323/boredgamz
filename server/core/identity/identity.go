// Package identity holds helpers for socket player IDs.
//
// Authenticated sockets use "<accountID>:<tabID>" (plus ":spectator" for
// observers), so one account can hold several tabs without them being
// mistaken for each other on reconnect.
package identity

import "strings"

// Account returns the account part of a socket player ID.
func Account(playerID string) string {
	account, _, _ := strings.Cut(playerID, ":")
	return account
}

// SameAccount reports whether two socket player IDs belong to one account, so
// matchmaking never seats a player against their own other tab.
func SameAccount(a, b string) bool {
	return Account(a) == Account(b)
}
