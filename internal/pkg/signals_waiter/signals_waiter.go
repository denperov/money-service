package signals_waiter

import (
	"context"
	"os"
	"os/signal"
)

func Wait(ctx context.Context, signals []os.Signal) os.Signal {
	sig := make(chan os.Signal, len(signals))
	signal.Notify(sig, signals...)

	select {
	case <-ctx.Done(): // cancellation
		return nil
	case s := <-sig:
		return s
	}
}
