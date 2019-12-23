package mcs

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	GAME_TICKER_SYNC_INTERVAL time.Duration = 30 * time.Second

	MARKET_OPEN_GAMETICK int64 = 4000

	MARKET_CLOSE_GAMETICK int64 = 12000

	TIME_QUERY_CMD string = "time query daytime"
)

type gameSession struct {
	console   *console
	terminate chan bool
	updates   chan logUpdate
	gametick  int64

	isMarketOpen bool
}

func newGameSession(c *console) *gameSession {
	return &gameSession{
		console:      c,
		terminate:    make(chan bool),
		updates:      make(chan logUpdate),
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
			}
		case <-stop:
			return
		}
	}
}

func (s *gameSession) stop() {
	s.terminate <- true
}
