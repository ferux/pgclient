package bot

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ferux/pgclient/api"
	"github.com/sirupsen/logrus"
)

// redirect uri:
// https://discordapp.com/api/oauth2/authorize?client_id=485069755368079382&permissions=6208&redirect_uri=https%3A%2F%2Fdsbot.loyso.art%2Foauth&scope=bot

// Discord implements bot for discord
type Discord struct {
	token string

	dg     *discordgo.Session
	client *api.Client
	l      *logrus.Entry
}

// NewDiscordBot creates a new bot
func NewDiscordBot(token string, client *api.Client, ll logrus.Level) (*Discord, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	l := logrus.New()
	l.SetLevel(ll)
	le := l.WithFields(logrus.Fields{
		"pkg":  "bot",
		"kind": "discord",
	})

	return &Discord{
		token:  token,
		client: client,
		l:      le,
	}, nil
}

// Run is not async.
func (d *Discord) Run() error {
	var err error

	d.dg, err = discordgo.New("Bot " + d.token)
	if err != nil {
		return err
	}

	d.dg.AddHandler(d.messageCreate())
	d.dg.AddHandler(d.readyState())

	err = d.dg.Open()
	if err != nil {
		return fmt.Errorf("can't connect to discord: %v", err)
	}

	d.l.Info("bashstupid is running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	return d.dg.Close()
}

func (d *Discord) messageCreate() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !strings.HasPrefix(m.Content, "!stupidbash") {
			return
		}

		msg, err := d.client.GetMessage()
		if err != nil {
			d.l.WithError(err).Error("can't get new message from server")
			// obviously we shouldn't send error messages to chats.
			_, errs := s.ChannelMessageSend(m.ChannelID, err.Error())
			if errs != nil {
				d.l.WithError(err).Error("can't send error to discord server")
			}
		}
		_, err = s.ChannelMessageSend(m.ChannelID, msg.GetText())
		if err != nil {
			d.l.WithError(err).Error("can't send message to discrod server")
		}
	}
}

func (d *Discord) readyState() func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		if err := s.UpdateStatus(0, "!stupidbash"); err != nil {
			d.l.WithError(err).Error("can't update status")
		}
	}
}
