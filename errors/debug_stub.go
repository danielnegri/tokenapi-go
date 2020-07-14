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

// +build !debug

package errors

import "bytes"

// These are stubs that disable stack collecting/printing
// when the debug build tag is not set.
// For documentation about these types and functions, see debug.go.

type stack struct{}

func (e *Error) populateStack()           {}
func (e *Error) printStack(*bytes.Buffer) {}
