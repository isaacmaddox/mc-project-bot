package project

import (
	"github.com/bwmarrin/discordgo"
	"github.com/isaacmaddox/mc-project-bot/cmd/project/resource"
)

var ResourceCommand = &discordgo.ApplicationCommandOption{
	Name: "resource",
	Description: "Interact with resources",
	Type: discordgo.ApplicationCommandOptionSubCommandGroup,
	Options: []*discordgo.ApplicationCommandOption{
		resource.NewResourceCommand,
	},
}

var ResourceHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options

		switch options[0].Name {
		case "new":
			resource.NewResourceHandler(discord, i)
		}
}