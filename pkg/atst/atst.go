package atst

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

type Manager = struct {
	programs []chan struct{}
}

func StartManager(programs []string) {
	var wg sync.WaitGroup

	for index, program := range programs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			cmd := exec.Command(program)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Printf("[%d]: %s\n", index, err)
				return
			}

			if err := cmd.Start(); err != nil {
				fmt.Printf("[%d]: %s\n", index, err)
				return
			}

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Printf("[%d]: %s\n", index, scanner.Text())
			}

			if err := cmd.Wait(); err != nil {
				fmt.Printf("[%d]: %s\n", index, err)
			}
		}()
	}

	wg.Wait()
}
