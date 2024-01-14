package keeper

import (
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// AccountKeeper is the interface contract that x/auth's keeper implements.
type AccountKeeper interface {
	// Return a new account with the next account number and the specified address. Does not save the new account to the store.
	NewAccountWithAddress(sdk.Context, sdk.AccAddress) types.AccountI

	// Return a new account with the next account number. Does not save the new account to the store.
	NewAccount(sdk.Context, types.AccountI) types.AccountI

	// Check if an account exists in the store.
	HasAccount(sdk.Context, sdk.AccAddress) bool

	// Retrieve an account from the store.
	GetAccount(sdk.Context, sdk.AccAddress) types.AccountI

	// GetAllAccounts returns all accounts in the accountKeeper.
	GetAllAccounts(sdk.Context) []types.AccountI

	// Set an account in the store.
	SetAccount(sdk.Context, types.AccountI)

	// Remove an account from the store.
	RemoveAccount(sdk.Context, types.AccountI)

	// Iterate over all accounts, calling the provided function. Stop iteration when it returns true.
	IterateAccounts(sdk.Context, func(types.AccountI) bool)

	types.QueryServer

	// Logger returns a module-specific logger.
	Logger(ctx sdk.Context) log.Logger

	// Fetch the public key of an account at a specified address
	GetPubKey(sdk.Context, sdk.AccAddress) (cryptotypes.PubKey, error)

	// Fetch the sequence of an account at a specified address.
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)

	// Fetch the next account number, and increment the internal counter.
	GetNextAccountNumber(sdk.Context) uint64

	// ValidatePermissions validates that the module account has been granted
	// permissions within its set of allowed permissions.
	ValidatePermissions(types.ModuleAccountI) error

	// GetModuleAddress returns an address based on the module name
	GetModuleAddress(string) sdk.AccAddress

	// GetModuleAddressAndPermissions returns an address and permissions based on the module name
	GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string)

	// GetModuleAccountAndPermissions gets the module account from the auth account store and its
	// registered permissions
	GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (types.ModuleAccountI, []string)

	// GetModuleAccount gets the module account from the auth account store, if the account does not
	// exist in the AccountKeeper, then it is created.
	GetModuleAccount(sdk.Context, string) types.ModuleAccountI

	// SetModuleAccount sets the module account to the auth account store
	SetModuleAccount(sdk.Context, types.ModuleAccountI)

	// MarshalAccount protobuf serializes an Account interface
	MarshalAccount(types.AccountI) ([]byte, error)

	// UnmarshalAccount returns an Account interface from raw encoded account
	// bytes of a Proto-based Account type
	UnmarshalAccount([]byte) (types.AccountI, error)

	// GetCodec return codec.Codec object used by the keeper
	GetCodec() codec.BinaryCodec

	// GetParams gets the auth module's parameters.
	GetParams(sdk.Context) types.Params

	// SetParams sets the auth module's parameters.
	SetParams(sdk.Context, types.Params)

	// HasAccountAddressByID checks account address exists by id.
	HasAccountAddressByID(ctx sdk.Context, id uint64) bool

	// GetAccountAddressById returns account address by id.
	GetAccountAddressByID(ctx sdk.Context, id uint64) string

	// InitGenesis - Init store state from genesis data
	InitGenesis(ctx sdk.Context, data types.GenesisState)

	// ExportGenesis returns a GenesisState for a given context and keeper
	ExportGenesis(ctx sdk.Context) *types.GenesisState
}

// accountKeeper encodes/decodes accounts using the go-amino (binary)
// encoding/decoding library.
type accountKeeper struct {
	key           storetypes.StoreKey
	cdc           codec.BinaryCodec
	paramSubspace paramtypes.Subspace
	permAddrs     map[string]types.PermissionsForAddress

	// The prototypical AccountI constructor.
	proto      func() types.AccountI
	addressCdc address.Codec
}

var _ AccountKeeper = &accountKeeper{}

// NewAccountKeeper returns a new AccountKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.Accounts.
// `maccPerms` is a map that takes accounts' addresses as keys, and their respective permissions as values. This map is used to construct
// types.PermissionsForAddress and is used in keeper.ValidatePermissions. Permissions are plain strings,
// and don't have to fit into any predefined structure. This auth module does not use account permissions internally, though other modules
// may use auth.Keeper to access the accounts permissions map.
func NewAccountKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramstore paramtypes.Subspace, proto func() types.AccountI,
	maccPerms map[string][]string, bech32Prefix string,
) AccountKeeper {
	// set KeyTable if it has not already been set
	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	}

	permAddrs := make(map[string]types.PermissionsForAddress)
	for name, perms := range maccPerms {
		permAddrs[name] = types.NewPermissionsForAddress(name, perms)
	}

	bech32Codec := newBech32Codec(bech32Prefix)

	return accountKeeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore,
		permAddrs:     permAddrs,
		addressCdc:    bech32Codec,
	}
}

// Logger returns a module-specific logger.
func (ak accountKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetPubKey Returns the PubKey of the account at address
func (ak accountKeeper) GetPubKey(ctx sdk.Context, addr sdk.AccAddress) (cryptotypes.PubKey, error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
	}

	return acc.GetPubKey(), nil
}

// GetSequence Returns the Sequence of the account at address
func (ak accountKeeper) GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
	}

	return acc.GetSequence(), nil
}

// GetNextAccountNumber returns and increments the global account number counter.
// If the global account number is not set, it initializes it with value 0.
func (ak accountKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(ak.key)

	bz := store.Get(types.GlobalAccountNumberKey)
	if bz == nil {
		// initialize the account numbers
		accNumber = 0
	} else {
		val := gogotypes.UInt64Value{}

		err := ak.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}

		accNumber = val.GetValue()
	}

	bz = ak.cdc.MustMarshal(&gogotypes.UInt64Value{Value: accNumber + 1})
	store.Set(types.GlobalAccountNumberKey, bz)

	return accNumber
}

// ValidatePermissions validates that the module account has been granted
// permissions within its set of allowed permissions.
func (ak accountKeeper) ValidatePermissions(macc types.ModuleAccountI) error {
	permAddr := ak.permAddrs[macc.GetName()]
	for _, perm := range macc.GetPermissions() {
		if !permAddr.HasPermission(perm) {
			return fmt.Errorf("invalid module permission %s", perm)
		}
	}

	return nil
}

// GetModuleAddress returns an address based on the module name
func (ak accountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	permAddr, ok := ak.permAddrs[moduleName]
	if !ok {
		return nil
	}

	return permAddr.GetAddress()
}

// GetModuleAddressAndPermissions returns an address and permissions based on the module name
func (ak accountKeeper) GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string) {
	permAddr, ok := ak.permAddrs[moduleName]
	if !ok {
		return addr, permissions
	}

	return permAddr.GetAddress(), permAddr.GetPermissions()
}

// GetModuleAccountAndPermissions gets the module account from the auth account store and its
// registered permissions
func (ak accountKeeper) GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (types.ModuleAccountI, []string) {
	addr, perms := ak.GetModuleAddressAndPermissions(moduleName)
	if addr == nil {
		return nil, []string{}
	}

	acc := ak.GetAccount(ctx, addr)
	if acc != nil {
		macc, ok := acc.(types.ModuleAccountI)
		if !ok {
			panic("account is not a module account")
		}
		return macc, perms
	}

	// create a new module account
	macc := types.NewEmptyModuleAccount(moduleName, perms...)
	maccI := (ak.NewAccount(ctx, macc)).(types.ModuleAccountI) // set the account number
	ak.SetModuleAccount(ctx, maccI)

	return maccI, perms
}

// GetModuleAccount gets the module account from the auth account store, if the account does not
// exist in the AccountKeeper, then it is created.
func (ak accountKeeper) GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI {
	acc, _ := ak.GetModuleAccountAndPermissions(ctx, moduleName)
	return acc
}

// SetModuleAccount sets the module account to the auth account store
func (ak accountKeeper) SetModuleAccount(ctx sdk.Context, macc types.ModuleAccountI) {
	ak.SetAccount(ctx, macc)
}

func (ak accountKeeper) decodeAccount(bz []byte) types.AccountI {
	acc, err := ak.UnmarshalAccount(bz)
	if err != nil {
		panic(err)
	}

	return acc
}

// MarshalAccount protobuf serializes an Account interface
func (ak accountKeeper) MarshalAccount(accountI types.AccountI) ([]byte, error) { // nolint:interfacer
	return ak.cdc.MarshalInterface(accountI)
}

// UnmarshalAccount returns an Account interface from raw encoded account
// bytes of a Proto-based Account type
func (ak accountKeeper) UnmarshalAccount(bz []byte) (types.AccountI, error) {
	var acc types.AccountI
	return acc, ak.cdc.UnmarshalInterface(bz, &acc)
}

// GetCodec return codec.Codec object used by the keeper
func (ak accountKeeper) GetCodec() codec.BinaryCodec { return ak.cdc }

// add getter for bech32Prefix
func (ak accountKeeper) getBech32Prefix() (string, error) {
	bech32Codec, ok := ak.addressCdc.(bech32Codec)
	if !ok {
		return "", fmt.Errorf("unable cast addressCdc to bech32Codec; expected %T got %T", bech32Codec, ak.addressCdc)
	}

	return bech32Codec.bech32Prefix, nil
}
