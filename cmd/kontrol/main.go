package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	errorcode "github.com/frantjc/go-error-code"
	"github.com/frantjc/kontrol/command"
)

func main() {
	var (
		ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		err       error
	)

	if err = command.NewKontrol().ExecuteContext(ctx); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}

	stop()
	os.Exit(errorcode.ExitCode(err))
}
