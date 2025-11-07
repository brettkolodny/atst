package cli

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type cmdOutput struct {
	Index int
	Text  string
}

func Simple(args []string) {
	cmds := make([]*exec.Cmd, 0)
	for _, cmdStr := range args {
		cmdParts := strings.Split(cmdStr, " ")
		if len(cmdParts) == 0 {
			continue
		}

		name := cmdParts[0]
		args := cmdParts[1:]

		cmds = append(cmds, exec.Command(name, args...))
	}

	printerDone := make(chan struct{})
	printerCh := make(chan (*cmdOutput))

	var cmdWg sync.WaitGroup
	for index, cmd := range cmds {
		cmdWg.Add(1)
		go func() {
			defer cmdWg.Done()

			stderr, err := cmd.StderrPipe()
			if err != nil {
				fmt.Printf("[%d]: Failed to run, %v\n", index, err)
				return
			}
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Printf("[%d]: Failed to run, %v\n", index, err)
				return
			}

			if err := cmd.Start(); err != nil {
				fmt.Printf("[%d]: Failed to run, %v\n", index, err)
				return
			}

			stderrScanner := bufio.NewScanner(stderr)
			stdoutScanner := bufio.NewScanner(stdout)

			// The Go docs state that "It is thus incorrect to call Wait before all reads from the pipe have completed."
			// so we create an aditional wait group to ensure the scanners are closed before calling cmd.Wait
			var scannerWg sync.WaitGroup
			scannerWg.Add(2)
			go func() {
				writeScannerOutput(stderrScanner, index, printerCh)
				scannerWg.Done()
			}()
			go func() {
				writeScannerOutput(stdoutScanner, index, printerCh)
				scannerWg.Done()
			}()
			scannerWg.Wait()

			if err = cmd.Wait(); err != nil {
				fmt.Printf("[%d] %v\n", index, err)
			} else {
				fmt.Printf("[%d] exited\n", index)
			}
		}()
	}

	go func() {
		for output := range printerCh {
			fmt.Printf("[%d] %s\n", output.Index, output.Text)
		}
		close(printerDone)
	}()

	cmdWg.Wait()
	close(printerCh)
	<-printerDone
}

func writeScannerOutput(scanner *bufio.Scanner, index int, ch chan *cmdOutput) {
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		out := scanner.Text()
		ch <- &cmdOutput{
			Index: index,
			Text:  out,
		}
	}
}
