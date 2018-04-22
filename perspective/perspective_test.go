package perspective

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	p := New()
	require.NoError(t, p.Update("obj1", 1, nil))
	require.NoError(t, p.Update("obj1", 2, nil))
	require.Error(t, p.Update("obj1", 0, nil))
	require.Error(t, p.Update("obj1", -1, nil))

	v, e := p.Version("obj1")
	require.NoError(t, e)
	require.Equal(t, 2, v)

	v, e = p.Version("objN")
	require.Error(t, e)
}
