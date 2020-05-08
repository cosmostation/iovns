package simulation

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/kv"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeDomainStore(t *testing.T) {
	cdc := makeTestCodec()

	domain := types.Domain{
		Name:         "",
		Admin:        nil,
		ValidUntil:   0,
		HasSuperuser: false,
		AccountRenew: 0,
		Broker:       nil,
	}
	account := types.Account{
		Domain:       "",
		Name:         "",
		Owner:        nil,
		ValidUntil:   0,
		Targets:      nil,
		Certificates: nil,
		Broker:       nil,
		MetadataURI:  "",
	}
	kvPairs := kv.Pairs{
		kv.Pair{Key: types.DomainStorePrefix, Value: cdc.MustMarshalBinaryLengthPrefixed(&domain)},
		kv.Pair{Key: types.AccountStorePrefix, Value: cdc.MustMarshalBinaryLengthPrefixed(&account)},
		kv.Pair{Key: types.OwnerToDomainPrefix, Value: cdc.MustMarshalBinaryBare([]byte{})},
		kv.Pair{Key: types.OwnerToAccountPrefix, Value: cdc.MustMarshalBinaryBare([]byte{})},
		kv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Domain", fmt.Sprintf("%v\n%v", domain, domain)},
		{"Account", fmt.Sprintf("%v\n%v", account, account)},
		{"OwnerToDomain", "[]\n[]"},
		{"OwnerToAccount", "[]\n[]"},
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
