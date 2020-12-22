package initproc

// based on https://github.com/pablo-ruth/go-init/blob/master/main.go
import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
)

func RunInit(kctrlArgs []string, log logr.Logger) {
	kctrlArgs = append([]string{"-c"}, kctrlArgs...)

	log.Info("starting zombie reaper")
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go removeZombies(ctx, &wg)

	cmd := exec.Command("kapp-controller", kctrlArgs...)
	err, retCode := run(cmd)
	if err != nil {
		log.Error(err, "Could not start controller")
	}

	cleanQuit(cancel, &wg, retCode)
}

func removeZombies(ctx context.Context, wg *sync.WaitGroup) {
	for {
		var status syscall.WaitStatus
		pid, _ := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)

		if pid <= 0 {
			time.Sleep(1 * time.Second)
		} else {
			continue
		}

		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
		}
	}
}

func run(cmd *exec.Cmd) (error, int) {
	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs)
	defer signal.Reset()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// forward signals to child's proc group
	go func() {
		for sig := range sigs {
			if sig != syscall.SIGCHLD {
				syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
			}
		}
	}()

	if err := cmd.Run(); err != nil {
		return err, 1
	}

	return nil, 0
}

func cleanQuit(cancel context.CancelFunc, wg *sync.WaitGroup, retCode int) {
	cancel()
	wg.Wait()

	os.Exit(retCode)
}
