package humanize_test

import (
	"testing"

	"github.com/jmhobbs/srv/internal/humanize"
	"github.com/stretchr/testify/assert"
)

func Test_Bytes(t *testing.T) {
	assert.Equal(t, "0 b", humanize.Bytes(0))
	assert.Equal(t, "420 b", humanize.Bytes(420))
	assert.Equal(t, "2.01 kb", humanize.Bytes(2056))
	assert.Equal(t, "1.00 mb", humanize.Bytes(1_048_576))
	assert.Equal(t, "1.00 gb", humanize.Bytes(1_073_741_824))
}
