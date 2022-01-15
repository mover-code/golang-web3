package jsonrpc

import (
	"testing"

	"github.com/mover-code/golang-web3/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNetVersion(t *testing.T) {
	testutil.MultiAddr(t, nil, func(s *testutil.TestServer, addr string) {
		c, _ := NewClient(addr)
		defer c.Close()

		_, err := c.Net().Version()
		assert.NoError(t, err)
	})
}

func TestNetListening(t *testing.T) {
	testutil.MultiAddr(t, nil, func(s *testutil.TestServer, addr string) {
		c, _ := NewClient(addr)
		defer c.Close()

		ok, err := c.Net().Listening()
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}

func TestNetPeerCount(t *testing.T) {
	testutil.MultiAddr(t, nil, func(s *testutil.TestServer, addr string) {
		c, _ := NewClient(addr)
		defer c.Close()

		count, err := c.Net().PeerCount()
		assert.NoError(t, err)
		assert.Equal(t, count, uint64(0))
	})
}
