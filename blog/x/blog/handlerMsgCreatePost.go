package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/fadeev/blog/x/blog/types"
)

func handleMsgCreatePost(ctx sdk.Context, k Keeper, msg MsgCreatePost) (*sdk.Result, error) {
	var post = types.Post{
		Creator: msg.Creator,
		ID:      msg.ID,
    Title: msg.Title,
    Body: msg.Body,
	}
	k.CreatePost(ctx, post)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
