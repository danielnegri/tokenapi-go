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
