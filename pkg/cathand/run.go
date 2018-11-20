package cathand

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

func RunWait(name string, arg ...string) {
	fmt.Println(name + " " + strings.Join(arg, " "))

	cmd := exec.Command(name, arg...)
	cmd.Start()
	cmd.Wait()
}

func RunWaitWrite(name string, arg ...string) []byte {
	fmt.Println(name + " " + strings.Join(arg, " "))

	bytes, err := exec.Command(name, arg...).Output()
	AssertError(err)

	return bytes
}

func RunSequentialSignalWait(signalCh chan os.Signal, doneCount chan int, cmdFactory func(int) *exec.Cmd) {
	running := true

	count := 0
	for i := 0; running; i++ {
		count = i
		processDoneCh := make(chan error)
		pidCh := make(chan int)

		go func(i int) {
			cmd := cmdFactory(i)
			cmd.Start()
			pidCh <- cmd.Process.Pid
			processDoneCh <- cmd.Wait()
		}(count)

		checking := true
		pid := <-pidCh

		for checking {
			select {
			case <-signalCh:
				checking = false
				running = false

				process, _ := os.FindProcess(pid)
				process.Signal(os.Signal(syscall.SIGINT))
				break;
			case <-processDoneCh:
				checking = false
				break;
			default:
			}
		}
	}

	doneCount <- count
}

func RunWaitKillWrite(result chan []byte, group *sync.WaitGroup, name string, arg ...string) {
	fmt.Println(name + " " + strings.Join(arg, " "))

	var buffer bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &buffer

	err := cmd.Start()
	AssertError(err)

	group.Wait()

	cmd.Process.Signal(os.Signal(syscall.SIGINT))
	cmd.Wait()

	result <- buffer.Bytes()
}
