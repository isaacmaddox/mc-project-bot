package project

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/isaacmaddox/mc-project-bot/db"
	"github.com/isaacmaddox/mc-project-bot/util"
)

var ContributeCommand = &discordgo.ApplicationCommandOption{
	Name:        "contribute",
	Description: "Contribute to a specific resource",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "project",
			Description:  "Which project to contribute to",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Name:         "resource",
			Description:  "Whatchu givin fool?",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Name:        "amount",
			Description: "How much are ye givin?",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

var ContributeHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		optionMap := util.MakeOptionMap(i.ApplicationCommandData().Options[0].Options)

		projectName := util.Extract(util.GetString(optionMap, "project"))
		resourceName := util.Extract(util.GetString(optionMap, "resource"))
		amount := util.FromUnit(util.Extract(util.GetString(optionMap, "amount")))

		if amount < 1 {
			err := discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Cannot contribute less than one item",
				},
			})

			if err != nil {
				log.Fatalf("could not respond to interaction: %v", err)
			}
		} else {
			var project db.Project
			project.Get(projectName)

			retResource := util.Extract(project.ContributeResource(resourceName, amount))

			err := discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						"# %s\n\nContributed %s of %s\n-# [%s] (%s of %s)",
						project.Name,
						util.ToUnit(amount),
						retResource.Name,
						util.MakeProgress(retResource.Amount, retResource.Goal, 15),
						util.ToUnit(retResource.Amount),
						util.ToUnit(retResource.Goal),
					),
				},
			})

			if err != nil {
				log.Fatalf("could not respond to interaction: %v", err)
			}
		}
	} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		options := []*discordgo.ApplicationCommandOptionChoice{}

		var selectedOption string
		for _, input := range i.ApplicationCommandData().Options[0].Options {
			if input.Focused {
				selectedOption = input.Name
				break
			}
		}

		switch selectedOption {
		case "project":
			for _, name := range db.GetProjectNames() {
				options = append(options, &discordgo.ApplicationCommandOptionChoice{
					Name:  name,
					Value: name,
				})
			}
		case "resource":
			optionMap := util.MakeOptionMap(i.ApplicationCommandData().Options[0].Options)
			projectName := util.Extract(util.GetString(optionMap, "project"))

			var project db.Project
			project.Get(projectName)

			for i, resource := range project.Resources {
				options = append(options, &discordgo.ApplicationCommandOptionChoice{
					Name:  resource.Name,
					Value: resource.Name,
				})
				if i == 24 {
					break
				}
			}
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
