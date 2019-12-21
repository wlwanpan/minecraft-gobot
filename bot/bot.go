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

	ErrServerAlreadyStopped = errors.New("aws server already stopped")
)

const (
	MCS_PROD_INSTANCE_ADDR string = "ec2-35-182-195-210.ca-central-1.compute.amazonaws.com:7777"
	MCS_DEV_INSTANCE_ADDR  string = ":7777"
	THE_BRAIN_CHANNEL_ID   string = "654873630806769674"
)

type mcsClient struct {
	grpcConn *grpc.ClientConn
	client   services.LauncherServiceClient
}

func (c *mcsClient) Start(ctx context.Context, in *services.StartConfig) (*services.ServiceResp, error) {
	return c.client.Start(ctx, in)
}

func (c *mcsClient) Stop(ctx context.Context) (*services.ServiceResp, error) {
	return c.client.Stop(ctx, &services.EmptyReq{})
}

func (c *mcsClient) Status(ctx context.Context) (*services.StatusResp, error) {
	return c.client.Status(ctx, &services.EmptyReq{})
}

func (c *mcsClient) initConn() error {
	addr := MCS_DEV_INSTANCE_ADDR
	if os.Getenv("IS_DEV") == "0" {
		addr = MCS_PROD_INSTANCE_ADDR
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.grpcConn = conn
	c.client = services.NewLauncherServiceClient(conn)
	return nil
}

func (c *mcsClient) closeConn() {
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
	if err := bot.mcsClient.initConn(); err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}
	defer bot.mcsClient.closeConn()

	config := &services.StartConfig{
		MemAlloc: 3,
	}
	resp, err := bot.mcsClient.Start(context.Background(), config)
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
			resp, err := bot.mcsClient.Status(context.Background())
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
	if err := bot.mcsClient.initConn(); err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}
	defer bot.mcsClient.closeConn()

	resp, err := bot.mcsClient.Stop(context.Background())
	if err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}

	sendMessageToChannel(s, m.ChannelID, resp.GetMessage())
}

func (bot *Bot) statusCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := bot.mcsClient.initConn(); err != nil {
		sendMessageToChannel(s, m.ChannelID, err.Error())
		return
	}
	defer bot.mcsClient.closeConn()

	status, err := bot.mcsClient.Status(context.Background())
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
