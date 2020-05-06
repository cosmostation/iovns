package simulation

import (
	"fmt"
	"testing"

	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeDistributionStore(t *testing.T) {
	cdc := makeTestCodec()

	domainName := "sucuk"
	domainAdmin := sdk.AccAddress(crypto.AddressHash([]byte("domainAadmin")))
	broker := sdk.AccAddress(crypto.AddressHash([]byte("broker")))
	domain := types.Domain{
		Name:         domainName,
		Admin:        domainAdmin,
		ValidUntil:   0,
		HasSuperuser: false,
		AccountRenew: 5,
		Broker:       broker,
	}

	accountName := "doner"
	accountOwner := sdk.AccAddress(crypto.AddressHash([]byte("accountOwner")))
	account := types.Account{
		Domain:       domainName,
		Name:         accountName,
		Owner:        accountOwner,
		ValidUntil:   0,
		Targets:      nil,
		Certificates: nil,
		Broker:       broker,
	}
	kvPairs := kv.Pairs{
		kv.Pair{Key: keeper.AccountByOwnerPrefix, Value: cdc.MustMarshalBinaryLengthPrefixed(account)},
		kv.Pair{Key: keeper.DomainByOwnerPrefix, Value: cdc.MustMarshalBinaryLengthPrefixed(domain)},
		kv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Account", fmt.Sprintf("%v\n%v", account, account)},
		{"Domain", fmt.Sprintf("%v\n%v", domain, domain)},
		{"other", ""},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
