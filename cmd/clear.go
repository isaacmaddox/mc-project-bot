package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/isaacmaddox/mc-project-bot/util"
)

var ClearCommand = &discordgo.ApplicationCommand{
	Name: "clear",
	Description: "Clear the current channel (use only in testing)",
}

var ClearHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Clearing... channel",
		},
	})
	
	for ok := true; ok; ok = true {
		messages := util.Extract(discord.ChannelMessages(i.ChannelID, 50, "", "", ""))
		if len(messages) == 0 { break }
		messageIds := []string{}

		for _, message := range messages {
			messageIds = append(messageIds, message.ID)
		}

		util.ErrorCheck(discord.ChannelMessagesBulkDelete(i.ChannelID, messageIds), "Error deleting messages: %v")
	}

	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Cleared channel",
		},
	})
}