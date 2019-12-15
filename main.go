package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/wlwanpan/minecraft-gobot/bot"
	"github.com/wlwanpan/minecraft-gobot/launcher"
)

func main() {
	p := argparse.NewParser("minecraft-bot", "Discord bot to manage a minecraft server hosted on aws.")

	// Bot commands
	botCmd := p.NewCommand("bot", "Start the discord bot.")

	// Launcher commands
	launcherCmd := p.NewCommand("launcher", "Start the launcher server.")

	if err := p.Parse(os.Args); err != nil {
		log.Fatal(p.Usage(err))
	}

	switch {
	case botCmd.Happened():
		bot := bot.New()
		bot.Run()
	case launcherCmd.Happened():
		if err := launcher.Start(); err != nil {
			log.Fatal(err)
		}
	}
}
