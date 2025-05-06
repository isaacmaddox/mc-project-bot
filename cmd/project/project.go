package project

import (
	"github.com/bwmarrin/discordgo"
)

var ProjectCommand = &discordgo.ApplicationCommand{
	Name:        "project",
	Description: "Interact with projects",
	Options: []*discordgo.ApplicationCommandOption{
		NewProjectCommand,
		OverviewCommand,
		ContributeCommand,
		ResourceCommand,
	},
}

var ProjectHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "new":
		NewProjectHandler(discord, i)
	case "overview":
		OverviewHandler(discord, i)
	case "contribute":
		ContributeHandler(discord, i)
	case "resource":
		ResourceHandler(discord, i)
	}
}
