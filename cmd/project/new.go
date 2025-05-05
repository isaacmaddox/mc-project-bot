package project

import (
	"github.com/bwmarrin/discordgo"
	"github.com/isaacmaddox/mc-project-bot/db"
	"github.com/isaacmaddox/mc-project-bot/util"
)

var NewProjectCommand = &discordgo.ApplicationCommandOption{
	Name:        "new",
	Description: "Create a new project",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "name",
			Description: "Name",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
		{
			Name:        "description",
			Description: "Description",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

var NewProjectHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := util.MakeOptionMap(i.ApplicationCommandData().Options[0].Options)

	name := util.Extract(util.GetString(optionMap, "name"))
	description := util.Extract(util.GetString(optionMap, "description"))

	var project db.Project
	project.Create(name, description)

	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "New project created:",
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       project.Name,
					Description: project.Description,
				},
			},
		},
	})
}
