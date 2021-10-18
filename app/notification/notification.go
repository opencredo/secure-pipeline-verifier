package notification

import "secure-pipeline-poc/app/config"


func Notify(cfg *config.Config, policyEvaluation []interface{})  {
   if !cfg.Notification.Enabled {
       return
   }
   var slack config.Slack
   if cfg.Notification.Slack != slack {
        NotifySlack(cfg.Notification.Slack.Channel, policyEvaluation)
   }
}