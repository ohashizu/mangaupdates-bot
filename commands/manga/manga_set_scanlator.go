package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func MangaSetScanlatorHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	_, inGuild := e.Guild()
	if !inGuild {
		adapter := &CommandEventAdapter{Event: e}
		return MangaSetScanlatorUserHandler(adapter, b)
	} else {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds: []discord.Embed{
					selectServerOrUserEmbed(
						"Set Scanlator",
						"Would you like to set a manga's scanlator from this server or your DMs?",
					),
				},
				Components: selectServerOrUserNestedComponents("manga", "set", "scanlator", ""),
			},
		)
	}
}

func MangaSetScanlatorUserHandler(
	e EventHandler,
	b *mubot.Bot,
) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckExists(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if user exists (MangaSetScanlatorUserHandler): %s",
				err.Error(),
			),
		)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(utils.DcErrorTechnicalErrorEmbed()).Build(),
		)
	}
	if !exists {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorMangaSetupNeededEmbed()).Build(),
		)
	}

	return MangaSetScanlatorUserSearchHandler(e, b, 1)
}

func MangaSetScanlatorServerHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
) error {
	serverId := int64(*e.GuildID())

	exists, err := utils.DbServerCheckExists(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if server exists (MangaSetScanlatorServerHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	if !exists {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorMangaSetupNeededEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	return MangaSetScanlatorServerSearchHandler(e, b, 1)
}

func MangaSetScanlatorUserSearchHandler(
	e EventHandler,
	b *mubot.Bot,
	page int,
) error {
	userId := int64(e.User().ID)

	user, err := utils.DbGetUser(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get user (MangaSetScanlatorUserHandler): %s",
				err.Error(),
			),
		)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(utils.DcErrorTechnicalErrorEmbed()).Build(),
		)
	}

	parsed := parsePaginationMangaList(user.Manga, page)

	searchResults, searchResultsFormatted := dbMangaSearchResultsEmbed(
		"Set Scanlator",
		parsed.MangaList,
		page,
	)
	if searchResultsFormatted == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(searchResults).Build(),
		)
	}
	dropdownSearchResults := dropdownDbMangaSearchResultsNestedComponents(
		"manga",
		"set",
		"scanlator",
		"user",
		searchResultsFormatted,
	)

	if !parsed.Pagination {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds:     []discord.Embed{searchResults},
				Components: dropdownSearchResults,
			},
		)
	}

	pagination := paginationMangaSearchResultsNestedComponents(
		"manga",
		"set",
		"scanlator",
		"user",
		parsed,
	)
	components := append(dropdownSearchResults, pagination...)
	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.MessageCreate{
			Embeds:     []discord.Embed{searchResults},
			Components: components,
		},
	)
}

func MangaSetScanlatorServerSearchHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
) error {
	serverId := int64(*e.GuildID())

	server, err := utils.DbGetServer(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get server (MangaSetScanlatorServerHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	parsed := parsePaginationMangaList(server.Manga, page)

	searchResults, searchResultsFormatted := dbMangaSearchResultsEmbed(
		"Set Scanlator",
		parsed.MangaList,
		page,
	)
	if searchResultsFormatted == nil {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	dropdownSearchResults := dropdownDbMangaSearchResultsNestedComponents(
		"manga",
		"set",
		"scanlator",
		"server",
		searchResultsFormatted,
	)

	if !parsed.Pagination {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults},
				Components: &dropdownSearchResults,
			},
		)
	}

	pagination := paginationMangaSearchResultsNestedComponents(
		"manga",
		"set",
		"scanlator",
		"server",
		parsed,
	)
	components := append(dropdownSearchResults, pagination...)
	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{searchResults},
			Components: &components,
		},
	)
}

func MangaSetScanlatorMangaSelectHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mode string,
) error {
	mangaId := e.StringSelectMenuInteractionData().Values[0]

	return MangaSetScanlatorGroupHandler(e, b, 1, mode, mangaId)
}

func MangaSetScanlatorGroupHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
	mode,
	mangaId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse manga ID (MangaSetScanlatorSelectGroupHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	groups, err := utils.MuGetSeriesGroups(b, intMangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get manga (MangaSetScanlatorUserSelectGroupHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	parsed := parsePaginationSeriesGroups(groups.GroupList, page)

	searchResults, searchResultsFormatted := selectSeriesGroupEmbed(
		"Set Scanlator",
		parsed.GroupList,
	)
	if searchResultsFormatted == nil {
		fmt.Println("no search results formatted")
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults.Build()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	dropdownSearchResults := dropdownSeriesGroupsNestedComponents(
		"manga",
		"set",
		"scanlator",
		mode,
		mangaId,
		searchResultsFormatted,
	)

	if !parsed.Pagination {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults.Build()},
				Components: &dropdownSearchResults,
			},
		)
	}

	pagination := paginationGroupListNestedComponents(
		"manga",
		"set",
		"scanlator",
		mode,
		mangaId,
		parsed,
	)

	components := append(dropdownSearchResults, pagination...)
	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{searchResults.Build()},
			Components: &components,
		},
	)

	return err
}

func MangaSetScanlatorGroupSelectHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mode,
	mangaId string,
) error {
	groupId := e.StringSelectMenuInteractionData().Values[0]

	intGroupId, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse group ID (MangaSetScanlatorGroupSelectHandler): %s",
				groupId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	components := selectConfirmGroupComponents("manga", "set", "scanlator", mode, mangaId, groupId)
	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{confirmGroupEmbed(b, "Set Scanlator", intGroupId)},
			Components: &components,
		},
	)

	return err
}

func MangaSetScanlatorGroupCancelHandler(
	e *handler.ComponentEvent,
) error {
	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{cancelEmbed("Set Scanlator")},
			Components: &[]discord.ContainerComponent{},
		},
	)
}

func MangaSetScanlatorGroupConfirmHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mode,
	mangaId,
	groupId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse manga id (MangaSetScanlatorGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	intGroupId, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse group id (MangaSetScanlatorGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	groupInfo, err := utils.MuGetGroupInfo(b, intGroupId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get group info (MangaSetScanlatorGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	if mode == "user" {
		userId := int64(e.User().ID)
		exists, err := utils.DbUserCheckGroupExists(b, userId, intMangaId, intGroupId)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to check group exists from user (MangaSetScanlatorGroupConfirmHandler): %s",
					err.Error(),
				),
			)
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}
		if exists {
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{groupExistsEmbed("Set Scanlator")},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}
		err = utils.DbUserAddGroup(b, userId, intMangaId, intGroupId, groupInfo.Name)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to set scanlator from user (MangaSetScanlatorGroupConfirmHandler): %s",
					err.Error(),
				),
			)
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}
	} else {
		serverId := int64(*e.GuildID())
		exists, err := utils.DbServerCheckGroupExists(b, serverId, intMangaId, intGroupId)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to check group exists from server (MangaSetScanlatorGroupConfirmHandler): %s",
					err.Error(),
				),
			)
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}
		if exists {
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{groupExistsEmbed("Set Scanlator")},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}

		err = utils.DbServerAddGroup(b, serverId, intMangaId, intGroupId, groupInfo.Name)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to set scanlator from server (MangaSetScanlatorGroupConfirmHandler): %s",
					err.Error(),
				),
			)
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}
	}

	seriesInfo, err := utils.MuGetSeriesInfo(b, intMangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get series info (MangaSetScanlatorGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				successMangaSetScanlatorEmbed("Set Scanlator", seriesInfo.Title, groupInfo.Name),
			},
			Components: &[]discord.ContainerComponent{},
		},
	)
}
