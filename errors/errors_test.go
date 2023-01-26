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

//go:build !debug
// +build !debug

package errors

import (
	"os"
	"os/exec"
	"testing"

	"github.com/danielnegri/tokenapi-go/ledger"
)

func TestDebug(t *testing.T) {
	// Test with -tags debug to run the tests in debug_test.go
	cmd := exec.Command("go", "test", "-tags", "debug")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("external go test failed: %v", err)
	}
}

func TestMarshal(t *testing.T) {
	token := "xPGvwdBqDrpFLXyMVf0ovQ"

	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := E(Op("Get"), token, IO, "network unreachable")

	// Nested error.
	e2 := E(Op("Read"), token, Other, e1)

	b := MarshalError(e2)
	e3 := UnmarshalError(b)

	in := e2.(*Error)
	out := e3.(*Error)
	// Compare elementwise.
	if in.Token != out.Token {
		t.Errorf("expected Token %q; got %q", in.Token, out.Token)
	}
	if in.Op != out.Op {
		t.Errorf("expected Op %q; got %q", in.Op, out.Op)
	}
	if in.Kind != out.Kind {
		t.Errorf("expected kind %d; got %d", in.Kind, out.Kind)
	}
	// Note that error will have lost type information, so just check its Error string.
	if in.Err.Error() != out.Err.Error() {
		t.Errorf("expected Err %q; got %q", in.Err, out.Err)
	}
}

func TestSeparator(t *testing.T) {
	defer func(prev string) {
		Separator = prev
	}(Separator)
	Separator = ":: "

	// Same pattern as above.
	token := ledger.Token("3oMUY0bSsieok9GKuSQKpQ")

	// Single error. No token is set, so we will have a zero-length field inside.
	e1 := E(Op("Get"), IO, "network unreachable")

	// Nested error.
	e2 := E(Op("Write"), token, Other, e1)

	want := "Write: 3oMUY0bSsieok9GKuSQKpQ: I/O error:: Get: network unreachable"
	got := errorAsString(e2)
	if got != want {
		t.Errorf("expected %q; got %q", want, e2)
	}
}

// errorAsString returns the string form of the provided error value.
// If the given string is an *Error, the stack information is removed
// before the value is stringified.
func errorAsString(err error) string {
	if e, ok := err.(*Error); ok {
		e2 := *e
		e2.stack = stack{}
		return e2.Error()
	}
	return err.Error()
}
