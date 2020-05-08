package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/tendermint/tendermint/libs/kv"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding auction type
func DecodeStore(cdc *codec.Codec, kvA, kvB kv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.AccountStorePrefix):
		var accA, accB types.Account
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &accA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &accB)
		return fmt.Sprintf("%v\n%v", accA, accB)
	case bytes.Equal(kvA.Key[:1], types.DomainStorePrefix):
		var dA, dB types.Domain
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &dA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &dB)
		return fmt.Sprintf("%v\n%v", dA, dB)
	case bytes.Equal(kvA.Key[:1], types.OwnerToAccountPrefix):
		var dA, dB []byte
		cdc.MustUnmarshalBinaryBare(kvA.Value, &dA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &dB)
		return fmt.Sprintf("%v\n%v", dA, dB)
	case bytes.Equal(kvA.Key[:1], types.OwnerToDomainPrefix):
		var dA, dB []byte
		cdc.MustUnmarshalBinaryBare(kvA.Value, &dA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &dB)
		return fmt.Sprintf("%v\n%v", dA, dB)
	default:
		panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
	}
}
