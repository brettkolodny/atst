package atst

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type Output struct {
	Command string
	Index   int
	Msg     string
}

type Program struct {
	Exec string
	Args []string
}

type Atst struct {
	Outputs []chan Output
	wg      *sync.WaitGroup
}

func (a Atst) Wait() {
	a.wg.Wait()
}

func Start(programs []Program) Atst {
	var wg sync.WaitGroup

	outputChannels := []chan Output{}

	for index, program := range programs {
		wg.Add(1)

		ch := make(chan Output)
		outputChannels = append(outputChannels, ch)

		go func() {
			defer wg.Done()

			cmd := exec.Command(program.Exec, program.Args...)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				ch <- Output{
					Command: program.Exec,
					Index:   index,
					Msg:     fmt.Sprintf("%s\n", err),
				}
				return
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				ch <- Output{
					Command: program.Exec,
					Index:   index,
					Msg:     fmt.Sprintf("%s\n", err),
				}
				return
			}

			if err := cmd.Start(); err != nil {
				ch <- Output{
					Command: program.Exec,
					Index:   index,
					Msg:     fmt.Sprintf("%s\n", err),
				}
				return
			}

			writeFromScanner := func(reader io.Reader) {
				scanner := bufio.NewScanner(reader)
				for scanner.Scan() {
					ch <- Output{
						Command: program.Exec,
						Index:   index,
						Msg:     scanner.Text(),
					}
				}

			}

			go writeFromScanner(stdout)
			go writeFromScanner(stderr)

			if err := cmd.Wait(); err != nil {
				ch <- Output{
					Command: program.Exec,
					Index:   index,
					Msg:     fmt.Sprintf("%s\n", err),
				}
			}
		}()
	}

	return Atst{
		Outputs: outputChannels,
		wg:      &wg,
	}
}
