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

package source

import (
	"context"
	"testing"
	"time"

	"github.com/danielnegri/adheretech/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	src := New(nil)
	assert.NotNil(t, src)
	assert.NotNil(t, src.httpClient)
	assert.Equal(t, src.httpClient.HostURL, DefaultURL)
	assert.Equal(t, src.httpClient.RetryCount, DefaultRetry)
	assert.Equal(t, src.httpClient.GetClient().Timeout, DefaultTimeout)
}

func Test_client_Check(t *testing.T) {
	ts := newTestServer(t, time.Now().Unix())
	defer ts.Close()

	assert.NoError(t, newTestSource(t, ts.URL).Check(context.Background()))

	err := newTestSource(t, discardURL).Check(context.Background())
	assert.Error(t, err)
	assert.True(t, errors.Is(errors.Internal, err))
}
