package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	atst "atst/pkg/atst"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "atst",
		Usage: "TODO",
		Action: func(cliCtx *cli.Context) error {
			_, cancel := context.WithCancel((cliCtx.Context))
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			done := make(chan bool)

			go func() {
				defer close(done)

				atst.StartManager(cliCtx.Args().Slice())

				done <- true
			}()

			for !<-done {
				select {
				case sig := <-sigChan:
					if sig == syscall.SIGINT {
						fmt.Printf("\nReceived signal: %v\n", sig)
						fmt.Println("Starting cleanup...")
						cancel()
					}
				}
			}

			fmt.Println("Returning")
			return nil
		},
	}

	app.Run(os.Args)

}
