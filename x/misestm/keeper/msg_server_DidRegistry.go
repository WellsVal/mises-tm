package keeper

import (
	"context"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mises-id/mises-tm/x/misestm/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
)

func (k msgServer) CreateDidRegistry(goCtx context.Context, msg *types.MsgCreateDidRegistry) (*types.MsgCreateDidRegistryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var DidRegistry = types.DidRegistry{
		Creator:       msg.Creator,
		Did:           msg.Did,
		PkeyDid:       msg.PkeyDid,
		PkeyType:      msg.PkeyType,
		PkeyMultibase: msg.PkeyMultibase,
		Version:       msg.Version,
	}
	ak := k.ak
	userMgr := NewUserMgrImpl(k.Keeper)
	misesAcc, err := userMgr.GetUserAccount(ctx, DidRegistry.Did)
	if misesAcc != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "account %s already exists", msg.Did)
	}
	if err != nil {
		return nil, err
	}
	addr, didType, err := types.AddrFromDid(DidRegistry.Did)
	if err != nil {
		return nil, err
	}
	baseAccount := ak.NewAccountWithAddress(ctx, addr)
	if _, ok := baseAccount.(*authtypes.BaseAccount); !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid account type; expected: BaseAccount, got: %T", baseAccount)
	}
	var acc authtypes.AccountI = baseAccount
	pubKeyBytes := base58.Decode(DidRegistry.PkeyMultibase)
	pubKey := secp256k1.PubKey{Key: pubKeyBytes}
	err = acc.SetPubKey(&pubKey)
	if err != nil {
		return nil, err
	}
	ak.SetAccount(ctx, acc)

	defer func() {
		telemetry.IncrCounter(1, "new", "account")
	}()

	regID := k.AppendDidRegistry(
		ctx,
		DidRegistry,
	)

	var infoID uint64
	if didType == types.DIDTypeUser {
		infoID = k.AppendUserInfo(
			ctx,
			types.UserInfo{
				Creator: DidRegistry.Creator,
				Uid:     DidRegistry.Did,
			},
		)

	} else if didType == types.DIDTypeApp {
		infoID = k.AppendAppInfo(
			ctx,
			types.AppInfo{
				Creator: DidRegistry.Creator,
				Appid:   DidRegistry.Did,
			},
		)
	} else {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "did %s not supported", msg.Did)
	}

	newMisesAcc := types.MisesAccount{
		MisesID:       DidRegistry.Did,
		DidRegistryID: regID,
		InfoID:        infoID,
		DidType:       didType,
	}
	k.SetMisesAccount(ctx, newMisesAcc)

	return &types.MsgCreateDidRegistryResponse{}, nil
}
