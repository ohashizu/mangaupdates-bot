package manga

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/jckli/mangaupdates-bot/utils"
)

type EventHandler interface {
	User() discord.User
	Respond(
		responseType discord.InteractionResponseType,
		data discord.InteractionResponseData,
		opts ...rest.RequestOpt,
	) error
	UpdateInteractionResponse(
		messageUpdate discord.MessageUpdate,
		opts ...rest.RequestOpt,
	) (*discord.Message, error)
	SelfUser() (discord.OAuth2User, bool)
}

type CommandEventAdapter struct {
	Event *handler.CommandEvent
}

func (c *CommandEventAdapter) User() discord.User {
	return c.Event.User()
}

func (c *CommandEventAdapter) Respond(
	responseType discord.InteractionResponseType,
	data discord.InteractionResponseData,
	opts ...rest.RequestOpt,
) error {
	return c.Event.Respond(responseType, data, opts...)
}

func (c *CommandEventAdapter) UpdateInteractionResponse(
	messageUpdate discord.MessageUpdate,
	opts ...rest.RequestOpt,
) (*discord.Message, error) {
	return c.Event.UpdateInteractionResponse(messageUpdate, opts...)
}

func (c *CommandEventAdapter) SelfUser() (discord.OAuth2User, bool) {
	return c.Event.Client().Caches().SelfUser()
}

type ComponentEventAdapter struct {
	Event *handler.ComponentEvent
}

func (c *ComponentEventAdapter) User() discord.User {
	return c.Event.User()
}

func (c *ComponentEventAdapter) Respond(
	responseType discord.InteractionResponseType,
	data discord.InteractionResponseData,
	opts ...rest.RequestOpt,
) error {
	return c.Event.Respond(responseType, data, opts...)
}

func (c *ComponentEventAdapter) UpdateInteractionResponse(
	messageUpdate discord.MessageUpdate,
	opts ...rest.RequestOpt,
) (*discord.Message, error) {
	return c.Event.UpdateInteractionResponse(messageUpdate, opts...)
}

func (c *ComponentEventAdapter) SelfUser() (discord.OAuth2User, bool) {
	return c.Event.Client().Caches().SelfUser()
}

type searchResultsFormatted struct {
	Title  string
	Year   string
	Rating float64
	Id     int
}

type dbMangaSearchResultsFormatted struct {
	Title string
	Id    int64
}

type parsedPaginationMangaList struct {
	Pagination  bool
	PrevPage    int
	CurrentPage int
	NextPage    int
	MaxPage     int
	MangaList   []utils.MDbManga
}

type parsedPaginationGroupList struct {
	Pagination  bool
	PrevPage    int
	CurrentPage int
	NextPage    int
	MaxPage     int
	GroupList   []utils.MuGroup
}
