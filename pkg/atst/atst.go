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

func Start(programs []Program) chan Output {
	var wg sync.WaitGroup

	ch := make(chan Output)

	for index, program := range programs {
		wg.Add(1)

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

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}
