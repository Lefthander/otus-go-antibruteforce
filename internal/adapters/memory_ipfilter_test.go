package adapters

import (
	"context"
	"net"
	"testing"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestIPFilterMemory(t *testing.T) {

	ctx := context.Background()

	ipstore := NewIPFilterMemory()

	_, testnet1, _ := net.ParseCIDR("10.10.0.0/24")

	_, testnet2, _ := net.ParseCIDR("20.20.0.0/24")

	_, testnet3, _ := net.ParseCIDR("30.30.0.0/24")

	t.Run("Add network to White list", func(t *testing.T) {
		ipstore.AddIPNetwork(ctx, *testnet1, true)
		ipl, err := ipstore.ListIPNetworks(ctx, true)
		assert.Equal(t, nil, err)
		assert.Contains(t, ipl, testnet1.String())
	})

	t.Run("Check white list for defined network", func(t *testing.T) {
		flag, err := ipstore.IsIPConform(ctx, net.ParseIP("10.10.0.1"))
		assert.Equal(t, true, flag)
		assert.Equal(t, errors.ErrIPFilterMatchedWhiteList, err)

	})

	t.Run("Add network to Black list", func(t *testing.T) {
		ipstore.AddIPNetwork(ctx, *testnet2, false)
		ipl, err := ipstore.ListIPNetworks(ctx, false)
		assert.Equal(t, nil, err)
		assert.Contains(t, ipl, testnet2.String())

	})

	t.Run("Check black list for defined network", func(t *testing.T) {
		flag, err := ipstore.IsIPConform(ctx, net.ParseIP("20.20.0.1"))
		assert.Equal(t, true, flag)
		assert.Equal(t, errors.ErrIPFilterMatchedBlackList, err)

	})
	t.Run("Check IP list from the undefined network", func(t *testing.T) {
		flag, err := ipstore.IsIPConform(ctx, net.ParseIP("30.30.0.1"))
		assert.Equal(t, false, flag)
		assert.Equal(t, errors.ErrIPFilterNoMatch, err)

	})

	t.Run("Delete defined network from White list", func(t *testing.T) {
		err := ipstore.DeleteIPNetwork(ctx, *testnet1, true)
		assert.Equal(t, nil, err)
		assert.NotContains(t, ipstore.WhiteIPList.Nets, "10.10.0.0/24")
	})

	t.Run("Delete defined network from Black list", func(t *testing.T) {
		err := ipstore.DeleteIPNetwork(ctx, *testnet2, false)
		assert.Equal(t, nil, err)
		assert.NotContains(t, ipstore.WhiteIPList.Nets, "20.20.0.0/24")
	})

	t.Run("Delete undefined network from White list", func(t *testing.T) {
		err := ipstore.DeleteIPNetwork(ctx, *testnet3, true)
		assert.Equal(t, errors.ErrIPFilterNetworkNotFound, err)
	})

	t.Run("Delete undefined network from Black list", func(t *testing.T) {
		err := ipstore.DeleteIPNetwork(ctx, *testnet3, false)
		assert.Equal(t, errors.ErrIPFilterNetworkNotFound, err)
	})

}
