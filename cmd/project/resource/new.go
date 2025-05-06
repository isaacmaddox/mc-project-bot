package resource

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/isaacmaddox/mc-project-bot/db"
	"github.com/isaacmaddox/mc-project-bot/util"
)

var NewResourceCommand = &discordgo.ApplicationCommandOption{
	Name:        "new",
	Description: "Add new resource requirement",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "project",
			Description:  "Which project to add resource to",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Name:        "resource",
			Description: "Whatchu want fool",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
		{
			Name:        "goal",
			Description: "How much do ye need?",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
		{
			Name:        "amount",
			Description: "How much do ye have?",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    false,
		},
	},
}

var NewResourceHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		NewResource(discord, i)
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

func NewResource(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := util.MakeOptionMap(i.ApplicationCommandData().Options[0].Options[0].Options)

	projectName := util.Extract(util.GetString(optionMap, "project"))
	resourceName := util.Extract(util.GetString(optionMap, "resource"))
	goal := util.FromUnit(util.Extract(util.GetString(optionMap, "goal")))
	amount := util.FromUnit(util.GetStringOptional(optionMap, "amount", "0"))

	var project db.Project
	found := project.Get(projectName)

	if !found {
		discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No project found with the name %s", projectName),
			},
		})
		return
	}

	resource := project.AddResource(resourceName, amount, goal)
	returnAmount := util.ToUnit(resource.Amount)

	if returnAmount == "" {
		returnAmount = "0"
	}

	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"# %s\n\n*Added resource*: **%s**\n-# [%s] (%s of %s)",
				project.Name,
				resource.Name,
				util.MakeProgress(resource.Amount, resource.Goal, 15),
				returnAmount,
				util.ToUnit(resource.Goal),
			),
		},
	})
}
