package cmd

import "github.com/bwmarrin/discordgo"

var HiCommand = &discordgo.ApplicationCommand{
	Name: "hi",
	Description: "Says hi back, like a gentleman",
}

var HiHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello there",
		},
	})
}