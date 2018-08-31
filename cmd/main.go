package main

import (
	"os"

	"github.com/ferux/pgclient"
	"github.com/ferux/pgclient/api"
	"github.com/ferux/pgclient/bot"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	l *logrus.Entry
)

func init() {
	l = logrus.New().WithField("pkg", "pgclient")
	if err := godotenv.Load(); err != nil {
		l.WithError(err).Fatal("can't load env variables")
	}
	if pgclient.ConnString = os.Getenv("GO_GRPC_SERVER"); pgclient.ConnString == "" {
		l.Fatal("connection String is not set")
	}
	// todo: it will be reworked lately
	if pgclient.DiscordBotToken = os.Getenv("GO_DISCORD_BOT"); pgclient.DiscordBotToken == "" {
		l.Fatal("token is not set")
	}
}

func main() {
	l.WithFields(logrus.Fields{
		"ver": pgclient.Version,
		"rev": pgclient.Revision,
		"env": pgclient.Environment,
	}).Info("client started")
	c := api.NewClient(pgclient.ConnString)
	if err := c.Run(); err != nil {
		l.WithError(err).Fatal("can't init connection")
	}
	defer func() {
		if err := c.Close(); err != nil {
			l.WithError(err).Warn("can't close connection to gRPC")
		}
		l.Info("Finishing operation")
	}()
	// for i := 0; i < 5; i++ {
	// 	msg, err := c.GetMessage()
	// 	if err != nil {
	// 		l.WithError(err).Fatal("can't get message")
	// 	}
	// 	l.WithTime(time.Now()).Infof("[%s] got message: %s", msg.GetId(), msg.GetText())
	// 	time.Sleep(time.Second * 1)
	// }
	// st, err := c.AskStatus()
	// if err != nil {
	// 	l.WithError(err).Fatal("can't get message")
	// }
	//l.WithTime(time.Now()).Infof("Server Status: %s", st.GetStatus())

	dbot, err := bot.NewDiscordBot(pgclient.DiscordBotToken, c, logrus.InfoLevel)
	if err != nil {
		l.WithError(err).Fatal("can't run dcbot")
	}

	l.WithError(dbot.Run()).Info("dsbot finished")
}
