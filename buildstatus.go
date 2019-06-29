package buildstatus

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/nlopes/slack"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/pubsub/v1"
)

func NotifyBuildStatus(ctx context.Context, m *pubsub.PubsubMessage) error {
	decoded, err := base64.StdEncoding.DecodeString(m.Data)
	if err != nil {
		return err
	}
	build := cloudbuild.Build{}
	err = json.Unmarshal(decoded, &build)
	if err != nil {
		return err
	}

	notifyStatus := map[string]bool{
		"STATUS_UNKNOWN": false,
		"QUEUED":         false,
		"WORKING":        false,
		"SUCCESS":        true,
		"FAILURE":        true,
		"INTERNAL_ERROR": true,
		"TIMEOUT":        true,
		"CANCELLED":      false,
	}

	if notifyStatus[build.Status] {
		err = pushSlack(build.Status)
		if err != nil {
			return err
		}
	}

	return nil
}

func pushSlack(msg string) error {
	token := ""
	api := slack.New(token)
	attachment := slack.Attachment{
		Pretext: msg,
	}
	cannelID := ""
	title := ""
	_, _, err := api.PostMessage(cannelID, slack.MsgOptionText(title, false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		return err
	}

	return nil
}
