package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func MangaRemoveHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	_, inGuild := e.Guild()
	if !inGuild {
		adapter := &CommandEventAdapter{Event: e}
		return MangaRemoveUserHandler(adapter, b)
	} else {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds: []discord.Embed{
					selectServerOrUserEmbed(
						"Remove Manga",
						"Would you like to remove a manga from this server or your DMs?",
					),
				},
				Components: selectServerOrUserComponents("manga", "remove", ""),
			},
		)
	}
}

func MangaRemoveUserHandler(
	e EventHandler,
	b *mubot.Bot,
) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckExists(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to check if user exists (MangaRemoveUserHandler): %s", err.Error()),
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

	return MangaRemoveUserSearchHandler(e, b, 1)
}

func MangaRemoveServerHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
) error {
	serverId := int64(*e.GuildID())

	exists, err := utils.DbServerCheckExists(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if server exists (MangaRemoveServerHandler): %s",
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

	return MangaRemoveServerSearchHandler(e, b, 1)
}

func MangaRemoveUserSearchHandler(
	e EventHandler,
	b *mubot.Bot,
	page int,
) error {
	userId := int64(e.User().ID)

	user, err := utils.DbGetUser(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get user (MangaRemoveUserHandler): %s",
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
		"Remove Manga",
		parsed.MangaList,
		page,
	)
	if searchResultsFormatted == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(searchResults).Build(),
		)
	}
	dropdownSearchResults := dropdownDbMangaSearchResultsComponents(
		"manga",
		"remove",
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

	pagination := paginationMangaSearchResultsComponents(
		"manga",
		"remove",
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

func MangaRemoveServerSearchHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
) error {
	serverId := int64(*e.GuildID())

	server, err := utils.DbGetServer(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get server (MangaRemoveServerHandler): %s",
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
		"Remove Manga",
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
	dropdownSearchResults := dropdownDbMangaSearchResultsComponents(
		"manga",
		"remove",
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

	pagination := paginationMangaSearchResultsComponents(
		"manga",
		"remove",
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

func MangaRemoveSelectHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mode string,
) error {
	mangaId := e.StringSelectMenuInteractionData().Values[0]

	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to parse manga ID (MangaRemoveSelectHandler): %s", mangaId),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	components := selectConfirmMangaComponents("manga", "remove", mode, mangaId)
	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{confirmMangaEmbed(b, "Remove Manga", intMangaId)},
			Components: &components,
		},
	)

	return err
}

func MangaRemoveCancelHandler(e *handler.ComponentEvent) error {
	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{cancelEmbed("Remove Manga")},
			Components: &[]discord.ContainerComponent{},
		},
	)
}

func MangaRemoveConfirmHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mode,
	mangaId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to parse manga ID (MangaRemoveConfirmHandler): %s", mangaId),
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
		err = utils.DbUserRemoveManga(b, userId, intMangaId)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to remove manga from user (MangaRemoveConfirmHandler): %s",
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
		err = utils.DbServerRemoveManga(b, serverId, intMangaId)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf("Failed to remove manga from server (MangaRemoveConfirmHandler): %s", err.Error()),
			)
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}
	}

	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{successMangaRemoveEmbed("Remove Manga")},
			Components: &[]discord.ContainerComponent{},
		},
	)
}
