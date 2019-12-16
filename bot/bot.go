package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/wlwanpan/minecraft-gobot/services"
	"google.golang.org/grpc"
)

var (
	ErrServerAlreadyRunning = errors.New("aws server already running")

	ErrServerAlreadyStopped = errors.New("aws serverr already stopped")
)

const (
	LAUNCHER_INSTANCE_ADDR        = "ec2-35-182-195-210.ca-central-1.compute.amazonaws.com:7777"
	THE_BRAIN_CHANNEL_ID   string = "654873630806769674"
)

type launcherClient struct {
	grpcConn *grpc.ClientConn
	client   services.LauncherServiceClient
}

func (c *launcherClient) initConn() error {
	conn, err := grpc.Dial(LAUNCHER_INSTANCE_ADDR, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.grpcConn = conn
	c.client = services.NewLauncherServiceClient(conn)
	return nil
}

func (c *launcherClient) closeConn() {
	if c.grpcConn == nil {
		log.Println("error closing conn: grpc conn must be open")
		return
	}
	if err := c.grpcConn.Close(); err != nil {
		log.Printf("error closing conn: %s", err)
		return
	}
}

type Bot struct {
	sync.Mutex
	sess           *discordgo.Session
	launcherClient *launcherClient
}

func New() *Bot {
	return &Bot{
		launcherClient: &launcherClient{},
	}
}

func (bot *Bot) Run() {
	sess, err := discordgo.New("Bot " + os.Getenv("DISCORD_KRAFFY_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.sess = sess

	// Add ws event handlers.
	bot.sess.AddHandler(bot.messageHandler)

	if err := bot.sess.Open(); err != nil {
		log.Fatal(err)
	}
	defer bot.Close()
	log.Println("Bot up and running...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func (bot *Bot) Close() {
	// Clean ups before closing bot.
	if bot.sess != nil {
		bot.sess.Close()
	}

	bot.launcherClient.closeConn()
}

func (bot *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		// Ignore bot own messages.
		return
	}
	if m.ChannelID != THE_BRAIN_CHANNEL_ID {
		// Ignore any other channel other than 'the-brain'
		return
	}

	log.Printf("Recv cmd: %s", m.Content)

	// Pipe out incoming commands.
	switch m.Content {
	case "start":
		bot.launchCmd(s, m)
	case "stop":
		bot.closeCmd(s, m)
	case "status":
		bot.statusCmd(s, m)
	default:
		log.Printf("Missing handler for command: %s", m.Content)
	}
}

func (bot *Bot) launchCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: need to build a connection manager to store/cache the active conns and
	// with ttl and handle conn dropout/timeout scenarios.
	if err := bot.launcherClient.initConn(); err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}
	defer bot.launcherClient.closeConn()

	c := bot.launcherClient.client
	config := &services.LaunchConfig{
		MemAlloc: 3,
	}
	resp, err := c.Launch(context.Background(), config)
	if err != nil {
		errMessage := fmt.Sprintf("Error: %s", err)
		sendMessageToChannel(s, m.ChannelID, errMessage)
		return
	}

	// Send init message: 'Launching server!'
	sendMessageToChannel(s, m.ChannelID, resp.GetMessage())

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			resp, err := c.Status(context.Background(), &services.EmptyReq{})
			if err != nil {
				log.Println(err)
				done <- true
			}

			if resp.GetMessage() == "RUNNING" {
				sendMessageToChannel(s, m.ChannelID, "Server up and running!")
				done <- true
			} else {
				sendMessageToChannel(s, m.ChannelID, "loading...")
			}
		}
	}
}

func (bot *Bot) closeCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := bot.launcherClient.initConn(); err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}
	defer bot.launcherClient.closeConn()

	resp, err := bot.launcherClient.client.Stop(context.Background(), &services.EmptyReq{})
	if err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}

	sendMessageToChannel(s, m.ChannelID, resp.GetMessage())
}

func (bot *Bot) statusCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := bot.launcherClient.initConn(); err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}
	defer bot.launcherClient.closeConn()

	status, err := bot.launcherClient.client.Status(context.Background(), &services.EmptyReq{})
	if err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}

	log.Printf("launcher status: %s", status.GetServerState())

	message := fmt.Sprintf("Server status: %s\n%s", status.GetServerState(), status.GetMessage())
	sendMessageToChannel(s, m.ChannelID, message)
}

func sendMessageToChannel(s *discordgo.Session, cid string, msg string) {
	wrappedMsg := fmt.Sprintf("```%s```", msg)
	s.ChannelMessageSend(cid, wrappedMsg)
}

// TODO: probably need to store/cache the aws client session
// instead of re initializing a session per op.
func getInstanceStatus() (*awsEC2StatusResp, error) {
	client, err := NewAwsClient()
	if err != nil {
		return nil, err
	}
	status, err := client.InstanceStatus()
	if err != nil {
		return nil, err
	}
	return status, nil
}

func startInstance() error {
	client, err := NewAwsClient()
	if err != nil {
		return err
	}
	status, err := client.InstanceStatus()
	if err != nil {
		if err == ErrNoRunningInstances {
			// All good start it.
		} else {
			return err
		}
	}
	if status.StateCode == INSTANCE_RUNNING_CODE {
		return ErrServerAlreadyRunning
	}
	if status.StateCode != INSTANCE_STOPPED_CODE {
		log.Printf("Can only start instance when its in a 'STOPPED' state: %s", status.State)
		return ErrServerAlreadyRunning
	}

	return client.StartInstance()
}

func stopInstance() error {
	client, err := NewAwsClient()
	if err != nil {
		return err
	}
	status, err := client.InstanceStatus()
	if err != nil {
		return err
	}
	if status.StateCode != INSTANCE_RUNNING_CODE {
		return ErrServerAlreadyStopped
	}

	return client.StopInstance()
}
