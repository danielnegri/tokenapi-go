// Copyright 2020 The Ledger Authors
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
