package executor

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
	"os"
	"testing"
	"time"
)

var testCtx sdk.Context
var testKey = sdk.NewKVStoreKey("test")
var testCdc *codec.Codec
var testKeeper keeper.Keeper
var testAccount = types.Account{
	Domain:     "a-super-domain",
	Name:       utils.StrPtr("a-super-account"),
	Owner:      keeper.CharlieKey,
	ValidUntil: 10000,
	Resources: []types.Resource{
		{
			URI:      "a-super-uri",
			Resource: "a-super-res",
		},
	},
	Certificates: []types.Certificate{types.Certificate("a-random-cert")},
	Broker:       nil,
	MetadataURI:  "metadata",
}

var aliceKey sdk.AccAddress
var bobKey sdk.AccAddress

func newTest() error {
	_, addr := utils.GeneratePrivKeyAddressPairs(2)
	aliceKey = addr[0]
	bobKey = addr[1]
	testCdc = codec.New()
	mdb := db.NewMemDB()
	ms := store.NewCommitMultiStore(mdb)
	ms.MountStoreWithDB(testKey, sdk.StoreTypeIAVL, mdb)
	err := ms.LoadLatestVersion()
	if err != nil {
		return err
	}
	testCtx = sdk.NewContext(ms, tmtypes.Header{Time: time.Now()}, true, log.NewNopLogger())
	testKeeper = keeper.NewKeeper(testCdc, testKey, nil, nil, nil)
	testKeeper.AccountStore(testCtx).Create(&testAccount)
	return nil
}

func TestMain(m *testing.M) {
	err := newTest()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
