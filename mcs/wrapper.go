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

const (
	// Game event specific const
	MARKET_OPEN_GAMETICK int64 = 4000

	MARKET_CLOSE_GAMETICK int64 = 12000
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

func (w *wrapper) startScheduler() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Check game tick related events.
			w.pushCmd("time query daytime")
		case <-w.done:
			return
		}
	}
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
	go w.startScheduler()

	return nil
}

func (w *wrapper) stop() error {
	if w.state != WRAPPER_STATE_ONLINE {
		return ErrServerAlreadyOffline
	}

	w.pushCmd("stop")

	// Guarentee that the process is killed after 5s delay.
	<-time.After(3 * time.Second)
	if err := w.console.kill(); err != nil {
		return err
	}

	w.sessMetadata = nil
	w.nextState(WRAPPER_STATE_OFFLINE)
	w.done <- true

	return nil
}

func (w *wrapper) processUpdateSess(l string) {
	if w.isOffline() {
		return
	}

	update := parseToLogUpdate(l)
	if update == nil {
		return
	}
	log.Printf("Update detected: action='%s', target='%s', message='%s'", update.action, update.target, update.message)

	switch update.action {
	case PLAYER_JOINED:
		w.sessMetadata.addConnectedPlayer(update.target)
		return
	case PLAYER_LEFT:
		w.sessMetadata.removeConnectedPlayer(update.target)
		return
	case SERVER_QUERY_TIME:
		gametick, _ := strconv.ParseInt(update.target, 10, 64)
		if gametick >= MARKET_OPEN_GAMETICK && gametick <= MARKET_CLOSE_GAMETICK {
			if !w.sessMetadata.isMarketOpen {
				c := fmt.Sprintf("say %s", "Market is now open!")
				w.pushCmd(c)
				w.sessMetadata.isMarketOpen = true
			}
		} else {
			if w.sessMetadata.isMarketOpen {
				c := fmt.Sprintf("say %s", "Market is now closed!")
				w.pushCmd(c)
				w.sessMetadata.isMarketOpen = false
			}
		}
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

		log.Printf("Raw log line: %s\n", line)

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
	log.Printf("pushing command=%s", cmd)
	return w.console.write(cmd)
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
