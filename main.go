package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var usage = `
Usage:
para 'cmd1' 'cmd2'...
`
var wg = sync.WaitGroup{}

func work(command string, id int) {
	wg.Add(1)
	cmdArr := strings.Split(command, " ")
	cmdName := cmdArr[0]
	cmdArgs := cmdArr[1:]

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("Worker %d | %s\n", id, scanner.Text())
		}
	}()

	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		fmt.Printf("id: %d %s", id, err)
	}
	wg.Done()
}

func main() {
	if os.Args[1] == "-h" {
		fmt.Println(usage)
		return
	}
	l := len(os.Args)
	for id, command := range os.Args[1 : l-1] {

		go work(command, id)
	}
	// idk why but it needs one blocking function.
	work(os.Args[l-1], l-2)
	wg.Wait()
}
