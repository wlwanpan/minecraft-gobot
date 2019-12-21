package main

import (
	"log"
	"os"

	"github.com/wlwanpan/minecraft-gobot/config"

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

	// Minecraft server commands
	mcsCmd := p.NewCommand("mcs", "Start the minecraft server.")

	if err := p.Parse(os.Args); err != nil {
		log.Fatal(p.Usage(err))
	}

	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	switch {
	case botCmd.Happened():
		log.Println("Booting bot using config: ", config.Cfg.Bot)
		bot, err := bot.New(*token)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(bot.Run())
	case mcsCmd.Happened():
		log.Println("Booting bot using config: ", config.Cfg.Mcs)
		log.Fatal(mcs.Start(config.Cfg.Mcs.Port))
	}
}
