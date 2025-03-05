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
			_, cancel := context.WithCancel((cliCtx.Context))
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			done := make(chan bool)

			go func() {
				defer close(done)

				programs := parsePrograms(os.Args)
				fmt.Printf("%v\n", programs)

				ch := atst.Start(programs)

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
