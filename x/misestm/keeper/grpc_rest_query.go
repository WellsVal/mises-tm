package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mises-id/mises-tm/x/misestm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.RestQueryServer = Keeper{}

func (k Keeper) QueryDid(c context.Context, req *types.RestQueryDidRequest) (*types.RestQueryDidResponse, error) {
 	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var DidRegistry types.DidRegistry
	ctx := sdk.UnwrapSDKContext(c)
	userMgr := NewUserMgrImpl(k)
	misesAcc, err := userMgr.GetUserAccount(ctx, req.MisesId)
	if err != nil {
		return nil, err
	}
	if misesAcc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "mises id %s not exists", req.MisesId)
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DidRegistryKey))
	k.cdc.Unmarshal(store.Get(GetDidRegistryIDBytes(misesAcc.DidRegistryID)), &DidRegistry)

	return &types.RestQueryDidResponse{DidRegistry: &DidRegistry}, nil
}

func (k Keeper) QueryUser(c context.Context, req *types.RestQueryUserRequest) (*types.RestQueryUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var UserInfo types.UserInfo
	ctx := sdk.UnwrapSDKContext(c)
	userMgr := NewUserMgrImpl(k)
	misesAcc, err := userMgr.GetUserAccount(ctx, req.MisesId)
	if err != nil {
		return nil, err
	}

	if misesAcc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "mises id %s not exists", req.MisesId)
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UserInfoKey))
	k.cdc.Unmarshal(store.Get(GetUserInfoIDBytes(misesAcc.UserInfoID)), &UserInfo)

	pubInfo := types.PublicUserInfo{
	}
	priInfo := types.PrivateUserInfo{
		EncData: UserInfo.EncData,
		Iv:   UserInfo.Iv,
	}

	return &types.RestQueryUserResponse{PubInfo: &pubInfo, PriInfo:&priInfo}, nil
}

// query user relations
func (k Keeper) QueryUserRelation(c context.Context, req *types.RestQueryUserRelationRequest) (*types.RestQueryUserRelationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	userMgr := NewUserMgrImpl(k)
	misesAcc, err := userMgr.GetUserAccount(ctx, req.MisesId)
	if err != nil {
		return nil, err
	}

	if misesAcc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "mises id %s not exists", req.MisesId)
	}
	pagination := req.Pagination 
	if pagination == nil {
		pagination = &query.PageRequest{
			Key: []byte(""),
			Limit: 100,
		}
	}
	relType := uint64(0)
	if req.Filter == "following" {
		relType = types.RelTypeBitFollow
	}
	UserRelations, err := userMgr.GetUserRelations(ctx, relType, req.MisesId, string(pagination.Key), int(pagination.Limit))
	if err != nil {
		return nil, err
	}
	misesList := []*types.MisesID{}
	for _, r := range UserRelations {
		misesList = append(misesList, &types.MisesID{MisesId: r.UidTo})
	}
	nextKey := ""
	if len(misesList) > 0 {
		nextKey = misesList[len(misesList)-1].MisesId
	}
	pageRes := &query.PageResponse{NextKey: []byte(nextKey)}

	resp := types.RestQueryUserRelationResponse{
		MisesList: misesList,
		Pagination:   pageRes,
	}
	return &resp, nil
}
