package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	atst "github.com/brettkolodny/atst/pkg/atst"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "atst",
		Description: "Run multiple CLI programs at the same time",
		Usage:       "atst <command...>",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "program",
				Aliases: []string{"p"},
				Usage:   "Specify a program to run",
			},
			&cli.StringSliceFlag{
				Name:    "arg",
				Aliases: []string{"a"},
				Usage:   "Specify an argument for the most recently defined program",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			// Create a cancellable context
			ctx, cancel := context.WithCancel(cliCtx.Context)
			defer cancel()

			// Set up signal handling
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			// Parse command line arguments
			programs := parsePrograms(os.Args)
			if len(programs) == 0 {
				return fmt.Errorf("no programs specified")
			}

			// Start a goroutine to handle signals
			go func() {
				sig := <-sigChan
				fmt.Printf("\nReceived signal: %v\n", sig)
				cancel() // Signal all processes to terminate
			}()

			// Start all programs
			a := atst.Start(programs)

			// Merge all output channels into one
			outputCh := make(chan atst.Output)
			for _, ch := range a.Outputs {
				go func(ch chan atst.Output) {
					for output := range ch {
						select {
						case outputCh <- output:
						case <-ctx.Done():
							return
						}
					}
				}(ch)
			}

			// Handle program output
			go func() {
				for output := range outputCh {
					fmt.Printf("[%d]: %s\n", output.Index, strings.TrimRight(output.Msg, "\n\r"))
				}
			}()

			// Wait for programs to complete or context to be cancelled
			done := make(chan struct{})
			go func() {
				a.Wait()
				close(done)
			}()

			// Wait for either completion or cancellation
			select {
			case <-done:
				fmt.Println("All programs completed")
			case <-ctx.Done():
				fmt.Println("Shutting down...")
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parsePrograms(args []string) []atst.Program {
	programs := []atst.Program{} // Initialize as empty slice, not nil
	var currentProgram string

	for i := 0; i < len(args); i++ {
		// Check if current arg is a program flag
		if args[i] == "-p" || args[i] == "--program" {
			if i+1 < len(args) {
				// The next item is the program name - we'll accept anything here, even if it looks like a flag
				progName := args[i+1]
				programs = append(programs, atst.Program{
					Exec: progName,
					Args: []string{},
				})
				currentProgram = progName
				i++ // Skip the program name
			}
		} else if args[i] == "-a" || args[i] == "--arg" {
			if i+1 < len(args) && currentProgram != "" {
				// Accept any argument value, even if it starts with - or --
				argValue := args[i+1]

				// Find the current program and add the arg
				for j := len(programs) - 1; j >= 0; j-- {
					if programs[j].Exec == currentProgram {
						programs[j].Args = append(programs[j].Args, argValue)
						break
					}
				}
				i++ // Skip the arg value
			}
		}
	}

	return programs
}

func isFlag(s string) bool {
	return len(s) > 0 && s[0] == '-'
}
