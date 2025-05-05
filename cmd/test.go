package cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/isaacmaddox/mc-project-bot/db"
	"github.com/isaacmaddox/mc-project-bot/util"
)

var TestCommand = &discordgo.ApplicationCommand{
	Name: "test",
	Description: "Used for testing only",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name: "project",
			Type: discordgo.ApplicationCommandOptionString,
			Description: "The name of the project",
			Required: true,
		},
	},
}

var TestHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	optionsMap := util.MakeOptionMap(i.ApplicationCommandData().Options)

	name := util.Extract(util.GetString(optionsMap, "project"))

	var project db.Project
	found := project.Get(name)

	log.Printf("%#v", project.Resources)
	
	if !found {
		discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("I couldn't find a project with the name %s", name),
			},
		})
		return
	}

	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("# %s\n*%s*", project.Name, project.Description),
		},
	})
}