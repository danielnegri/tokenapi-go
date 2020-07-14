package valid

import (
	"testing"

	"github.com/danielnegri/adheretech/ledger"
)

func TestToken(t *testing.T) {
	tests := []struct {
		token ledger.Token
		valid bool
	}{
		{token: "", valid: false},
		{token: "xPGvwdBqDrpFLXyMVf0ovQ", valid: true},
		{token: "_-kFu9fparYLZtyNBDH9vg", valid: false},
		{token: "3oMUY0bSsieok9GKuSQKpQ", valid: true},
	}
	for _, test := range tests {
		err := Token(test.token)
		if test.valid == (err == nil) {
			continue
		}

		t.Errorf("%q: expected valid=%t; got error %v", test.token, test.valid, err)
	}
}
