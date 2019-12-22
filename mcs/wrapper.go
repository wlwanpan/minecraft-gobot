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

	"github.com/wlwanpan/minecraft-gobot/config"
)

type wrapperState int

const (
	NEWLINE_BYTE byte = '\n'

	MEM_CONVERSION int = 1024

	// mcs initialized but no operation was performed yet, sleeping state.
	WRAPPER_STATE_OFFLINE wrapperState = iota
	// Minecraft server.jar successfully loaded and stdout 'DONE' is caught.
	WRAPPER_STATE_ONLINE
	// Minecraft serve.jar is still loading, has not yet caught 'Help' stdout.
	WRAPPER_STATE_LOADING
)

var (
	ErrServerAlreadyLoading = errors.New("server is already loading")

	ErrServerAlreadyOnline = errors.New("server is already online")

	ErrServerAlreadyOffline = errors.New("server is already offline")
)

var wrapperStateMap = map[wrapperState]string{
	WRAPPER_STATE_OFFLINE: "offline",
	WRAPPER_STATE_ONLINE:  "online",
	WRAPPER_STATE_LOADING: "loading",
}

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

func (c *console) kill() error {
	return c.execCmd.Process.Kill()
}

type sessionMetadata struct {
	connectedPlayers []string
}

func (meta *sessionMetadata) addConnectedPlayer(p string) {
	meta.connectedPlayers = append(meta.connectedPlayers, p)
}

func (meta *sessionMetadata) removeConnectedPlayer(p string) {
	n := []string{}
	for _, player := range meta.connectedPlayers {
		if player != p {
			n = append(n, player)
		}
	}
	meta.connectedPlayers = n
}

type wrapper struct {
	sync.RWMutex
	state   wrapperState
	console *console

	lastLogLine  string
	sessMetadata *sessionMetadata
}

func newWrapper() *wrapper {
	return &wrapper{
		state:   WRAPPER_STATE_OFFLINE,
		console: &console{},
	}
}

func (w *wrapper) isLoading() bool {
	w.RLock()
	defer w.RUnlock()
	return w.state == WRAPPER_STATE_LOADING
}

func (w *wrapper) isOnline() bool {
	w.RLock()
	defer w.RUnlock()
	return w.state == WRAPPER_STATE_ONLINE
}

func (w *wrapper) isOffline() bool {
	w.RLock()
	defer w.RUnlock()
	return w.state == WRAPPER_STATE_OFFLINE
}

func (w *wrapper) start(mem int) error {
	if w.isLoading() {
		return ErrServerAlreadyLoading
	}
	if w.isOnline() {
		return ErrServerAlreadyOnline
	}

	if err := w.console.execJava(mem); err != nil {
		return err
	}

	w.nextState(WRAPPER_STATE_LOADING)
	w.sessMetadata = &sessionMetadata{}

	go w.processCmdOut()

	return nil
}

func (w *wrapper) stop() error {
	if w.state != WRAPPER_STATE_ONLINE {
		return ErrServerAlreadyOffline
	}

	w.pushCmd("/stop\n") // Seems like stdin is not working ???
	<-time.After(3 * time.Second)

	if err := w.console.kill(); err != nil {
		return err
	}

	w.nextState(WRAPPER_STATE_OFFLINE)
	w.sessMetadata = nil

	return nil
}

func (w *wrapper) processUpdateSess(log string) {
	if w.isOffline() {
		return
	}

	update := parseToLogUpdate(log)
	if update == nil {
		return
	}
	switch update.action {
	case PLAYER_JOINED:
		w.sessMetadata.addConnectedPlayer(update.target)
		return
	case PLAYER_LEFT:
		w.sessMetadata.removeConnectedPlayer(update.target)
		return
	default:
		w.lastLogLine = update.message
	}
}

func (w *wrapper) sessionSummary() string {
	if w.isOffline() {
		return ""
	}
	if len(w.sessMetadata.connectedPlayers) == 0 {
		return w.lastLogLine
	}

	return fmt.Sprintf("Players online: %s", strings.Join(w.sessMetadata.connectedPlayers, ", "))
}

func (w *wrapper) processCmdOut() {
	for {
		line, err := w.console.cmdout.ReadString(NEWLINE_BYTE)
		if err != nil {
			if err == io.EOF {
				log.Printf("EOF reached: %s", line)
				break
			}
			log.Println(err)
		}

		log.Println(line)

		// TODO: Move the "Done" cond here to log_update
		if strings.Contains(line, "Done") {
			w.nextState(WRAPPER_STATE_ONLINE)
			if w.sessMetadata == nil {
				w.sessMetadata = &sessionMetadata{}
			}
			continue
		}

		w.processUpdateSess(line)
	}
}

func (w *wrapper) pushCmd(cmd string) error {
	n, err := w.console.cmdin.WriteString(cmd)
	if err != nil {
		return err
	}
	if n != len(cmd) {
		log.Println("Error pushing command to stdin")
		return nil
	}

	log.Println("Successfully pushed.")
	return nil
}

func (w *wrapper) nextState(s wrapperState) {
	w.Lock() // Should probably be a read/write lock instead
	defer w.Unlock()

	from := wrapperStateMap[w.state]
	to := wrapperStateMap[s]

	switch w.state {
	case WRAPPER_STATE_OFFLINE:
		if s != WRAPPER_STATE_LOADING {
			log.Printf("Invalid transition: %s -> %s", from, to)
			return
		}
	case WRAPPER_STATE_ONLINE:
		if s != WRAPPER_STATE_OFFLINE {
			log.Printf("Invalid transition: %s -> %s", from, to)
			return
		}
	case WRAPPER_STATE_LOADING:
		// all good!
	default:
		log.Fatalf("Current state: %s not handled", from)
	}

	log.Printf("State transition: %s -> %s", from, to)
	w.state = s
}
