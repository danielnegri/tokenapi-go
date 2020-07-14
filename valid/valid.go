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
