package mcs

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

var (
	ErrSavingGameTimedOut = errors.New("saving game timed out")
)

const (
	GAME_TICKER_SYNC_INTERVAL time.Duration = 30 * time.Second

	GAME_SESSION_SAVE_TIMEOUT time.Duration = 5 * time.Second

	MARKET_OPEN_GAMETICK int64 = 2000

	MARKET_CLOSE_GAMETICK int64 = 9000

	// Minecraft server commands
	TIME_QUERY_CMD string = "time query daytime"

	SAVE_GAME_CMD string = "save-all flush"
)

type gameSession struct {
	console   *console
	terminate chan bool
	updates   chan logUpdate
	saves     chan bool
	gametick  int64

	isMarketOpen bool
}

func newGameSession(c *console) *gameSession {
	return &gameSession{
		console:      c,
		terminate:    make(chan bool),
		updates:      make(chan logUpdate),
		saves:        make(chan bool),
		gametick:     0,
		isMarketOpen: false,
	}
}

func (s *gameSession) start() {
	stop := make(chan bool, 3)

	go s.processUpdates(stop)

	gameTicker := time.NewTicker(1 * time.Second)
	defer gameTicker.Stop()
	go func() {
		for {
			select {
			case <-gameTicker.C:
				// Minecraft game tick runs at a fixed rate of 20 ticks per second.
				// reference: https://minecraft.gamepedia.com/Tick
				s.gametick += 20
				s.checkGameTickEvents()
			case <-stop:
				return
			}
		}
	}()

	s.console.write(TIME_QUERY_CMD)
	queryTicker := time.NewTicker(GAME_TICKER_SYNC_INTERVAL)
	defer queryTicker.Stop()
	go func() {
		for {
			select {
			case <-queryTicker.C:
				s.console.write(TIME_QUERY_CMD)
			case <-stop:
				return
			}
		}
	}()

	<-s.terminate
	stop <- true
	stop <- true
	stop <- true
}

func (s *gameSession) save() error {
	s.console.write(SAVE_GAME_CMD)
	select {
	case <-s.saves:
		return nil
	case <-time.After(GAME_SESSION_SAVE_TIMEOUT):
		return ErrSavingGameTimedOut
	}
}

func (s *gameSession) checkGameTickEvents() {
	if s.gametick >= MARKET_OPEN_GAMETICK && s.gametick <= MARKET_CLOSE_GAMETICK {
		if !s.isMarketOpen {
			s.isMarketOpen = true
			c := fmt.Sprintf("say %s", "Market is now open!")
			s.console.write(c)
		}
	} else {
		if s.isMarketOpen {
			s.isMarketOpen = false
			c := fmt.Sprintf("say %s", "Market is now closed!")
			s.console.write(c)
		}
	}
}

func (s *gameSession) processUpdates(stop chan bool) {
	for {
		select {
		case update := <-s.updates:
			switch update.action {
			case SERVER_QUERY_TIME:
				realGametick, _ := strconv.ParseInt(update.target, 10, 64)
				log.Printf("Syncing game ticks, current-tick=%d, real-tick=%d", s.gametick, realGametick)

				s.gametick = realGametick
			case SERVER_SAVED_THE_GAME:
				s.saves <- true
			}
		case <-stop:
			fmt.Println("Stopping processUpdates")
			return
		}
	}
}

func (s *gameSession) stop() {
	s.terminate <- true
}
