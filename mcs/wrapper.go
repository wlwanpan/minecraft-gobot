package mcs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/looplab/fsm"
	"github.com/wlwanpan/minecraft-gobot/config"
)

const (
	NEWLINE_BYTE byte = '\n'

	MEM_CONVERSION int = 1024
	// mcs initialized but no operation was performed yet, sleeping state.
	WRAPPER_STATE_OFFLINE string = "offline"
	// Minecraft server.jar successfully loaded and stdout 'DONE' is caught.
	WRAPPER_STATE_ONLINE string = "online"
	// Minecraft serve.jar is still loading, has not yet caught 'Help' stdout.
	WRAPPER_STATE_LOADING string = "loading"
	// Minecraft server.jar is still running, but in the process of 'stopping'.
	WRAPPER_STATE_STOPPING string = "stopping"
)

var (
	ErrServerAlreadyLoading = errors.New("server is already loading")

	ErrServerAlreadyOnline = errors.New("server is already online")

	ErrServerAlreadyOffline = errors.New("server is already offline")

	ErrSessionAlreadyStopped = errors.New("session already stopped")
)

func generateJavaRunCmd(ramAllocInGig int) *exec.Cmd {
	ramAllocInMb := strconv.Itoa(ramAllocInGig * MEM_CONVERSION)
	initialMemAlloc := strings.Join([]string{"-Xmx", ramAllocInMb, "M"}, "")
	maxMemAlloc := strings.Join([]string{"-Xms", ramAllocInMb, "M"}, "")
	serverjar := config.Cfg.Mcs.Serverjar
	return exec.Command("java", initialMemAlloc, maxMemAlloc, "-jar", serverjar, "nogui")
}

type console struct {
	execCmd *exec.Cmd
	cmdout  *bufio.Reader
	cmdin   *bufio.Writer
}

func (c *console) execJava(mem int) error {
	c.execCmd = generateJavaRunCmd(mem)

	stdout, err := c.execCmd.StdoutPipe()
	if err != nil {
		return err
	}
	c.cmdout = bufio.NewReader(stdout)

	stdin, err := c.execCmd.StdinPipe()
	if err != nil {
		return err
	}
	c.cmdin = bufio.NewWriter(stdin)

	return c.execCmd.Start()
}

func (c *console) write(cmd string) error {
	wrappedCmd := fmt.Sprintf("%s\r\n", cmd)
	_, err := c.cmdin.WriteString(wrappedCmd)
	if err != nil {
		return err
	}
	return c.cmdin.Flush()
}

func (c *console) readLine() (string, error) {
	return c.cmdout.ReadString(NEWLINE_BYTE)
}

func (c *console) kill() error {
	return c.execCmd.Process.Kill()
}

type wrapper struct {
	sync.RWMutex
	stateMachine *fsm.FSM
	console      *console
	done         chan bool

	gameSess    *gameSession
	lastLogLine string
}

func newWrapper() *wrapper {
	return &wrapper{
		console: &console{},
		done:    make(chan bool),
		// stateMachine: offline -> loading -> online -> stopping >>
		stateMachine: fsm.NewFSM(
			WRAPPER_STATE_OFFLINE,
			fsm.Events{
				fsm.EventDesc{
					Name: WRAPPER_STATE_OFFLINE,
					Src:  []string{WRAPPER_STATE_STOPPING},
					Dst:  WRAPPER_STATE_LOADING,
				},
				fsm.EventDesc{
					Name: WRAPPER_STATE_LOADING,
					Src:  []string{WRAPPER_STATE_OFFLINE},
					Dst:  WRAPPER_STATE_ONLINE,
				},
				fsm.EventDesc{
					Name: WRAPPER_STATE_ONLINE,
					Src:  []string{WRAPPER_STATE_LOADING},
					Dst:  WRAPPER_STATE_STOPPING,
				},
				fsm.EventDesc{
					Name: WRAPPER_STATE_STOPPING,
					Src:  []string{WRAPPER_STATE_ONLINE},
					Dst:  WRAPPER_STATE_OFFLINE,
				},
			},
			map[string]fsm.Callback{
				"enter_state": func(e *fsm.Event) {
					log.Printf("State changes: %s -> %s", e.Src, e.Dst)
				},
			},
		),
	}
}

func (w *wrapper) isLoading() bool {
	return w.stateMachine.Current() == WRAPPER_STATE_LOADING
}

func (w *wrapper) isStopping() bool {
	return w.stateMachine.Current() == WRAPPER_STATE_STOPPING
}

func (w *wrapper) isOnline() bool {
	return w.stateMachine.Current() == WRAPPER_STATE_ONLINE
}

func (w *wrapper) isOffline() bool {
	return w.stateMachine.Current() == WRAPPER_STATE_OFFLINE
}

func (w *wrapper) startGameSession() {
	if w.gameSess != nil {
		w.gameSess.stop()
	}

	log.Println("Starting new game session")
	w.gameSess = newGameSession(w.console)
	w.gameSess.start()
}

func (w *wrapper) stopGameSession() error {
	if w.gameSess == nil {
		return ErrSessionAlreadyStopped
	}

	if err := w.saveGameSession(); err != nil {
		return err
	}

	log.Printf("Stopping current game session, game-tick=%d", w.gameSess.gametick)
	w.gameSess.stop()
	w.gameSess = nil
	return nil
}

func (w *wrapper) saveGameSession() error {
	err := w.gameSess.save()
	if err == ErrSavingGameTimedOut {
		// TODO: fix 'save the game' log actionable.
		return nil
	}
	return err
}

func (w *wrapper) pushUpdateToGameSession(update logUpdate) {
	if w.gameSess == nil {
		return
	}
	w.gameSess.updates <- update
}

func (w *wrapper) start(mem int) error {
	if w.isLoading() {
		return ErrServerAlreadyLoading
	}
	if w.isOnline() {
		return ErrServerAlreadyOnline
	}

	w.stateMachine.SetState(WRAPPER_STATE_LOADING)
	if err := w.console.execJava(mem); err != nil {
		w.stateMachine.SetState(WRAPPER_STATE_OFFLINE)
		return err
	}

	go w.processCmdOut() // exits on EOF

	return nil
}

func (w *wrapper) stop() error {
	if !w.isOnline() {
		return ErrServerAlreadyOffline
	}
	w.stateMachine.SetState(WRAPPER_STATE_STOPPING)

	if err := w.stopGameSession(); err != nil {
		log.Printf("err stopping game session: err='%s'", err)
		w.stateMachine.SetState(WRAPPER_STATE_ONLINE)
		return err
	}

	w.stateMachine.SetState(WRAPPER_STATE_OFFLINE)

	// TODO: move to game session stop
	w.pushCmd("stop")
	<-time.After(5 * time.Second)

	w.console.kill()
	w.done <- true

	return nil
}

func (w *wrapper) processLogLine(line string) {
	if w.isOffline() {
		return
	}

	update, err := parseToLogUpdate(line)
	if err != nil {
		log.Printf("err parsing log: %s", err)
		return
	}

	log.Printf("Update detected: action='%s', target='%s', message='%s'", update.action, update.target, update.message)
	w.lastLogLine = update.message

	switch update.action {
	case SERVER_DONE:
		w.stateMachine.SetState(WRAPPER_STATE_ONLINE)
		go w.startGameSession()
	case SERVER_QUERY_TIME, SERVER_SAVED_THE_GAME:
		w.pushUpdateToGameSession(update)
	default:
	}
}

func (w *wrapper) processCmdOut() {
	for {
		line, err := w.console.readLine()
		if err != nil {
			if err == io.EOF {
				log.Printf("EOF reached! exiting log='%s'", line)
				return
			}
			log.Println(err)
			return
		}

		log.Printf("Raw log='%s'", line)
		w.processLogLine(line)
	}
}

func (w *wrapper) pushCmd(cmd string) error {
	log.Printf("pushing command=%s", cmd)
	return w.console.write(cmd)
}
