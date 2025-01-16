package manga

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func errorMangaSetupNeededEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription("Please run the setup command first.").
		SetColor(0xff4f4f).
		Build()
	return embed
}

func selectServerOrUserComponents(command, subcommand, title string) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewSecondaryButton(
				"Server",
				"/"+command+"/"+subcommand+"/mode/server/"+title,
			),
			discord.NewSecondaryButton(
				"User (DMs)",
				"/"+command+"/"+subcommand+"/mode/user/"+title,
			),
		},
	}
}

func selectServerOrUserNestedComponents(
	command, subcommand, nested, title string,
) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewSecondaryButton(
				"Server",
				"/"+command+"/"+subcommand+"/"+nested+"/mode/server/"+title,
			),
			discord.NewSecondaryButton(
				"User (DMs)",
				"/"+command+"/"+subcommand+"/"+nested+"/mode/user/"+title,
			),
		},
	}
}

func selectServerOrUserEmbed(embedTitle, embedDescription string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(embedDescription).
		SetColor(0x3083e3).
		Build()

	return embed

}

func searchResultsEmbed(
	b *mubot.Bot,
	embedTitle, mangaTitle string,
) (discord.Embed, []searchResultsFormatted) {
	searchResults, err := utils.MuPostSearchSeries(b, mangaTitle)
	if err != nil {
		embed := discord.NewEmbedBuilder().
			SetTitle("Error").
			SetDescription("Failed to search for series. Try again later.").
			SetColor(0xff4f4f).
			Build()
		return embed, nil
	}

	description := "Select a manga from the search results:\n"
	if len(searchResults.Results) == 0 {
		description = "No results found for: `" + mangaTitle + "`. Try again or input a full https://mangaupdates.com link."
		return discord.NewEmbedBuilder().
			SetTitle(embedTitle).
			SetDescription(description).
			SetColor(0x3083e3).
			Build(), nil
	}

	allResults := []searchResultsFormatted{}
	for i, result := range searchResults.Results {
		if i >= 25 {
			break
		}
		description += fmt.Sprintf(
			"%d. %s (%s, Rating: %.2f)\n",
			i+1,
			utils.ParseHTMLEntities(result.Record.Title),
			result.Record.Year,
			result.Record.BayesianRating,
		)

		allResults = append(allResults, searchResultsFormatted{
			Title:  utils.ParseHTMLEntities(result.Record.Title),
			Year:   result.Record.Year,
			Rating: result.Record.BayesianRating,
			Id:     result.Record.SeriesID,
		})
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(description).
		SetColor(0x3083e3).
		Build()
	return embed, allResults
}

func dropdownSearchResultsComponents(
	command, subcommand, mode string,
	results []searchResultsFormatted,
) []discord.ContainerComponent {
	options := []discord.StringSelectMenuOption{}
	for i, result := range results {
		description := result.Year
		if result.Rating != 0 {
			if description != "" {
				description += ", "
			}
			description += fmt.Sprintf("Rating: %.2f", result.Rating)
		}
		options = append(options, discord.StringSelectMenuOption{
			Label:       fmt.Sprintf("%d. %s", i+1, utils.TruncateString(result.Title, 50)),
			Description: description,
			Value:       strconv.Itoa(result.Id),
		})
	}

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.StringSelectMenuComponent{
				CustomID:    "/" + command + "/" + subcommand + "/select/" + mode,
				Placeholder: "Select a Manga",
				Options:     options,
			},
		},
	}
}

func confirmMangaEmbed(b *mubot.Bot, embedTitle string, mangaId int64) discord.Embed {
	seriesInfo, err := utils.MuGetSeriesInfo(b, mangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to get series info (searchMangaAddHandler): %s", err.Error()),
		)
		return discord.NewEmbedBuilder().
			SetTitle("Error").
			SetDescription("Could not get series info, please try again later.").
			SetColor(0xff4f4f).
			Build()
	}

	description, err := utils.MuCleanupDescription(seriesInfo.Description)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to cleanup description (searchMangaAddHandler): %s", err.Error()),
		)
		return utils.DcErrorTechnicalErrorEmbed()
	}

	authorArray := []string{}
	for _, author := range seriesInfo.Authors {
		authorArray = append(authorArray, author.Name)
	}
	authorString := strings.Join(authorArray, ", ")

	embed := discord.NewEmbedBuilder().
		SetColor(0x3083e3).
		SetTitle(embedTitle).
		SetDescription(fmt.Sprintf("**Is `%s` the correct manga?**", seriesInfo.Title)).
		AddField("Description", description, false).
		AddField("Author(s)", authorString, true).
		AddField("Year", seriesInfo.Year, true).
		SetImage(seriesInfo.Image.URL.Original).
		Build()

	return embed
}

func selectConfirmMangaComponents(
	command, subcommand, mode, mangaId string,
) []discord.ContainerComponent {

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton(
				"Cancel",
				"/"+command+"/"+subcommand+"/confirm/select/"+mode+"/"+mangaId+"/cancel",
			),
			discord.NewSuccessButton(
				"Confirm",
				"/"+command+"/"+subcommand+"/confirm/select/"+mode+"/"+mangaId+"/confirm",
			),
		},
	}
}

func cancelEmbed(embedTitle string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription("Successfully cancelled.").
		SetColor(0xff4f4f).
		Build()
	return embed
}

func successMangaAddEmbed(embedTitle, mangaTitle string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(fmt.Sprintf("Successfully added `%s` to your manga list.", mangaTitle)).
		SetColor(0x3083e3).
		Build()
	return embed
}

func mangaExistsEmbed(embedTitle string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription("This manga is already in your list.").
		SetColor(0xff4f4f).
		Build()
	return embed
}

func dbMangaSearchResultsEmbed(
	embedTitle string,
	userManga []utils.MDbManga,
	page int,
) (discord.Embed, []dbMangaSearchResultsFormatted) {
	description := "Select a manga from your manga list:\n"
	if len(userManga) == 0 {
		description = "No manga found in your list."
		return discord.NewEmbedBuilder().
			SetTitle(embedTitle).
			SetDescription(description).
			SetColor(0x3083e3).
			Build(), nil
	}

	allResults := []dbMangaSearchResultsFormatted{}
	for i, result := range userManga {
		if i >= 25 {
			break
		}
		n := (page-1)*25 + i + 1
		str := fmt.Sprintf(
			"%d. %s\n",
			n,
			result.Title,
		)
		description += str

		allResults = append(allResults, dbMangaSearchResultsFormatted{
			Title: str,
			Id:    result.Id,
		})
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(description).
		SetColor(0x3083e3).
		Build()
	return embed, allResults
}

func dropdownDbMangaSearchResultsComponents(
	command, subcommand, mode string,
	results []dbMangaSearchResultsFormatted,
) []discord.ContainerComponent {
	options := []discord.StringSelectMenuOption{}
	for _, result := range results {
		options = append(options, discord.StringSelectMenuOption{
			Label: utils.TruncateString(result.Title, 50),
			Value: strconv.Itoa(int(result.Id)),
		})
	}

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.StringSelectMenuComponent{
				CustomID:    "/" + command + "/" + subcommand + "/select/" + mode,
				Placeholder: "Select a Manga",
				Options:     options,
			},
		},
	}
}

func dropdownDbMangaSearchResultsNestedComponents(
	command, subcommand, nested, mode string,
	results []dbMangaSearchResultsFormatted,
) []discord.ContainerComponent {
	options := []discord.StringSelectMenuOption{}
	for _, result := range results {
		options = append(options, discord.StringSelectMenuOption{
			Label: utils.TruncateString(result.Title, 50),
			Value: strconv.Itoa(int(result.Id)),
		})
	}

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.StringSelectMenuComponent{
				CustomID:    "/" + command + "/" + subcommand + "/" + nested + "/manga/select/" + mode,
				Placeholder: "Select a Manga",
				Options:     options,
			},
		},
	}
}

func paginationMangaSearchResultsComponents(
	command, subcommand, mode string,
	p parsedPaginationMangaList,
) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton(
				"",
				"/"+command+"/"+subcommand+"/search/mode/"+mode+"/"+strconv.Itoa(p.PrevPage),
			).
				WithEmoji(discord.ComponentEmoji{Name: "◀"}).
				WithDisabled(p.PrevPage == -1),
			discord.NewSecondaryButton(fmt.Sprintf("%d/%d", p.CurrentPage, p.MaxPage), "page-counter").
				WithDisabled(true),
			discord.NewSuccessButton(
				"",
				"/"+command+"/"+subcommand+"/search/mode/"+mode+"/"+strconv.Itoa(p.NextPage),
			).
				WithEmoji(discord.ComponentEmoji{Name: "▶"}).
				WithDisabled(p.NextPage == -1),
		},
	}
}

func paginationMangaSearchResultsNestedComponents(
	command, subcommand, nested, mode string,
	p parsedPaginationMangaList,
) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton(
				"",
				"/"+command+"/"+subcommand+"/"+nested+"/search/mode/"+mode+"/"+strconv.Itoa(
					p.PrevPage,
				),
			).
				WithEmoji(discord.ComponentEmoji{Name: "◀"}).
				WithDisabled(p.PrevPage == -1),
			discord.NewSecondaryButton(fmt.Sprintf("%d/%d", p.CurrentPage, p.MaxPage), "page-counter").
				WithDisabled(true),
			discord.NewSuccessButton(
				"",
				"/"+command+"/"+subcommand+"/"+nested+"/search/mode/"+mode+"/"+strconv.Itoa(
					p.NextPage,
				),
			).
				WithEmoji(discord.ComponentEmoji{Name: "▶"}).
				WithDisabled(p.NextPage == -1),
		},
	}
}

func dbMangaListEmbed(
	embedTitle string,
	userManga []utils.MDbManga,
) (*discord.EmbedBuilder, []dbMangaSearchResultsFormatted) {
	description := ""
	if len(userManga) == 0 {
		description = "No manga found in your list."
		return discord.NewEmbedBuilder().
			SetTitle(embedTitle).
			SetDescription(description).
			SetColor(0x3083e3), nil
	}

	allResults := []dbMangaSearchResultsFormatted{}
	for i, result := range userManga {
		if i >= 25 {
			break
		}
		str := fmt.Sprintf(
			"• %s\n",
			result.Title,
		)
		description += str

		allResults = append(allResults, dbMangaSearchResultsFormatted{
			Title: str,
			Id:    result.Id,
		})
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(description).
		SetColor(0x3083e3)
	return embed, allResults
}

func paginationMangaListComponents(
	command, subcommand, mode string,
	p parsedPaginationMangaList,
) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton(
				"",
				"/"+command+"/"+subcommand+"/p/mode/"+mode+"/"+strconv.Itoa(p.PrevPage),
			).
				WithEmoji(discord.ComponentEmoji{Name: "◀"}).
				WithDisabled(p.PrevPage == -1),
			discord.NewSecondaryButton(fmt.Sprintf("%d/%d", p.CurrentPage, p.MaxPage), "page-counter").
				WithDisabled(true),
			discord.NewSuccessButton(
				"",
				"/"+command+"/"+subcommand+"/p/mode/"+mode+"/"+strconv.Itoa(p.NextPage),
			).
				WithEmoji(discord.ComponentEmoji{Name: "▶"}).
				WithDisabled(p.NextPage == -1),
		},
	}
}

func parsePaginationMangaList(
	mangaList []utils.MDbManga,
	page int,
) parsedPaginationMangaList {
	const pageSize = 25
	totalMangas := len(mangaList)
	totalPages := (totalMangas + pageSize - 1) / pageSize

	if totalMangas <= pageSize {
		return parsedPaginationMangaList{
			Pagination:  false,
			PrevPage:    -1,
			CurrentPage: 1,
			NextPage:    -1,
			MaxPage:     1,
			MangaList:   mangaList,
		}
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > totalMangas {
		endIndex = totalMangas
	}
	var prevPage, nextPage int
	if page > 1 {
		prevPage = page - 1
	} else {
		prevPage = -1
	}
	if page < totalPages {
		nextPage = page + 1
	} else {
		nextPage = -1
	}

	return parsedPaginationMangaList{
		Pagination:  true,
		PrevPage:    prevPage,
		CurrentPage: page,
		NextPage:    nextPage,
		MaxPage:     totalPages,
		MangaList:   mangaList[startIndex:endIndex],
	}
}

func successMangaRemoveEmbed(embedTitle string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription("Manga successfully removed from your manga list.").
		SetColor(0x3083e3).
		Build()
	return embed
}

func parsePaginationSeriesGroups(
	seriesGroups []utils.MuGroup,
	page int,
) parsedPaginationGroupList {
	const pageSize = 25
	totalMangas := len(seriesGroups)
	totalPages := (totalMangas + pageSize - 1) / pageSize

	if totalMangas <= pageSize {
		return parsedPaginationGroupList{
			Pagination:  false,
			PrevPage:    -1,
			CurrentPage: 1,
			NextPage:    -1,
			MaxPage:     1,
			GroupList:   seriesGroups,
		}
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > totalMangas {
		endIndex = totalMangas
	}
	var prevPage, nextPage int
	if page > 1 {
		prevPage = page - 1
	} else {
		prevPage = -1
	}
	if page < totalPages {
		nextPage = page + 1
	} else {
		nextPage = -1
	}

	return parsedPaginationGroupList{
		Pagination:  true,
		PrevPage:    prevPage,
		CurrentPage: page,
		NextPage:    nextPage,
		MaxPage:     totalPages,
		GroupList:   seriesGroups[startIndex:endIndex],
	}
}

func paginationGroupListNestedComponents(
	command, subcommand, nested, mode, mangaId string,
	p parsedPaginationGroupList,
) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton(
				"",
				"/"+command+"/"+subcommand+"/"+nested+"/groups/p/"+mangaId+"/mode/"+mode+"/"+strconv.Itoa(
					p.PrevPage,
				),
			).
				WithEmoji(discord.ComponentEmoji{Name: "◀"}).
				WithDisabled(p.PrevPage == -1),
			discord.NewSecondaryButton(fmt.Sprintf("%d/%d", p.CurrentPage, p.MaxPage), "page-counter").
				WithDisabled(true),
			discord.NewSuccessButton(
				"",
				"/"+command+"/"+subcommand+"/"+nested+"/groups/p/"+mangaId+"/mode/"+mode+"/"+strconv.Itoa(
					p.NextPage,
				),
			).
				WithEmoji(discord.ComponentEmoji{Name: "▶"}).
				WithDisabled(p.NextPage == -1),
		},
	}
}

func selectSeriesGroupEmbed(
	embedTitle string,
	seriesGroup []utils.MuGroup,
) (*discord.EmbedBuilder, []dbMangaSearchResultsFormatted) {
	description := ""
	if len(seriesGroup) == 0 {
		description = "No groups found scanlating this manga."
		return discord.NewEmbedBuilder().
			SetTitle(embedTitle).
			SetDescription(description).
			SetColor(0x3083e3), nil
	}

	allResults := []dbMangaSearchResultsFormatted{}
	for i, result := range seriesGroup {
		if i >= 25 {
			break
		}
		str := fmt.Sprintf(
			"• %s\n",
			result.Name,
		)
		description += str

		allResults = append(allResults, dbMangaSearchResultsFormatted{
			Title: result.Name,
			Id:    int64(result.GroupID),
		})
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(description).
		SetColor(0x3083e3)
	return embed, allResults
}

func dropdownSeriesGroupsNestedComponents(
	command, subcommand, nested, mode, mangaId string,
	results []dbMangaSearchResultsFormatted,
) []discord.ContainerComponent {
	options := []discord.StringSelectMenuOption{}
	for _, result := range results {
		options = append(options, discord.StringSelectMenuOption{
			Label: utils.TruncateString(result.Title, 50),
			Value: strconv.Itoa(int(result.Id)),
		})
	}

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.StringSelectMenuComponent{
				CustomID:    "/" + command + "/" + subcommand + "/" + nested + "/groups/select/" + mode + "/" + mangaId,
				Placeholder: "Select a Scanlator Group",
				Options:     options,
			},
		},
	}
}

func confirmGroupEmbed(b *mubot.Bot, embedTitle string, groupId int64) discord.Embed {
	groupInfo, err := utils.MuGetGroupInfo(b, groupId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to get series info (searchMangaAddHandler): %s", err.Error()),
		)
		return discord.NewEmbedBuilder().
			SetTitle("Error").
			SetDescription("Could not get series info, please try again later.").
			SetColor(0xff4f4f).
			Build()
	}

	activeStatement := "*A group is designated as inactive if it hasn't released in the past 6 months*"

	var active string
	if groupInfo.Active {
		active = "Yes"
	} else {
		active = "No"
	}

	embed := discord.NewEmbedBuilder().
		SetColor(0x3083e3).
		SetTitle(embedTitle).
		SetDescription(fmt.Sprintf("**Is `%s` the correct group?**", groupInfo.Name)).
		AddField("Active", fmt.Sprintf("%s (%s)", active, activeStatement), false)

	if groupInfo.Social.Site != "" {
		embed.AddField("Website", groupInfo.Social.Site, true)
	}

	if groupInfo.Social.Discord != "" {
		embed.AddField("Discord", groupInfo.Social.Discord, true)
	}

	if groupInfo.Social.Twitter != "" {
		embed.AddField("Twitter", groupInfo.Social.Twitter, true)
	}

	return embed.Build()
}

func selectConfirmGroupComponents(
	command, subcommand, nested, mode, mangaId, groupId string,
) []discord.ContainerComponent {

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton(
				"Cancel",
				"/"+command+"/"+subcommand+"/"+nested+"/confirm/select/"+mode+"/"+mangaId+"/"+groupId+"/cancel",
			),
			discord.NewSuccessButton(
				"Confirm",
				"/"+command+"/"+subcommand+"/"+nested+"/confirm/select/"+mode+"/"+mangaId+"/"+groupId+"/confirm",
			),
		},
	}
}

func groupExistsEmbed(embedTitle string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription("This scanlator group is already added for this manga.").
		SetColor(0xff4f4f).
		Build()
	return embed
}

func successMangaSetScanlatorEmbed(embedTitle, mangaName, groupName string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(fmt.Sprintf("Scanlator group, `%s`, successfully added for the manga `%s`.", groupName, mangaName)).
		SetColor(0x3083e3).
		Build()
	return embed
}
