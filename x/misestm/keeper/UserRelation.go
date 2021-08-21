package keeper

import (
	"encoding/binary"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mises-id/mises-tm/x/misestm/types"
)

// GetUserRelationCount get the total number of UserRelation
func (k Keeper) GetUserRelationCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserRelationCountKey))
	byteKey := types.KeyPrefix(types.UserRelationCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetUserRelationCount set the total number of UserRelation
func (k Keeper) SetUserRelationCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserRelationCountKey))
	byteKey := types.KeyPrefix(types.UserRelationCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendUserRelation appends a UserRelation in the store with a new id and update the count
func (k Keeper) AppendUserRelation(
	ctx sdk.Context,
	UserRelation types.UserRelation,
) uint64 {
	// Create the UserRelation
	count := k.GetUserRelationCount(ctx)

	// Set the ID of the appended value
	UserRelation.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserRelationKey))
	appendedValue := k.cdc.MustMarshal(&UserRelation)
	store.Set(GetUserRelationIDBytes(UserRelation.Id), appendedValue)

	// Update UserRelation count
	k.SetUserRelationCount(ctx, count+1)

	return count
}

// SetUserRelation set a specific UserRelation in the store
func (k Keeper) SetUserRelation(ctx sdk.Context, UserRelation types.UserRelation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserRelationKey))
	b := k.cdc.MustMarshal(&UserRelation)
	store.Set(GetUserRelationIDBytes(UserRelation.Id), b)
}

// GetUserRelation returns a UserRelation from its id
func (k Keeper) GetUserRelation(ctx sdk.Context, id uint64) types.UserRelation {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserRelationKey))
	var UserRelation types.UserRelation
	k.cdc.MustUnmarshal(store.Get(GetUserRelationIDBytes(id)), &UserRelation)
	return UserRelation
}

// HasUserRelation checks if the UserRelation exists in the store
func (k Keeper) HasUserRelation(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserRelationKey))
	return store.Has(GetUserRelationIDBytes(id))
}

// GetUserRelationOwner returns the creator of the UserRelation
func (k Keeper) GetUserRelationOwner(ctx sdk.Context, id uint64) string {
	return k.GetUserRelation(ctx, id).Creator
}

// RemoveUserRelation removes a UserRelation from the store
func (k Keeper) RemoveUserRelation(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserRelationKey))
	store.Delete(GetUserRelationIDBytes(id))
}

// GetAllUserRelation returns all UserRelation
func (k Keeper) GetAllUserRelation(ctx sdk.Context) (list []types.UserRelation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserRelationKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UserRelation
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetUserRelationIDBytes returns the byte representation of the ID
func GetUserRelationIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetUserRelationIDFromBytes returns ID in uint64 format from a byte array
func GetUserRelationIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
