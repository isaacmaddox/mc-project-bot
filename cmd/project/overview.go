package project

import (
	"fmt"
	"log"

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
			Name:         "name",
			Description:  "Name",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Name:        "page",
			Description: "Page number",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Required:    false,
		},
	},
}

var OverviewHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		optionMap := util.MakeOptionMap(i.ApplicationCommandData().Options[0].Options)

		name := util.Extract(util.GetString(optionMap, "name"))
		page := max(util.GetIntOptional(optionMap, "page", 1), 1)

		OverviewPaginatedHandler(discord, i, name, page)
	} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		options := []*discordgo.ApplicationCommandOptionChoice{}

		for _, name := range db.GetProjectNames() {
			options = append(options, &discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: name,
			})
		}

		err := discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: options,
			},
		})

		if err != nil {
			log.Fatalf("Error getting autocomplete: %v", err)
		}
	}
}

var OverviewPaginatedHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate, name string, page int64) {
	var project db.Project
	project.Get(name)

	resourceList := []string{""}
	var projectAmount int
	var projectGoal int

	var areMorePages bool = false

	resources := project.Resources

	for _, resource := range resources {
		projectAmount += min(resource.Amount, resource.Goal)
		projectGoal += resource.Goal

		percentage := min(float32(resource.Amount)/float32(resource.Goal)*100, 100)
		progressBar := util.MakeProgress(resource.Amount, resource.Goal, 15)

		resourceString := fmt.Sprintf("%s (%.0f%%)\n-# [%s] (%s of %s)\n\n",
			resource.Name,
			percentage,
			progressBar,
			util.ToUnit(resource.Amount),
			util.ToUnit(resource.Goal),
		)

		var lastItem *string = &resourceList[len(resourceList)-1]

		if len(*lastItem)+len(resourceString) <= 1000 {
			*lastItem += resourceString
		} else {
			if int64(len(resourceList)) == page {
				areMorePages = true
				break
			}
			resourceList = append(resourceList, resourceString)
		}
	}

	progressBar := util.MakeProgress(projectAmount, projectGoal, 25)
	progressPercent := (float32(projectAmount) / float32(projectGoal)) * 100

	header := fmt.Sprintf("# %s\n[%s] (%.0f%%)\n\n",
		project.Name,
		progressBar,
		progressPercent,
	)

	var err error
	var messageType discordgo.InteractionResponseType

	if i.Type == discordgo.InteractionMessageComponent {
		messageType = discordgo.InteractionResponseUpdateMessage
	} else {
		messageType = discordgo.InteractionResponseChannelMessageWithSource
	}

	err = discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: messageType,
		Data: &discordgo.InteractionResponseData{
			Content: header + resourceList[page-1],
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Previous",
							Style:    discordgo.SecondaryButton,
							CustomID: fmt.Sprintf("project_overview_%d_%s", page-1, name),
							Disabled: page == 1,
						},
						discordgo.Button{
							Label:    "Next",
							Style:    discordgo.SecondaryButton,
							CustomID: fmt.Sprintf("project_overview_%d_%s", page+1, name),
							Disabled: !areMorePages,
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Fatalf("could not respond to interaction: %v", err)
	}
}
