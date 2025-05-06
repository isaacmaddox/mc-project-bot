package project

import (
	"fmt"
	"log"
	"slices"

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
		{
			Name:         "sortby",
			Description:  "Sort list of resources",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     false,
			Autocomplete: true,
		},
	},
}

var OverviewHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		optionMap := util.MakeOptionMap(i.ApplicationCommandData().Options[0].Options)

		name := util.Extract(util.GetString(optionMap, "name"))
		page := max(util.GetIntOptional(optionMap, "page", 1), 1)
		sortBy := util.GetStringOptional(optionMap, "sortby", "")

		OverviewPaginatedHandler(discord, i, name, page, sortBy)
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
		case "name":
			for i, name := range db.GetProjectNames() {
				options = append(options, &discordgo.ApplicationCommandOptionChoice{
					Name:  name,
					Value: name,
				})
				if i == 24 {
					break
				}
			}
		case "sortby":
			options = []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Greatest amount",
					Value: "greatest amount",
				},
				{
					Name:  "Least amount",
					Value: "least amount",
				},
				{
					Name:  "Most needed",
					Value: "most needed",
				},
				{
					Name:  "Least needed",
					Value: "least needed",
				},
				{
					Name:  "Greatest goal",
					Value: "greatest goal",
				},
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

var OverviewPaginatedHandler = func(discord *discordgo.Session, i *discordgo.InteractionCreate, name string, page int64, sortBy string) {
	var project db.Project
	project.Get(name)

	resourceList := []string{""}
	var projectAmount int
	var projectGoal int

	var areMorePages bool = false

	resources := project.Resources

	switch sortBy {
	case "greatest amount":
		slices.SortFunc(resources,
			func(a *db.Resource, b *db.Resource) int {
				ret := max(min(b.Amount-a.Amount, 1), -1)
				if ret == 0 {
					ret = max(min(b.Goal-a.Goal, 1), -1)
				}

				return ret
			},
		)
	case "least amount":
		slices.SortFunc(resources,
			func(a *db.Resource, b *db.Resource) int {
				ret := max(min(a.Amount-b.Amount, 1), -1)
				if ret == 0 {
					ret = max(min(a.Goal-b.Goal, 1), -1)
				}

				return ret
			},
		)
	case "most needed":
		slices.SortFunc(resources,
			func(a *db.Resource, b *db.Resource) int {
				aNeeded := a.Goal - a.Amount
				bNeeded := b.Goal - b.Amount

				ret := max(min(bNeeded-aNeeded, 1), -1)
				if ret == 0 {
					ret = max(min(b.Goal-a.Goal, 1), -1)
				}

				return ret
			},
		)
	case "least needed":
		slices.SortFunc(resources,
			func(a *db.Resource, b *db.Resource) int {
				aNeeded := a.Goal - a.Amount
				bNeeded := b.Goal - b.Amount

				ret := max(min(aNeeded-bNeeded, 1), -1)
				if ret == 0 {
					ret = max(min(a.Goal-b.Goal, 1), -1)
				}

				return ret
			},
		)
	case "greatest goal", "":
		slices.SortFunc(resources,
			func(a *db.Resource, b *db.Resource) int {
				ret := max(min(b.Goal-a.Goal, 1), -1)
				if ret == 0 {
					ret = max(min(a.Amount-b.Amount, 1), -1)
				}

				return ret
			},
		)
	}

	for _, resource := range resources {
		projectAmount += min(resource.Amount, resource.Goal)
		projectGoal += resource.Goal

		if !areMorePages {
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
				}
				resourceList = append(resourceList, resourceString)
			}
		}
	}

	var header string
	if projectGoal != 0 {
		progressBar := util.MakeProgress(projectAmount, projectGoal, 25)
		progressPercent := (float32(projectAmount) / float32(projectGoal)) * 100

		header = fmt.Sprintf("# %s\n-# **%s**\n[%s] (%.0f%%)\n\n",
			project.Name,
			project.Description,
			progressBar,
			progressPercent,
		)
	} else {
		header = fmt.Sprintf("# %s\n-# **%s**\n add resources with /project resource new %v\n",
			project.Name,
			project.Description,
			project.Name,
		)
	}

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
			Content: header + resourceList[page-1] + "-# .",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Previous",
							Style:    discordgo.SecondaryButton,
							CustomID: fmt.Sprintf("project_overview_%d_%s_%s", page-1, sortBy, name),
							Disabled: page == 1,
						},
						discordgo.Button{
							Label:    "Next",
							Style:    discordgo.SecondaryButton,
							CustomID: fmt.Sprintf("project_overview_%d_%s_%s", page+1, sortBy, name),
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
