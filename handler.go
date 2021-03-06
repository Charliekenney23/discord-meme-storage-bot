package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var handleGetMemeMessage = func(sess *discordgo.Session, evt *discordgo.MessageCreate) {
	msg := evt.Message
	contentParts := strings.Fields(strings.TrimSpace(strings.ToLower(msg.Content)))

	if !(len(contentParts) >= 1) || contentParts[0] != "!meme" {
		return
	}

	if len(msg.Attachments) > 0 {
		return
	}

	chanID := msg.ChannelID

	if !(len(contentParts) <= 2) {
		log.WithFields(logrus.Fields{
			"messageId": msg.ID,
			"content":   msg.Content,
		}).Error("Could not find meme name to get")

		res := "God damnit! I couldn't figure out what the name of the meme you " +
			"wanted me to get! When fetching memes from storage, format your messages" +
			"like so `!meme <meme_name>`. (:"
		memeBot.sendTextMessage(chanID, res)
		return
	}

	name := contentParts[1]
	image, err := storage.Get(name)
	if err != nil {
		memeBot.sendTextMessage(chanID, "God damnit! I couldn't get that meme for some reason.")
		return
	}

	memeBot.sendFileMessage(chanID, image)
}

var handlePutMemeMessage = func(sess *discordgo.Session, evt *discordgo.MessageCreate) {
	msg := evt.Message
	contentParts := strings.Fields(strings.TrimSpace(strings.ToLower(msg.Content)))

	if !(len(contentParts) >= 1) || contentParts[0] != "!meme" {
		return
	}

	if !(len(msg.Attachments) >= 1) {
		return
	}

	chanID := msg.ChannelID

	if !(len(contentParts) <= 2) {
		log.WithFields(logrus.Fields{
			"messageId": msg.ID,
			"content":   msg.Content,
		}).Error("Could not find meme name to save")

		res := "God damnit! I couldn't figure out what name you wanted me to " +
			"save this meme under! When sending an image in, format the message " +
			"like so `!meme <meme_name>`. (:"
		memeBot.sendTextMessage(chanID, res)
		return
	}

	name := contentParts[1]

	if msg.Attachments != nil && len(msg.Attachments) != 0 {
		uri := msg.Attachments[0].URL
		err := storage.Put(name, uri)
		if err != nil {
			log.WithFields(logrus.Fields{
				"messageID": msg.ID,
				"uri":       uri,
				"name":      name,
			}).Error("Could not save image")
			memeBot.sendTextMessage(chanID, "God damnit! I couldn't save that for some reason. ):")
			return
		}

		memeBot.sendTextMessage(chanID, fmt.Sprintf("Cool! I saved your meme as `%s`!", name))
	} else {
		log.WithFields(logrus.Fields{
			"messageId": msg.ID,
		}).Error("No imaged attached")

		memeBot.sendTextMessage(chanID, "God damnit! I couldn't find a meme to save on that message. ):")
		return
	}
}
