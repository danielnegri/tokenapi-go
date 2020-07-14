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

package sync

import (
	"math"
	"sync"
)

// WaitGroup has the same role and close to the
// same API as the Golang sync.WaitGroup but adds a limit of
// the amount of goroutines started concurrently.
type WaitGroup struct {
	Size int

	current chan struct{}
	wg      sync.WaitGroup
}

// NewWaitGroup creates a WaitGroup.
// The limit parameter is the maximum amount of
// goroutines which can be started concurrently.
func NewWaitGroup(limit int) WaitGroup {
	size := math.MaxInt32 // 2^32 - 1
	if limit > 0 {
		size = limit
	}
	return WaitGroup{
		Size: size,

		current: make(chan struct{}, size),
		wg:      sync.WaitGroup{},
	}
}

// Add increments the internal WaitGroup counter.
// It can be blocking if the limit of spawned goroutines
// has been reached. It will stop blocking when Done is
// been called.
//
// See sync.WaitGroup documentation for more information.
func (s *WaitGroup) Add() {
	s.current <- struct{}{}
	s.wg.Add(1)
}

// Done decrements the WaitGroup counter.
// See sync.WaitGroup documentation for more information.
func (s *WaitGroup) Done() {
	<-s.current
	s.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
// See sync.WaitGroup documentation for more information.
func (s *WaitGroup) Wait() {
	s.wg.Wait()
}
