package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store"
	"github.com/bwmarrin/discordgo"
)

type GatewayCmd struct{}

func (c *GatewayCmd) Run(store *store.Store) error {
	bot, err := connectBot(store)
	if err != nil {
		return err
	}

	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		channelID, err := store.GetChannelID()
		if err != nil {
			log.Println("get channel ID:", err)
			return
		} else if i.ChannelID != channelID {
			return
		}

		if i.Type == discordgo.InteractionMessageComponent {
			data := i.MessageComponentData()
			lastScheduleMessageID, err := store.GetLastScheduleMessageID()
			if err != nil {
				logAndReportError(err, "get lsmID", s, i.Interaction)
				return
			}
			if lastScheduleMessageID != i.Message.ID {
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You can’t change your schedule in the past! Try using today’s signup: https://discord.com/channels/" + i.GuildID + "/" + i.ChannelID + "/" + lastScheduleMessageID,
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}); err != nil {
					log.Println("interaction response:", err)
				}
				return
			}

			hours := make([]int, 0, len(data.Values))

			switch data.CustomID {
			case "signup_selection":
				for _, hourString := range data.Values {
					hour, err := strconv.Atoi(hourString)
					if err != nil {
						logAndReportError(err, "atoi", s, i.Interaction)
						return
					}
					hours = append(hours, hour)
				}
			case "remove_signup":
				// leave hours empty
			default:
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "The `" + data.CustomID + "` component you have interacted with is invalid. You should probably tell <@311598715817426945> so she can fix it.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}); err != nil {
					log.Println("interaction response:", err)
				}
				return
			}

			err = store.SetUserArrivals(i.Member.User.ID, hours)
			if err != nil {
				logAndReportError(err, "set user arrivals", s, i.Interaction)
				return
			}

			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			}); err != nil {
				log.Println("interaction response:", err)
			}

			signupEmbedIndex, err := store.GetSignupEmbedIndex()
			if err != nil {
				log.Println("get signup embed index:", err)
				return
			}

			if signupEmbedIndex < len(i.Message.Embeds) {
				embeds := i.Message.Embeds
				signupEmbed := embeds[signupEmbedIndex]
				arrivals, err := store.GetAllArrivals()
				if err != nil {
					log.Println("get all arrivals:", err)
				}
				if len(arrivals) == 0 {
					signupEmbed.Description = "No one has signed up yet!"
					signupEmbed.Fields = nil
				} else {
					signupEmbed.Description = ""
					signupEmbed.Fields = nil
					for _, hour := range arrivals {
						signupEmbed.Fields = append(signupEmbed.Fields, &discordgo.MessageEmbedField{
							Name:  hourNames[hour.Hour],
							Value: "<@" + strings.Join(hour.ArrivingUsers, ">\n<@") + ">",
						})
					}
				}
				_, err = s.ChannelMessageEditEmbeds(i.ChannelID, i.Message.ID, embeds)
				if err != nil {
					log.Println("edit message:", err)
				}
			}
		}
	})

	err = bot.Open()
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	log.Println("Gateway connected")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
	log.Println("Closing gateway connection")
	return bot.Close()
}

func logAndReportError(err error, context string, s *discordgo.Session, interaction *discordgo.Interaction) {
	errMessage := fmt.Sprintf("%s: %s", context, err)
	log.Println(errMessage)
	if err := s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "An error occurred. You should probably tell <@311598715817426945>.\n`" + errMessage + "`",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Println("interaction response:", err)
	}
}
