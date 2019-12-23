package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/wlwanpan/minecraft-gobot/config"
	"github.com/wlwanpan/minecraft-gobot/services"
)

type Bot struct {
	sync.Mutex
	sess      *discordgo.Session
	mcsClient *mcsClient
}

func New(token string) (*Bot, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		sess:      sess,
		mcsClient: &mcsClient{},
	}

	// Add discord ws message handlers.
	b.sess.AddHandler(b.messageHandler)
	go b.mcsClient.initConn()

	return b, nil
}

func (bot *Bot) Run() error {
	if err := bot.sess.Open(); err != nil {
		return err
	}

	defer bot.Close()
	log.Println("Bot up and running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	return nil
}

func (bot *Bot) Close() {
	if bot.sess != nil {
		bot.sess.Close()
	}
	if bot.mcsClient != nil {
		bot.mcsClient.closeConn()
	}
}

func (bot *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		// Ignore bot own messages.
		return
	}

	if !isValidChannelID(m.ChannelID) {
		// Ignore any other channel not whitelisted in the config file.
		log.Printf("ignoring message=%s, from=%s", m.Content, m.ChannelID)
		return
	}

	log.Printf("Receiving command=%s, from=%s", m.Content, s.State.User.Username)

	// Probably need to be a cancellable context
	ctx := context.Background()

	// Pipe out incoming commands.
	switch m.Content {
	case "start":
		bot.startCmd(ctx, s, m)
	case "stop":
		bot.closeCmd(ctx, s, m)
	case "status":
		bot.statusCmd(ctx, s, m)
	default:
		log.Printf("Unknown command=%s", m.Content)
	}
}

func (bot *Bot) startCmd(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := bot.mcsClient.checkConn(ctx); err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}

	config := &services.StartConfig{
		MemAlloc: 3,
	}
	resp, err := bot.mcsClient.Start(ctx, config)
	if err != nil {
		errMessage := fmt.Sprintf("Error: %s", err)
		sendMessageToChannel(s, m.ChannelID, errMessage)
		return
	}

	sendMessageToChannel(s, m.ChannelID, resp.GetMessage())

	ticker := time.NewTicker(3 * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			resp, err := bot.mcsClient.Status(ctx)
			if err != nil {
				log.Println(err)
				done <- true
			}

			switch resp.GetServerState() {
			case "online":
				sendMessageToChannel(s, m.ChannelID, "Server up and running!")
				done <- true
			case "loading":
				sendMessageToChannel(s, m.ChannelID, resp.GetMessage())
			default:
				message := fmt.Sprintf("Error, server state: %s", resp.GetServerState())
				sendMessageToChannel(s, m.ChannelID, message)
				done <- true
			}
		}
	}
}

func (bot *Bot) closeCmd(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := bot.mcsClient.checkConn(ctx); err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}

	resp, err := bot.mcsClient.Stop(ctx)
	if err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}

	sendMessageToChannel(s, m.ChannelID, resp.GetMessage())
}

func (bot *Bot) statusCmd(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := bot.mcsClient.checkConn(ctx); err != nil {
		message := fmt.Sprintf("Error connecting to minecraft server: %s", err)
		if err == ErrMcsNotResponding {
			message = "Server status: offline"
		}
		sendMessageToChannel(s, m.ChannelID, message)
		return
	}

	status, err := bot.mcsClient.Status(ctx)
	if err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}

	log.Printf("mcs status: %s", status.GetServerState())

	message := fmt.Sprintf("Server status: %s\n", status.GetServerState())
	sendMessageToChannel(s, m.ChannelID, message)
}

func sendMessageToChannel(s *discordgo.Session, cid string, msg string) {
	wrappedMsg := fmt.Sprintf("```%s```", msg)
	s.ChannelMessageSend(cid, wrappedMsg)
}

func isValidChannelID(cid string) bool {
	for _, chanID := range config.Cfg.Bot.WhitelistedChannelIDS {
		if cid == chanID {
			return true
		}
	}
	return false
}
