package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandIdString(t *testing.T) {
	// arrange
	utils := New()

	// act
	result := utils.GetRandIdString(10)

	//accert
	assert.Len(t, result, 10)
	assert.Regexp(t, "[0-9a-zA-Z]", result)

}
