import os
import discord
from discord.ext import commands 
from discord.commands import Option, SlashCommandGroup
from core.mongodb import Mongo
from core.mangaupdates import MangaUpdates
from core.utils import Util
import validators

mangaupdates = MangaUpdates()
mongo = Mongo()
util = Util()

timeoutError = discord.Embed(title="Error", description="You didn't respond in time! Please rerun the command.", color=0xff4f4f)
ghuser = os.environ.get("GITHUB_USER")

class Mode(discord.ui.View):
    def __init__(self):
        super().__init__(timeout=15.0)
        self.value = None
        self.interaction = None

    @discord.ui.button(label=f'User (DMs)', style=discord.ButtonStyle.grey)
    async def confirm(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = "user"
        self.interaction = interaction
        self.stop()

    @discord.ui.button(label=f'Server', style=discord.ButtonStyle.grey)
    async def cancel(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = "server"
        self.interaction = interaction
        self.stop()

class Confirm(discord.ui.View):
    def __init__(self):
        super().__init__(timeout=15.0)
        self.value = None
        self.interaction = None

    @discord.ui.button(label="Confirm", style=discord.ButtonStyle.green)
    async def confirm(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = True
        self.interaction = interaction
        self.stop()

    @discord.ui.button(label="Cancel", style=discord.ButtonStyle.red)
    async def cancel(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = False
        self.interaction = interaction
        self.stop()

# manga add views
class SelectMangaView(discord.ui.View):
    def __init__(self, manga_list, mode):
        super().__init__(timeout=15.0)
        self.select_manga = SelectMangaWConfirm(manga_list=manga_list, mode=mode)
        self.add_item(self.select_manga)

class SelectMangaWConfirm(discord.ui.Select):
    def __init__(self, manga_list, mode):
        manga_desc = []
        for manga in manga_list:
            manga_desc.append(discord.SelectOption(label=manga["name"], description=manga["description"]))
        super().__init__(
            placeholder="Choose a manga series...",
            min_values=1,
            max_values=1,
            options=manga_desc
        )
        self.manga_list = manga_list
        self.finish = None
        self.modeval = mode

    async def callback(self, interaction: discord.Interaction):
        value = int(self.values[0][:2].replace(".", "")) - 1
        title = self.manga_list[value]["info"]["title"]
        fulldescription = self.manga_list[value]["info"]["description"]
        image = self.manga_list[value]["info"]["image"]["url"]["original"]
        result = discord.Embed(title=f"Did you mean to add `{title}`?", color=0x3083e3, description=util.format_mu_description(fulldescription))
        result.set_image(url=image)
        confirm = Confirm()
        await interaction.response.edit_message(embed=result, view=confirm)
        mangaid = self.manga_list[value]["id"]
        manganame = title
        self.finish = True
        await confirm.wait()
        if confirm.value is None:
            await interaction.message.edit(embed=timeoutError, view=None)
            return
        elif confirm.value is False:
            cancelEmbed = discord.Embed(title=f"Canceled", color=0x3083e3, description="Successfully canceled.")
            await confirm.interaction.response.edit_message(embed=cancelEmbed, view=None)
            return
        else:
            if self.modeval == "user":
                mangaindb = await mongo.check_manga_exist_user(confirm.interaction.user.id, mangaid)
            elif self.modeval == "server":
                mangaindb = await mongo.check_manga_exist_server(confirm.interaction.guild.id, mangaid)
            if mangaindb is True:
                mangaExist = discord.Embed(title="Add Manga", color=0x3083e3, description="This manga is already added to your list.")
                await confirm.interaction.response.edit_message(embed=mangaExist, view=None)
            elif mangaindb is False:
                if self.modeval == "user":
                    await mongo.add_manga_user(confirm.interaction.user.id, mangaid, manganame)
                elif self.modeval == "server":
                    await mongo.add_manga_server(confirm.interaction.guild.id, mangaid, manganame)
                mangaAdded = discord.Embed(title="Add Manga", color=0x3083e3, description="Manga succesfully added.")
                await confirm.interaction.response.edit_message(embed=mangaAdded, view=None)

# manga remove views
class SelectMangaRemoveView(discord.ui.View):
    def __init__(self, manga_list, mode):
        super().__init__(timeout=15.0)
        self.select_manga = SelectMangaRemove(manga_list=manga_list, mode=mode)
        self.add_item(self.select_manga)

class SelectMangaRemove(discord.ui.Select):
    def __init__(self, manga_list, mode):
        manga_desc = []
        for manga in manga_list:
            manga_desc.append(discord.SelectOption(label=manga["dropdownTitle"]))
        super().__init__(
            placeholder="Choose a manga series...",
            min_values=1,
            max_values=1,
            options=manga_desc
        )
        self.manga_list = manga_list
        self.finish = None
        self.modeval = mode

    async def callback(self, interaction: discord.Interaction):
        value = int(self.values[0][:2].replace(".", "")) - 1
        title = self.manga_list[value]["title"]
        search_data = await mangaupdates.series_info(self.manga_list[value]["id"])
        fulldescription = search_data["description"]
        image = search_data["image"]["url"]["original"]
        result = discord.Embed(title=f"Are you sure you want to remove `{title}`?", color=0x3083e3, description=util.format_mu_description(fulldescription))
        result.set_image(url=image)
        confirm = Confirm()
        await interaction.response.edit_message(embed=result, view=confirm)
        mangaid = self.manga_list[value]["id"]
        self.finish = True
        await confirm.wait()
        if confirm.value is None:
            await interaction.message.edit(embed=timeoutError, view=None)
            return
        elif confirm.value is False:
            cancelEmbed = discord.Embed(title=f"Canceled", color=0x3083e3, description="Successfully canceled.")
            await confirm.interaction.response.edit_message(embed=cancelEmbed, view=None)
            return
        else:
            if self.modeval == "user":
                await mongo.remove_manga_user(confirm.interaction.user.id, mangaid)
            elif self.modeval == "server":
                await mongo.remove_manga_server(confirm.interaction.guild.id, mangaid)
            mangaRemoved = discord.Embed(title="Remove Manga", color=0x3083e3, description="Manga succesfully removed.")
            await confirm.interaction.response.edit_message(embed=mangaRemoved, view=None)

# manga setgroup views
class SelectMangaSetGroupView(discord.ui.View):
    def __init__(self, manga_list, mode):
        super().__init__(timeout=15.0)
        self.select_manga = SelectMangaSetGroup(manga_list=manga_list, mode=mode)
        self.add_item(self.select_manga)

class SelectMangaSetGroup(discord.ui.Select):
    def __init__(self, manga_list, mode):
        manga_desc = []
        for manga in manga_list:
            manga_desc.append(discord.SelectOption(label=manga["dropdownTitle"]))
        super().__init__(
            placeholder="Choose a manga series...",
            min_values=1,
            max_values=1,
            options=manga_desc
        )
        self.manga_list = manga_list
        self.finish = None
        self.modeval = mode

    async def callback(self, interaction: discord.Interaction):
        value = int(self.values[0][:2].replace(".", "")) - 1
        mangaid = self.manga_list[value]["id"]
        title = self.manga_list[value]["title"]
        series_groups = await mangaupdates.series_groups(mangaid)
        i = 1
        description = f"Select the scanlator group you want to set for `{title}`.\n"
        group_list = []
        for group in series_groups["group_list"]:
            name = util.format_group_name(group["name"])
            description += f"{i}. {name}\n"
            group_list.append({"id": group["group_id"], "dropdownTitle": f"{i}. {name}", "name": name, "mangaTitle": title, "mangaid": mangaid})
            i += 1
        groupEmbed = discord.Embed(title="Set Scanlator Group", color=0x3083e3, description=description)
        group_drop = SelectScanGroupView(group_list=group_list, mode=self.modeval)
        await interaction.response.edit_message(embed=groupEmbed, view=group_drop)
        self.finish = True
        await group_drop.wait()
        if group_drop.select_group.finish is None:
            await interaction.message.edit(embed=timeoutError, view=None)
            return
        else:
            return

class SelectScanGroupView(discord.ui.View):
    def __init__(self, group_list, mode):
        super().__init__(timeout=15.0)
        self.select_group = SelectScanGroup(group_list=group_list, mode=mode)
        self.add_item(self.select_group)

class SelectScanGroup(discord.ui.Select):
    def __init__(self, group_list, mode):
        group_desc = []
        for group in group_list:
            group_desc.append(discord.SelectOption(label=group["dropdownTitle"]))
        super().__init__(
            placeholder="Choose a scanlator group...",
            min_values=1,
            max_values=1,
            options=group_desc
        )
        self.group_list = group_list
        self.finish = None
        self.modeval = mode

    async def callback(self, interaction: discord.Interaction):
        value = int(self.values[0][:2].replace(".", "")) - 1
        name = self.group_list[value]["name"]
        mangaTitle = self.group_list[value]["mangaTitle"]
        result = discord.Embed(title=f"Confirm", description=f"Are you sure you want to set `{name}` as the scanlator group for `{mangaTitle}`?", color=0x3083e3)
        confirm = Confirm()
        await interaction.response.edit_message(embed=result, view=confirm)
        mangaid = self.group_list[value]["mangaid"]
        groupid = self.group_list[value]["id"]
        self.finish = True
        await confirm.wait()
        if confirm.value is None:
            await interaction.message.edit(embed=timeoutError, view=None)
            return
        elif confirm.value is False:
            cancelEmbed = discord.Embed(title=f"Canceled", color=0x3083e3, description="Successfully canceled.")
            await confirm.interaction.response.edit_message(embed=cancelEmbed, view=None)
            return
        else:
            if self.modeval == "user":
                await mongo.set_scan_group_user(confirm.interaction.user.id, mangaid, groupid, name)
            elif self.modeval == "server":
                await mongo.set_scan_group_server(confirm.interaction.guild.id, mangaid, groupid, name)
            completeEmbed = discord.Embed(title="Successfully Set Scanlator Group", color=0x3083e3, description=f"Scanlator group for `{mangaTitle}` has been set to `{name}`.")
            await confirm.interaction.response.edit_message(embed=completeEmbed, view=None)

class MangaMain(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
    
    manga = SlashCommandGroup("manga", "Manga related commands")

    @manga.command(name="add", description="Adds manga to your list to be tracked")
    async def add(self, ctx, manga: Option(str, description="The name of the manga series (can use mangaupdates links)", required=True)):
        if ctx.guild is not None:
            modeEmbed = discord.Embed(title="Add Manga", color=0x3083e3, description="Do you want this manga added to your list or this server's list?")
            mode = Mode()
            await ctx.respond(embed=modeEmbed, view=mode)
            await mode.wait()
            if mode.value is None:
                await mode.interaction.response.edit_message(embed=timeoutError, view=None)
            else:
                modeval = mode.value
        else:
            modeval = "user"
            mode = None
        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the setup command first.")
        if modeval == "user":
            userExist = await mongo.check_user_exist(ctx.author.id)
            if userExist is False:
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=setupError, view=None)
                else:
                    await ctx.respond(embed=setupError, view=None)
                return
        elif modeval == "server":
            serverExist = await mongo.check_server_exist(ctx.guild.id)
            if ctx.author.guild_permissions.administrator is False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to add manga. You need `Administrator` permission to use this.")
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=permissionError, view=None)
                else:
                    await ctx.respond(embed=permissionError, view=None)
                return
            else:
                if serverExist is False:
                    if mode is not None:
                        await mode.interaction.response.edit_message(embed=setupError, view=None)
                    else:
                        await ctx.respond(embed=setupError, view=None)
                    return
        if validators.url(manga) is True:
            link = manga.partition("https://www.mangaupdates.com/series/")[2]
            mangaid = link.partition("/")[0]
            mangaid = await mangaupdates.convert_new_id(mangaid)
            series_info = await mangaupdates.series_info(mangaid)
            manganame = series_info["title"]
            if modeval == "user":
                mangaindb = await mongo.check_manga_exist_user(ctx.author.id, mangaid)
            elif modeval == "server":
                mangaindb = await mongo.check_manga_exist_server(ctx.guild.id, mangaid)
            if mangaindb is True:
                mangaExist = discord.Embed(title="Add Manga", color=0x3083e3, description="This manga is already added to your list.")
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=mangaExist, view=None)
                else:
                    await ctx.respond(embed=mangaExist, view=None)
                return
            elif mangaindb is False:
                if modeval == "user":
                    await mongo.add_manga_user(ctx.author.id, mangaid, manganame)
                elif modeval == "server":
                    await mongo.add_manga_server(ctx.guild.id, mangaid, manganame)
                mangaAdded = discord.Embed(title="Add Manga", color=0x3083e3, description="Manga succesfully added.")
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=mangaAdded, view=None)
                else:
                    await ctx.respond(embed=mangaAdded, view=None)
        elif validators.url(manga) is not True:
            search_results = []
            description = "Select the manga you want to add to your list.\n"
            search_data = await mangaupdates.search_series(manga)
            if search_data["results"] == []:
                resultError = discord.Embed(title="Error", color=0xff4f4f, description="No mangas were found.")
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=resultError)
                else:
                    await ctx.respond(embed=resultError)
                return
            elif search_data["results"] != []:
                count = 1
                for manga_data in search_data["results"]:
                    manga = manga_data["record"]
                    name = manga["title"]
                    year = manga["year"]
                    rating = manga["bayesian_rating"]
                    description += f"{count}. {name} ({year}, Rating: {rating})\n"
                    search_results.append({"id": manga["series_id"], "name": f"{count}. {name}", "description": f"{year}, Rating: {rating}", "info": manga})
                    count += 1
                searchEmbed = discord.Embed(title="Search Results", color=0x3083e3, description=description)
                manga_drop = SelectMangaView(manga_list=search_results, mode=modeval)
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=searchEmbed, view=manga_drop)
                else:
                    search = await ctx.respond(embed=searchEmbed, view=manga_drop)
                await manga_drop.wait()
                if manga_drop.select_manga.finish is None:
                    if mode is not None:
                        await mode.interaction.message.edit(embed=timeoutError, view=None)
                    else:
                        await search.edit(embed=timeoutError, view=None)
                    return
                else:
                    return
        else:
            completeError = discord.Embed(title="Error", color=0xff4f4f, description=f"Something went wrong. Create an issue here for support: https://github.com/{ghuser}/mangaupdates-bot")
            await mode.interaction.response.edit_message(embed=completeError)
        
    @manga.command(name="remove", description="Removes a manga series from your list")
    async def remove(self, ctx):
        if ctx.guild is not None:
            modeEmbed = discord.Embed(title="Remove Manga", color=0x3083e3, description="Do you want to remove a manga from your list or this server's list?")
            mode = Mode()
            await ctx.respond(embed=modeEmbed, view=mode)
            await mode.wait()
            if mode.value is None:
                await mode.interaction.response.edit_message(embed=timeoutError, view=None)
            else:
                modeval = mode.value
        else:
            modeval = "user"
            mode = None
        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the setup command first.")
        noManga = discord.Embed(title="Error", color=0xff4f4f, description="You have no manga added to your list. Please add some manga first.")
        if modeval == "user":
            userExist = await mongo.check_user_exist(ctx.author.id)
            if userExist is False:
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=setupError, view=None)
                else:
                    await ctx.respond(embed=setupError, view=None)
                return
            mangaList = await mongo.get_manga_list_user(ctx.author.id)
            if mangaList is None:
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=noManga, view=None)
                else:
                    await ctx.respond(embed=noManga, view=None)
                return
        elif modeval == "server":
            if ctx.author.guild_permissions.administrator is False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to remove manga. You need `Administrator` permission to use this.")
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=permissionError, view=None)
                else:
                    await ctx.respond(embed=permissionError, view=None)
                return
            else:
                serverExist = await mongo.check_server_exist(ctx.guild.id)
                if serverExist is False:
                    if mode is not None:
                        await mode.interaction.response.edit_message(embed=setupError, view=None)
                    else:
                        await ctx.respond(embed=setupError, view=None)
                    return
                mangaList = await mongo.get_manga_list_server(ctx.guild.id)
                if mangaList is None:
                    if mode is not None:
                        await mode.interaction.response.edit_message(embed=noManga, view=None)
                    else:
                        await ctx.respond(embed=noManga, view=None)
                    return
        i = 1
        description = "Select the manga you want to remove.\n"
        manga_list = []
        for manga in mangaList:
            description += f"{i}. {manga['title']}\n"
            manga_list.append({"id": manga["id"], "dropdownTitle": f"{i}. {manga['title']}", "title": manga["title"]})
            i += 1
        removeEmbed = discord.Embed(title="Remove Manga", color=0x3083e3, description=description)
        manga_drop = SelectMangaRemoveView(manga_list=manga_list, mode=modeval)
        if mode is not None:
            await mode.interaction.response.edit_message(embed=removeEmbed, view=manga_drop)
        else:
            remove = await ctx.respond(embed=removeEmbed, view=manga_drop)
        await manga_drop.wait()
        if manga_drop.select_manga.finish is None:
            if mode is not None:
                await mode.interaction.message.edit(embed=timeoutError, view=None)
            else:
                await remove.edit(embed=timeoutError, view=None)
            return
        else:
            return

    @manga.command(name="list", description="Lists all manga in your list")
    async def list(self, ctx):
        if ctx.guild is not None:
            modeEmbed = discord.Embed(title="Manga List", color=0x3083e3, description="Do you want to see your manga list or this server's manga list?")
            mode = Mode()
            await ctx.respond(embed=modeEmbed, view=mode)
            await mode.wait()
            if mode.value is None:
                await mode.interaction.response.edit_message(embed=timeoutError, view=None)
            else:
                modeval = mode.value
        else:
            modeval = "user"
            mode = None
        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the setup command first.")
        if modeval == "user":
            userExist = await mongo.check_user_exist(ctx.author.id)
            if userExist is False:
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=setupError, view=None)
                else:
                    await ctx.respond(embed=setupError, view=None)
                return
            name = ctx.author.name
            icon = ctx.author.display_avatar
            mangaList = await mongo.get_manga_list_user(ctx.author.id)
        elif modeval == "server":
            serverExist = await mongo.check_server_exist(ctx.guild.id)
            if serverExist is False:
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=setupError, view=None)
                else:
                    await ctx.respond(embed=setupError, view=None)
                return
            name = ctx.guild.name
            if ctx.guild.icon is not None:
                    icon = ctx.guild.icon.url
            else:
                icon = "https://cdn.discordapp.com/embed/avatars/0.png"
            mangaList = await mongo.get_manga_list_server(ctx.guild.id)
        description = ""
        if mangaList is None:
            description="You have no manga added to your list."
        else:
            for manga in mangaList:
                description += f"• {manga['title']}\n"
        mangaListEmbed = discord.Embed(title=f"{name}'s Manga List", color=0x3083e3, description=description)
        mangaListEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        mangaListEmbed.set_thumbnail(url=icon)
        if mode is not None:
            await mode.interaction.response.edit_message(embed=mangaListEmbed, view=None)
        else:
            await ctx.respond(embed=mangaListEmbed, view=None)

    @manga.command(name="setgroup", description="Sets a manga's scan group")
    async def setgroup(self, ctx):
        if ctx.guild is not None:
            modeEmbed = discord.Embed(title="Set Scanlator Group", color=0x3083e3, description="Do you want to set your manga's scan groups or the server's scan groups?")
            mode = Mode()
            await ctx.respond(embed=modeEmbed, view=mode)
            await mode.wait()
            if mode.value is None:
                await mode.interaction.response.edit_message(embed=timeoutError, view=None)
            else:
                modeval = mode.value
        else:
            modeval = "user"
            mode = None
        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the setup command first.")
        noManga = discord.Embed(title="Error", color=0xff4f4f, description="You have no manga added to your list. Please add some manga first.")
        if modeval == "user":
            userExist = await mongo.check_user_exist(ctx.author.id)
            if userExist is False:
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=setupError, view=None)
                else:
                    await ctx.respond(embed=setupError, view=None)
                return
            mangaList = await mongo.get_manga_list_user(ctx.author.id)
            if mangaList is None:
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=noManga, view=None)
                else:
                    await ctx.respond(embed=noManga, view=None)
                return
        elif modeval == "server":
            if ctx.author.guild_permissions.administrator is False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to remove manga. You need `Administrator` permission to use this.")
                if mode is not None:
                    await mode.interaction.response.edit_message(embed=permissionError, view=None)
                else:
                    await ctx.respond(embed=permissionError, view=None)
                return
            else:
                serverExist = await mongo.check_server_exist(ctx.guild.id)
                if serverExist is False:
                    if mode is not None:
                        await mode.interaction.response.edit_message(embed=setupError, view=None)
                    else:
                        await ctx.respond(embed=setupError, view=None)
                    return
                mangaList = await mongo.get_manga_list_server(ctx.guild.id)
                if mangaList is None:
                    if mode is not None:
                        await mode.interaction.response.edit_message(embed=noManga, view=None)
                    else:
                        await ctx.respond(embed=noManga, view=None)
                    return
        i = 1
        description = "Select the manga you want to set the scanlator group for.\n"
        manga_list = []
        for manga in mangaList:
            description += f"{i}. {manga['title']}\n"
            manga_list.append({"id": manga["id"], "dropdownTitle": f"{i}. {manga['title']}", "title": manga["title"]})
            i += 1
        selectMangaEmbed = discord.Embed(title="Set Scanlator Group", color=0x3083e3, description=description)
        manga_drop = SelectMangaSetGroupView(manga_list=manga_list, mode=modeval)
        if mode is not None:
            await mode.interaction.response.edit_message(embed=selectMangaEmbed, view=manga_drop)
        else:
            search = await ctx.respond(embed=selectMangaEmbed, view=manga_drop)
        await manga_drop.wait()
        if manga_drop.select_manga.finish is None:
            if mode is not None:
                await mode.interaction.message.edit(embed=timeoutError, view=None)
            else:
                await search.edit(embed=timeoutError, view=None)
            return
        else:
            return

    @manga.command(name="test", description="Tests sending manga to your server")
    async def testsending(self, ctx):
        if ctx.guild is not None:
            channel = await mongo.get_channel(ctx.guild.id)
            if channel is None:
                noChannelError = discord.Embed(title="Error", color=0xff4f4f, description="You have no channel set up. Please set one up with the `setup` command.")
                await ctx.respond(embed=noChannelError, view=None)
                return
            else:
                try:
                    channelObject = self.bot.get_channel(channel)
                    if channelObject is None:
                        noChannelError = discord.Embed(title="Error", color=0xff4f4f, description="Channel doesn't exist anymore, please reset the channel with the `server setchannel` command.")
                        await ctx.respond(embed=noChannelError, view=None)
                        return
                    mangaTestEmbed = discord.Embed(title=f"Testing Update", url="https://picsiv.hayasaka.moe/", description=f"This is a test alert for the MangaUpdates Bot.", color=0x3083e3)
                    mangaTestEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
                    mangaTestEmbed.add_field(name="Chapter", value="c.1", inline=True)
                    mangaTestEmbed.add_field(name="Group", value="The MangaUpdates Bot Team", inline=True)
                    mangaTestEmbed.add_field(name="Scanlator Link", value="https://picsiv.hayasaka.moe/", inline=False)
                    await channelObject.send(embed=mangaTestEmbed, view=None)
                    mainResponseEmbed = discord.Embed(title="Test Manga Update Message", color=0x3083e3, description="Congrats! Manga can be sent to your server.")
                    await ctx.respond(embed=mainResponseEmbed, view=None)
                except:
                    mainResponseEmbed = discord.Embed(title="Test Manga Update Message", color=0xff4f4f, description="Oops! Something went wrong. Give the bot permissions to send messages in the channel you have set and then try again.")
                    await ctx.respond(embed=mainResponseEmbed, view=None)
                    return   
        else:
            notServerError = discord.Embed(title="Error", color=0xff4f4f, description="You need to be in a server to use this command.")
            await ctx.respond(embed=notServerError, view=None)
            return
        
def setup(bot):
    bot.add_cog(MangaMain(bot))