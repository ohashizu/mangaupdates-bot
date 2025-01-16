package manga

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func MangaListHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	_, inGuild := e.Guild()
	if !inGuild {
		adapter := &CommandEventAdapter{Event: e}
		return MangaListUserHandler(adapter, b)
	} else {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds: []discord.Embed{
					selectServerOrUserEmbed(
						"Manga List",
						"Would you like to view the manga list from this server or your DMs?",
					),
				},
				Components: selectServerOrUserComponents("manga", "list", ""),
			},
		)
	}
}

func MangaListUserHandler(
	e EventHandler,
	b *mubot.Bot,
) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckExists(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to check if user exists (MangaListUserHandler): %s", err.Error()),
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

	return MangaListUserListHandler(e, b, 1)
}

func MangaListServerHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
) error {
	serverId := int64(*e.GuildID())

	exists, err := utils.DbServerCheckExists(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if server exists (MangaListServerHandler): %s",
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

	return MangaListServerListHandler(e, b, 1)
}

func MangaListUserListHandler(
	e EventHandler,
	b *mubot.Bot,
	page int,
) error {
	userId := int64(e.User().ID)

	user, err := utils.DbGetUser(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get user (MangaListUserListHandler): %s",
				err.Error(),
			),
		)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(utils.DcErrorTechnicalErrorEmbed()).Build(),
		)
	}

	parsed := parsePaginationMangaList(user.Manga, page)

	usr := e.User()
	mangaList, mangaListFormatted := dbMangaListEmbed(
		fmt.Sprintf("%s's Manga List", usr.EffectiveName()),
		parsed.MangaList,
	)
	avatar := usr.AvatarURL()
	if *avatar == "" {
		*avatar = e.User().DefaultAvatarURL()
	}
	botUser, _ := e.SelfUser()
	builtML := mangaList.SetThumbnail(*avatar).
		SetAuthor("MangaUpdates", "", *botUser.AvatarURL()).
		Build()

	if mangaListFormatted == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(builtML).Build(),
		)
	}

	if !parsed.Pagination {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds: []discord.Embed{builtML},
			},
		)
	}

	pagination := paginationMangaListComponents(
		"manga",
		"list",
		"user",
		parsed,
	)
	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.MessageCreate{
			Embeds:     []discord.Embed{builtML},
			Components: pagination,
		},
	)
}

func MangaListServerListHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
) error {
	serverId := int64(*e.GuildID())

	server, err := utils.DbGetServer(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get server (MangaListServerListHandler): %s",
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

	guild, _ := e.Guild()
	mangaList, mangaListFormatted := dbMangaListEmbed(
		fmt.Sprintf("%s's Manga List", guild.Name),
		parsed.MangaList,
	)
	avatar := guild.IconURL()
	if *avatar == "" {
		*avatar = "https://cdn.discordapp.com/embed/avatars/0.png"
	}
	botUser, _ := e.Client().Caches().SelfUser()
	builtML := mangaList.SetThumbnail(*avatar).
		SetAuthor("MangaUpdates", "", *botUser.AvatarURL()).
		Build()

	if mangaListFormatted == nil {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{builtML},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	if !parsed.Pagination {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{builtML},
			},
		)
	}

	pagination := paginationMangaListComponents(
		"manga",
		"list",
		"server",
		parsed,
	)
	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{builtML},
			Components: &pagination,
		},
	)
}
