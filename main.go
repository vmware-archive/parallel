package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var (
	silent         = flag.Bool("s", false, "suppress all command output")
	burstCount     = flag.Int("b", 10, "burst count")
	iterationCount = flag.Int("i", 5, "iteration count")
)

func main() {
	if len(os.Args) < 2 {
		panic("need command to run!")
	}

	flag.Parse()

	errChan := make(chan error, *burstCount)

	for i := 0; i < *iterationCount; i++ {
		fmt.Println(fmt.Sprintf("%d/%d", i+1, *iterationCount))
		for b := 0; b < *burstCount; b++ {
			cmdPath := flag.Args()[0]

			go func() {
				errChan <- runCmd(cmdPath, flag.Args()[1:]...)
			}()
		}

		for b := 0; b < *burstCount; b++ {
			err := <-errChan
			if err != nil {
				panic(err)
			}
		}

		time.Sleep(1 * time.Second)
	}

}

func runCmd(cmdPath string, args ...string) error {
	cmd := exec.Command(cmdPath, args...)

	if !*silent {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil

}
