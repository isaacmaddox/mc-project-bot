package util

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type OptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func MakeOptionMap(options []*discordgo.ApplicationCommandInteractionDataOption) *OptionMap {
	optionMap := make(OptionMap, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return &optionMap
}

func GetString(optionMap *OptionMap, key string) (string, error) {
	if opt, ok := (*optionMap)[key]; ok {
		return opt.StringValue(), nil
	}
	return "", errors.New("ruh roh")
}

func GetInt(optionMap *OptionMap, key string) (int64, error) {
	if opt, ok := (*optionMap)[key]; ok {
		return opt.IntValue(), nil
	}
	return -1, errors.New("ruh roh")
}

func GetStringOptional(optionMap *OptionMap, key string, backup string) string {
	if opt, ok := (*optionMap)[key]; ok {
		return opt.StringValue()
	}

	return backup
}

func GetIntOptional(optionMap *OptionMap, key string, backup int64) int64 {
	if opt, ok := (*optionMap)[key]; ok {
		return opt.IntValue()
	}

	return backup
}
