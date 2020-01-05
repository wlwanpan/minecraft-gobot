package mcs

import (
	"errors"
	"regexp"
	"strings"
)

const (
	PLAYER_JOINED string = "joined"
	PLAYER_LEFT   string = "left"

	SERVER_QUERY_TIME      string = "time is"
	SERVER_PREPARING_SPAWN string = "Preparing spawn"
	SERVER_PREPARING_LEVEL string = "Preparing level"
	SERVER_SAVED_THE_GAME  string = "Saved the game"
)

var (
	ErrNoActionableLogUpdate = errors.New("no actionables in the log line")
)

var actionables = map[string]*regexp.Regexp{
	PLAYER_JOINED:     regexp.MustCompile(`]: (?s)(.*) joined the game`),
	PLAYER_LEFT:       regexp.MustCompile(`]: (?s)(.*) left the game`),
	SERVER_QUERY_TIME: regexp.MustCompile(`]: The time is (?s)(.*)\n`),
}

// TODO: move the 'Done' log here < for startCmd.
var gameStateChanges = []string{
	SERVER_SAVED_THE_GAME,
}

var defaultActionables = regexp.MustCompile(`/INFO]: (?s)(.*)\n`)

type logUpdate struct {
	action  string
	target  string
	message string
}

func parseToLogUpdate(l string) (logUpdate, error) {
	// Try and parse the raw log as an actionable.
	for action, reg := range actionables {
		if !strings.Contains(l, action) {
			continue
		}
		return logUpdate{
			action: action,
			target: reg.FindStringSubmatch(l)[1],
		}, nil
	}

	// Try and parse the raw log as an gameStateChange.
	for _, stateChange := range gameStateChanges {
		if !strings.Contains(l, stateChange) {
			continue
		}
		return logUpdate{action: SERVER_SAVED_THE_GAME}, nil
	}

	r := defaultActionables.FindStringSubmatch(l)
	if len(r) < 2 {
		return logUpdate{}, ErrNoActionableLogUpdate
	}

	return logUpdate{
		message: r[1],
	}, nil
}
