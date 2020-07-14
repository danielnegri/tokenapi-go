// Package valid does validation of various data types.
// For the most part, its functions are used by servers and enforce
// stronger constraints than client code needs to follow.
package valid

import (
	"strings"

	"github.com/danielnegri/adheretech/errors"
	"github.com/danielnegri/adheretech/ledger"
)

// Token verifies that the token is a valid non-empty string.
// It also requires that the string does not start with a "-".
func Token(token ledger.Token) error {
	const op errors.Op = "valid.Token"

	if len(token) == 0 {
		return errors.E(op, errors.Invalid, token)
	}

	if strings.Contains(string(token), "-") {
		return errors.E(op, errors.Invalid, token, "cannot contain dash")
	}

	return nil
}
