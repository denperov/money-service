package stop_signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type stopSignal struct {
	signalHandler func()

	stop func()
	done <-chan struct{}
}

func New(
	signalHandler func(),
) *stopSignal {
	ctx, cancel := context.WithCancel(context.Background())
	return &stopSignal{
		signalHandler: signalHandler,
		stop:          cancel,
		done:          ctx.Done(),
	}
}

func (w *stopSignal) Start(_ context.Context) {
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	sig := make(chan os.Signal, len(signals))
	signal.Notify(sig, signals...)

	go func() {
		select {
		case <-w.done: // cancellation
			return
		case <-sig:
			w.signalHandler()
			return
		}
	}()
}

func (w *stopSignal) Stop() {
	w.stop()
}
