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

type sessionMetadata struct {
	connectedPlayers []string

	// TODO: move to new game session state.
	isMarketOpen bool
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
	done    chan bool

	gameSess    *gameSession
	lastLogLine string
}

func newWrapper() *wrapper {
	return &wrapper{
		state:   WRAPPER_STATE_OFFLINE,
		console: &console{},
		done:    make(chan bool),
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

func (w *wrapper) startGameSession() {
	if w.gameSess != nil {
		w.gameSess.stop()
	}

	log.Println("Starting new game session")
	w.gameSess = newGameSession(w.console)
	go w.gameSess.start()
}

func (w *wrapper) stopGameSession() {
	if w.gameSess == nil {
		return
	}

	log.Printf("Stopping current game session, game-tick=%d", w.gameSess.gametick)
	w.gameSess.stop()
	w.gameSess = nil
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

	go w.processCmdOut() // exits on EOF

	return nil
}

func (w *wrapper) stop() error {
	if w.state != WRAPPER_STATE_ONLINE {
		return ErrServerAlreadyOffline
	}

	// Guarentee that the process is killed after 5s delay.
	defer w.console.kill()

	// Dont like this but works for now.
	w.pushCmd("save-all flush")
	<-time.After(5 * time.Second)
	w.pushCmd("stop")
	<-time.After(5 * time.Second)

	w.nextState(WRAPPER_STATE_OFFLINE)
	w.stopGameSession()
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

	if w.gameSess != nil {
		w.gameSess.updates <- update
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

		// TODO: Move the "Done" cond here to log_update
		if strings.Contains(line, "Done") {
			w.nextState(WRAPPER_STATE_ONLINE)
			w.startGameSession()
			continue
		}

		w.processLogLine(line)
	}
}

func (w *wrapper) pushCmd(cmd string) error {
	log.Printf("pushing command=%s", cmd)
	return w.console.write(cmd)
}

func (w *wrapper) nextState(s wrapperState) {
	w.Lock()
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
