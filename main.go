package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/wlwanpan/minecraft-gobot/bot"
	"github.com/wlwanpan/minecraft-gobot/mcs"
)

func main() {
	p := argparse.NewParser("minecraft-bot", "Discord bot to manage a minecraft server hosted on aws.")

	// Bot commands
	botCmd := p.NewCommand("bot", "Start the discord bot.")
	token := botCmd.String("t", "token", &argparse.Options{
		Required: false,
		Default:  os.Getenv("DISCORD_KRAFFY_TOKEN"),
		Help: `A discord token to auth the bot else it will read
			   from the env variable: 'DISCORD_KRAFFY_TOKEN'`,
	})

	// Launcher commands
	mcsCmd := p.NewCommand("mcs", "Start the minecraft server.")

	// Common flags
	dev := p.Flag("d", "dev", &argparse.Options{
		Required: false,
		Default:  true,
	})

	if err := p.Parse(os.Args); err != nil {
		log.Fatal(p.Usage(err))
	}

	if *dev {
		os.Setenv("IS_DEV", "1")
	} else {
		os.Setenv("IS_DEV", "0")
	}

	switch {
	case botCmd.Happened():
		bot, err := bot.New(*token)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(bot.Run())
	case mcsCmd.Happened():
		log.Fatal(mcs.Start())
	}
}
