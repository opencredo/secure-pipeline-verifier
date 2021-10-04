package notification

import (
	"github.com/slack-go/slack"
	"os"
	"secure-pipeline-poc/app/config"
	"strings"
)

const channel = "secure-pipeline"

func Notify(messages []string){
	if len(messages) < 1{
		return
	}
	token := os.Getenv(config.SlackToken)
	client := slack.New(token)

	text := strings.Join(messages, "\n\n")

	opt := slack.MsgOptionText(text, true)
	_, _, _, err := client.SendMessage(channel, opt)
	if err != nil {
		panic("Slack couldn't send a message!")
	}
}
