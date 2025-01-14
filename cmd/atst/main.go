package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	atst "atst/pkg/atst"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "atst",
		Description: "Run multiple CLI programs at the same time",
		Usage:       "atst <command...>",
		Action: func(cliCtx *cli.Context) error {
			_, cancel := context.WithCancel((cliCtx.Context))
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			done := make(chan bool)

			go func() {
				defer close(done)

				ch := atst.Start(cliCtx.Args().Slice())

				for v := range ch {

					fmt.Printf("[%d]: %s\n", v.Index, strings.TrimRight(v.Msg, "\n\r"))
				}

				done <- true
			}()

			for !<-done {
				select {
				case sig := <-sigChan:
					if sig == syscall.SIGINT {
						cancel()
					}
				}
			}

			return nil
		},
	}

	app.Run(os.Args)

}
