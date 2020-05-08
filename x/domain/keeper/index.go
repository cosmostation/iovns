package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
)

// indexStore returns the indexing store from the module's kvstore
func indexStore(store sdk.KVStore) sdk.KVStore {
	return prefix.NewStore(store, types.IndexStorePrefix)
}

// domainIndexStore returns the kvstore space that maps
// owner to domain keys
func domainIndexStore(store sdk.KVStore) sdk.KVStore {
	return prefix.NewStore(store, types.OwnerToDomainPrefix)
}

// accountIndexStore returns the kvstore space that maps
// owner to accounts
func accountIndexStore(store sdk.KVStore) sdk.KVStore {
	return prefix.NewStore(store, types.OwnerToAccountPrefix)
}

func (k Keeper) unmapAccountToOwner(ctx sdk.Context, account types.Account) {
	// get store
	store := accountIndexStore(indexStore(ctx.KVStore(k.storeKey)))

	// check if key exists TODO remove panic
	key := types.GetOwnerToAccountKey(account.Owner, account.Domain, account.Name)
	if !store.Has(key) {
		panic(fmt.Sprintf("missing store key: %s", key))
	}
	// delete key
	store.Delete(key)
}

// mapAccountToOwner maps accounts to an owner
func (k Keeper) mapAccountToOwner(ctx sdk.Context, account types.Account) {
	// get store
	store := accountIndexStore(indexStore(ctx.KVStore(k.storeKey)))
	key := types.GetOwnerToAccountKey(account.Owner, account.Domain, account.Name)
	// check if key exists TODO remove panic
	if store.Has(key) {
		panic(fmt.Sprintf("existing store key: %s", key))
	}
	// set key
	store.Set(key, []byte{})
}

func (k Keeper) iterAccountToOwner(ctx sdk.Context, address sdk.AccAddress, do func(key []byte) bool) {
	// get store
	store := accountIndexStore(indexStore(ctx.KVStore(k.storeKey)))
	// get iterator
	iterator := sdk.KVStorePrefixIterator(store, types.IndexAddr(address))
	defer iterator.Close()
	// iterate keys
	for ; iterator.Valid(); iterator.Next() {
		// do action
		keepGoing := do(iterator.Key())
		// keep going?
		if !keepGoing {
			return
		}
	}
}

func (k Keeper) mapDomainToOwner(ctx sdk.Context, domain types.Domain) {
	// get store
	store := domainIndexStore(indexStore(ctx.KVStore(k.storeKey)))
	// get unique key
	key := types.GetOwnerToDomainKey(domain.Admin, domain.Name)
	// check if key exists TODO remove panic
	if store.Has(key) {
		panic(fmt.Sprintf("existing store key: %s", key))
	}
	// set key
	store.Set(key, []byte{})
}

func (k Keeper) unmapDomainToOwner(ctx sdk.Context, domain types.Domain) {
	// get store
	store := domainIndexStore(indexStore(ctx.KVStore(k.storeKey)))
	// check if key exists TODO remove panic
	key := types.GetOwnerToDomainKey(domain.Admin, domain.Name)
	if !store.Has(key) {
		panic(fmt.Sprintf("missing store key: %s", key))
	}
	// delete key
	store.Delete(key)
}

// iterDomainToOwner iterates over all the domains owned by address
// and returns the unique keys
func (k Keeper) iterDomainToOwner(ctx sdk.Context, address sdk.AccAddress, do func(key []byte) bool) {
	// get store
	store := domainIndexStore(indexStore(ctx.KVStore(k.storeKey)))
	// get iterator
	iterator := sdk.KVStorePrefixIterator(store, types.IndexAddr(address))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if !do(iterator.Key()) {
			return
		}
	}
}
