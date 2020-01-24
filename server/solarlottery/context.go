package solarlottery

import (
	"context"
	"fmt"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/config"
)

var apiContextKey = config.Repository + "/" + fmt.Sprintf("%T", solarLottery{})

func Context(ctx context.Context, sl SolarLottery) context.Context {
	ctx = context.WithValue(ctx, apiContextKey, sl)
	return ctx
}

func FromContext(ctx context.Context) SolarLottery {
	return ctx.Value(apiContextKey).(SolarLottery)
}
