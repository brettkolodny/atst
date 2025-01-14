package atst

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type Output = struct {
	Command string
	Index   int
	Msg     string
}

func Start(programs []string) chan Output {
	var wg sync.WaitGroup

	ch := make(chan Output)

	for index, program := range programs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			args := strings.Split(program, " ")
			cmd := exec.Command(args[0], args[1:]...)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				ch <- Output{
					Command: program,
					Index:   index,
					Msg:     fmt.Sprintf("%s\n", err),
				}
				return
			}

			if err := cmd.Start(); err != nil {
				ch <- Output{
					Command: program,
					Index:   index,
					Msg:     fmt.Sprintf("%s\n", err),
				}
				return
			}

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				ch <- Output{
					Command: program,
					Index:   index,
					Msg:     scanner.Text(),
				}
			}

			if err := cmd.Wait(); err != nil {
				ch <- Output{
					Command: program,
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
