package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	ctx context.Context
}

var testContext testStruct

func TestHandler(t *testing.T) {
	ok, err := HandleRequest(testContext.ctx)

	assert.IsType(t, nil, err)
	assert.Equal(t, "ok", ok)
}
