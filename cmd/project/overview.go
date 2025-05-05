package project

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/isaacmaddox/mc-project-bot/db"
	"github.com/isaacmaddox/mc-project-bot/util"
)

var OverviewCommand = &discordgo.ApplicationCommandOption{
	Name:        "overview",
	Description: "See progress so far!",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "name",
			Description: "Name",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

var OverviewHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := util.MakeOptionMap(i.ApplicationCommandData().Options[0].Options)

	name := util.Extract(util.GetString(optionMap, "name"))

	var project db.Project
	project.Get(name)

	var resourceList string
	var totalGoal int
	var totalAmount int

	resources := project.Resources
	for _, resource := range resources {
		totalAmount += resource.Amount
		totalGoal += resource.Goal

		resourceList += fmt.Sprintf(
			"%s\n-# [%s] (%s of %s)\n\n",
			resource.Name,
			util.Make_progress(resource.Amount, resource.Goal, 10),
			util.To_unit(resource.Amount),
			util.To_unit(resource.Goal),
		)
	}

	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"# %s\n[%s] (%d%%)\n\n%s",
				project.Name,
				util.Make_progress(totalAmount, totalGoal, 15),
				int(float32(totalAmount)/float32(totalGoal)*100),
				resourceList,
			),
		},
	})
}
