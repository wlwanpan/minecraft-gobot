package launcher

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type launcherState int

const (
	NEWLINE_BYTE byte = '\n'

	MEM_CONVERSION int = 1024

	// launcher initialized but no operation was performed yet, sleeping state.
	LAUNCHER_STATE_INIT launcherState = iota
	// Minecraft server.jar successfully loaded and stdout 'DONE' is caught.
	LAUNCHER_STATE_READY
	// Minecraft serve.jar is still loading, has not yet caught 'Help' stdout.
	LAUNCHER_STATE_LOADING
)

var (
	ErrNotReady = errors.New("launcher is not ready")
)

func generateJavaRunCmd(ramAllocInGig int) *exec.Cmd {
	ramAllocInMb := strconv.Itoa(ramAllocInGig * MEM_CONVERSION)
	initialMemAlloc := strings.Join([]string{"-Xmx", ramAllocInMb, "M"}, "")
	maxMemAlloc := strings.Join([]string{"-Xms", ramAllocInMb, "M"}, "")
	return exec.Command("java", initialMemAlloc, maxMemAlloc, "-jar", "server.jar", "nogui")
}

type launcher struct {
	sync.Mutex
	execCmd   *exec.Cmd     // Keep a ref to raw exec cmd in case.
	currState launcherState // Keep track of the current launcher state.

	cmdin  chan string
	cmdout chan string
}

func newLauncher(memAlloc int) *launcher {
	return &launcher{
		execCmd:   generateJavaRunCmd(memAlloc),
		currState: LAUNCHER_STATE_INIT,
		cmdin:     make(chan string),
		cmdout:    make(chan string),
	}
}

func (l *launcher) transitionState(to launcherState) {
	l.Lock() // Should probably be a read/write lock instead

	// State: init -> loading
	switch l.currState {
	case LAUNCHER_STATE_INIT:
		if to != LAUNCHER_STATE_LOADING {
			log.Printf("Invalid state transition: %d -> %d", LAUNCHER_STATE_INIT, to)
			return
		}
	case LAUNCHER_STATE_LOADING:
		if to != LAUNCHER_STATE_READY {
			log.Printf("Invalid state transition: %d -> %d", LAUNCHER_STATE_LOADING, to)
			return
		}
	case LAUNCHER_STATE_READY:
		if to != LAUNCHER_STATE_INIT {
			log.Printf("Invalid state transition: %d -> %d", LAUNCHER_STATE_READY, to)
			return
		}
	default:
		log.Fatalf("Current state not handled: %d", l.currState)
		return
	}

	log.Printf("State transition: %d -> %d", l.currState, to)
	l.currState = to
	l.Unlock()
}

func (l *launcher) Launch(ctx context.Context) error {
	stdout, err := l.execCmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := l.execCmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := l.execCmd.Start(); err != nil {
		return err
	}

	l.transitionState(LAUNCHER_STATE_LOADING)

	// Spawn 2 go routine to handle std in and out.
	go l.processStdOut(ctx, stdout)
	go l.processStdIn(ctx, stdin)

	return nil
}

func (l *launcher) Stop(ctx context.Context) error {
	if err := l.execCmd.Process.Kill(); err != nil {
		return err
	}
	l.transitionState(LAUNCHER_STATE_INIT)
	return nil
}

func (l *launcher) processStdOut(ctx context.Context, stdout io.Reader) {
	bufr := bufio.NewReader(stdout)
	for {
		line, err := bufr.ReadString(NEWLINE_BYTE)
		if err != nil {
			if err == io.EOF {
				log.Printf("EOF reached: %s", line)
				break
			}
			log.Println(err)
		}

		log.Println(line)
		// Read stdout successful process output here.
		if strings.Contains(line, "Done") {
			l.transitionState(LAUNCHER_STATE_READY)
			continue
		}
	}
}

func (l *launcher) processStdIn(ctx context.Context, stdin io.Writer) {
	for {
		select {
		case cmd := <-l.cmdin:
			n, err := io.WriteString(stdin, cmd)
			if err != nil {
				log.Println(err)
				continue
			}
			if n != len(cmd) {
				log.Println("Error pushing command to stdin")
				continue
			}
			log.Println("Successfully pushed.")
		case <-ctx.Done():
			break
		}
	}
}
