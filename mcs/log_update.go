package mcs

import (
	"regexp"
	"strings"
)

const (
	PLAYER_JOINED string = "joined"
	PLAYER_LEFT   string = "left"

	SERVER_QUERY_TIME      string = "time is"
	SERVER_PREPARING_SPAWN string = "Preparing spawn"
	SERVER_PREPARING_LEVEL string = "Preparing level"
)

var actionables = map[string]*regexp.Regexp{
	PLAYER_JOINED:     regexp.MustCompile(`]: (?s)(.*) joined the game`),
	PLAYER_LEFT:       regexp.MustCompile(`]: (?s)(.*) left the game`),
	SERVER_QUERY_TIME: regexp.MustCompile(`]: The time is (?s)(.*)\n`),
}

var defaultActionables = regexp.MustCompile(`/INFO]: (?s)(.*)\n`)

type logUpdate struct {
	action  string
	target  string
	message string
}

func parseToLogUpdate(l string) *logUpdate {
	for action, reg := range actionables {
		if !strings.Contains(l, action) {
			continue
		}
		return &logUpdate{
			action: action,
			target: reg.FindStringSubmatch(l)[1],
		}
	}

	r := defaultActionables.FindStringSubmatch(l)
	if len(r) < 2 {
		return nil
	}

	return &logUpdate{
		message: r[1],
	}
}
