package webtoken_test

import (
	"testing"
	"time"

	"github.com/josephspurrier/octane/example/app/lib/webtoken"
	"github.com/stretchr/testify/assert"
)

func TestClock(t *testing.T) {
	assert.True(t, time.Since(new(webtoken.Clock).Now()) < time.Second)
}
