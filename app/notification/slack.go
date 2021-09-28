package notification

import (
    "github.com/slack-go/slack"
    "os"
    "strings"
)

const channel = "secure-pipeline"


func Notify(messages []string) {

    if messages == nil{
        return
    }

    token := os.Getenv("SLACK_TOKEN")
    client := slack.New(token)

    text := strings.Join(messages, "\n\n")

    opt := slack.MsgOptionText(text, true)
    _, _, _, err := client.SendMessage(channel, opt)
    if err != nil {
        panic("Slack couldn't send a message!")
    }
}
