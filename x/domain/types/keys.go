package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
)

var (
	// DomainStorePrefix is the prefix used to define the prefixed store containing domain data
	DomainStorePrefix = []byte{0x00}
	// AccountPrefixStore is the prefix used to define the prefixed store containing account data
	AccountStorePrefix = []byte{0x01}
	// IndexStorePrefix is the prefix used to defines the prefixed store containing indexing data
	IndexStorePrefix = []byte{0x02}

	// OwnerToAccountPrefix is the prefix that matches owners to accounts
	OwnerToAccountPrefix = []byte{0x04}
	// OwnerToAccountIndexSeparator is the separator used to map owner address + domain + account name
	OwnerToAccountIndexSeparator = []byte(":")
	// OwnerToDomainPrefix is the prefix that matches owners to domains
	OwnerToDomainPrefix = []byte{0x05}
	// OwnerToDomainIndexSeparator is the separator used to map owner address + domain
	OwnerToDomainIndexSeparator = []byte(":")
)

// GetDomainPrefixKey returns the domain prefix byte key
func GetDomainPrefixKey(domainName string) []byte {
	return []byte(domainName)
}

// GetAccountKey returns the account byte key by its name
func GetAccountKey(accountName string) []byte {
	return []byte(accountName)
}

// AccountKeyToString converts account key bytes to string
func AccountKeyToString(accountKeyBytes []byte) string {
	return string(accountKeyBytes)
}

// GetOwnerToAccountKey generates the unique key that maps owner to account
func GetOwnerToAccountKey(owner types.AccAddress, domain string, account string) []byte {
	// get index bytes of addr
	addr := IndexAddr(owner)
	// generate unique key
	return bytes.Join([][]byte{addr, []byte(domain), []byte(account)}, OwnerToAccountIndexSeparator)
}

func IndexAddr(addr types.AccAddress) []byte {
	x := addr.String()
	return []byte(x)
}

func AccAddrFromIndex(indexedAddr []byte) types.AccAddress {
	accAddr, err := types.AccAddressFromBech32(string(indexedAddr))
	if err != nil {
		panic(err)
	}
	return accAddr
}

// GetOwnerToDomainKey generates the unique key that maps owner to domain
func GetOwnerToDomainKey(owner types.AccAddress, domain string) []byte {
	addrBytes := IndexAddr(owner)
	return bytes.Join([][]byte{addrBytes, []byte(domain)}, OwnerToDomainIndexSeparator)
}

// SplitOwnerToAccountKey takes an indexed owner to account key and splits it
// into owner address, domain name and account name
func SplitOwnerToAccountKey(key []byte) (addr types.AccAddress, domain string, account string) {
	splitBytes := bytes.SplitN(key, OwnerToAccountIndexSeparator, 3)
	if len(splitBytes) != 3 {
		panic(fmt.Sprintf("unexpected split length: %d", len(splitBytes)))
	}
	// convert back to their original types
	addr, domain, account = AccAddrFromIndex(splitBytes[0]), string(splitBytes[1]), string(splitBytes[2])
	return
}

// SplitOwnerToDomainKey takes an indexed owner to domain key
// and splits it into owner address and domain name
func SplitOwnerToDomainKey(key []byte) (addr types.AccAddress, domain string) {
	splitBytes := bytes.SplitN(key, OwnerToDomainIndexSeparator, 2)
	if len(splitBytes) != 2 {
		panic(fmt.Sprintf("expected split lenght: %d", len(splitBytes)))
	}
	// convert back to their original types
	addr, domain = AccAddrFromIndex(splitBytes[0]), string(splitBytes[1])
	return
}
