package main

import (
	"database/sql"
	"fmt"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store"
	"github.com/bwmarrin/discordgo"
)

type DoctorCmd struct{}

func (*DoctorCmd) Run(store *store.Store) error {
	bot, err := connectBot(store)
	if err != nil {
		return err
	}

	user, err := bot.User("@me")
	if err != nil {
		fmt.Println("Could not login bot:", err)
		return nil
	}
	fmt.Println("Bot auth successful:", user.Username)

	channelID, err := store.GetChannelID()
	if err == sql.ErrNoRows {
		fmt.Println("No channel ID set!")
		return nil
	} else if err != nil {
		fmt.Println("Could not access channel ID:", err)
		return nil
	}

	channel, err := bot.Channel(channelID)
	if err != nil {
		fmt.Println("Could not retrieve channel information:", err)
		return nil
	}
	fmt.Println("Channel:", channel.Name)
	if channel.Type != discordgo.ChannelTypeGuildText {
		fmt.Println("(warning: channel is not a standard text channel)")
	}

	guild, err := bot.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Could not retrieve guild information:", err)
		return nil
	}
	member, err := bot.GuildMember(guild.ID, user.ID)
	if err != nil {
		fmt.Println("Could not retrieve guild member information:", err)
		return nil
	}

	channelPermissions := computeOverwrites(computeBasePermissions(member, guild), member, channel)
	canViewChannel := channelPermissions&discordgo.PermissionViewChannel == discordgo.PermissionViewChannel
	canSendMessages := channelPermissions&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages
	canReadMessageHistory := channelPermissions&discordgo.PermissionReadMessageHistory == discordgo.PermissionReadMessageHistory
	canEmbedLinks := channelPermissions&discordgo.PermissionEmbedLinks == discordgo.PermissionEmbedLinks
	if canViewChannel && canSendMessages && canReadMessageHistory && canEmbedLinks {
		fmt.Println("All necessary channel permissions are available")
	} else {
		fmt.Println("Some necessary channel permissions are missing:")
		if !canViewChannel {
			fmt.Println("- View Channel")
		}
		if !canSendMessages {
			fmt.Println("- Send Messages")
		}
		if !canReadMessageHistory {
			fmt.Println("- Read Message History")
		}
		if !canEmbedLinks {
			fmt.Println("- Embed Links")
		}
	}

	workoutNames, err := store.GetWorkoutNames()
	if err != nil {
		fmt.Println("Could not load workout names:", err)
	} else {
		fmt.Println("Workout names:", workoutNames)
	}

	return nil
}

// https://discord.com/developers/docs/topics/permissions#permission-overwrites

func computeBasePermissions(member *discordgo.Member, guild *discordgo.Guild) int64 {
	if guild.OwnerID == member.User.ID {
		return discordgo.PermissionAll
	}

	roles := indexMap(guild.Roles, func(r *discordgo.Role) string {
		return r.ID
	})

	roleEveryone := roles[guild.ID]
	permissions := roleEveryone.Permissions

	for _, roleID := range member.Roles {
		role := roles[roleID]
		permissions |= role.Permissions
	}

	if permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		return discordgo.PermissionAll
	}

	return permissions
}

func computeOverwrites(basePermissions int64, member *discordgo.Member, channel *discordgo.Channel) int64 {
	if basePermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		return discordgo.PermissionAll
	}

	overwrites := indexMap(channel.PermissionOverwrites, func(o *discordgo.PermissionOverwrite) string {
		return o.ID
	})

	permissions := basePermissions
	overwriteEveryone, ok := overwrites[channel.GuildID]
	if ok {
		permissions &^= overwriteEveryone.Deny
		permissions |= overwriteEveryone.Allow
	}

	var roleAllow, roleDeny int64
	for _, roleID := range member.Roles {
		overwriteRole, ok := overwrites[roleID]
		if ok {
			roleAllow |= overwriteRole.Allow
			roleDeny |= overwriteRole.Deny
		}
	}

	permissions &^= roleDeny
	permissions |= roleAllow

	overwriteMember, ok := overwrites[member.User.ID]
	if ok {
		permissions &^= overwriteMember.Deny
		permissions |= overwriteMember.Allow
	}

	return permissions
}

func indexMap[K comparable, V any](values []V, keyFunc func(V) K) map[K]V {
	result := make(map[K]V, len(values))
	for _, v := range values {
		result[keyFunc(v)] = v
	}

	return result
}
