package perspective

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	p := New()
	require.NoError(t, p.Update("obj1", 1))
	require.NoError(t, p.Update("obj1", 2))
	require.Error(t, p.Update("obj1", 0))
	require.Error(t, p.Update("obj1", -1))

	v, e := p.Version("obj1")
	require.NoError(t, e)
	require.Equal(t, 2, v)

	v, e = p.Version("objN")
	require.Error(t, e)
}
