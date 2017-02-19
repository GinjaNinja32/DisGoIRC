package bot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	discord "github.com/bwmarrin/discordgo"
)

type DiscordConfig struct {
	Token string "token"
}

var (
	dBotID      string
	dSession    *discord.Session
	dGuilds     = map[string]string{}
	dGuildChans = map[string]map[string]string{}
)

func dInit() {
	d, err := discord.New(fmt.Sprintf("Bot %s", conf.Discord.Token))
	dSession = d
	if err != nil {
		log.Fatalf("Failed to initialise Discord session: %s", err)
	}

	u, err := dSession.User("@me")
	if err != nil {
		log.Fatalf("Failed to get own Discord user: %s", err)
	}

	dBotID = u.ID

	guilds, err := dSession.UserGuilds()
	if err != nil {
		log.Fatalf("Failed to get guilds: %s", err)
	}

	for _, g := range guilds {
		chans, err := dSession.GuildChannels(g.ID)
		if err != nil {
			log.Fatalf("Failed to get channels for %s: %s", g.Name, err)
		}

		dGuilds[g.Name] = g.ID
		dGuildChans[g.Name] = map[string]string{}
		for _, c := range chans {
			if c.Type == "text" {
				dGuildChans[g.Name][c.Name] = c.ID
			}
		}
	}

	dSession.AddHandler(dMessageCreate)

	err = dSession.Open()
	if err != nil {
		log.Fatalf("Failed to connect to Discord: %s", err)
	}

	log.Infof("Connected to Discord")
}

func dMessageCreate(s *discord.Session, m *discord.MessageCreate) {
	if m.Author.ID == dBotID {
		return
	}

	c, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Errorf("Failed to get channel for incoming message with CID %s: %s", m.ChannelID, err)
		return
	}

	guildID := c.GuildID

	g, err := s.Guild(guildID)
	if err != nil {
		log.Errorf("Failed to get guild with ID %s: %s", guildID, err)
		return
	}

	channel := fmt.Sprintf("%s#%s", g.Name, c.Name)

	if m.Content != "" {
		message := m.Content

		// Channels
		for _, c := range g.Channels {
			if c.Type != "text" {
				continue
			}
			find := fmt.Sprintf("<#%s>", c.ID)
			replace := fmt.Sprintf("#%s", c.Name)
			message = strings.Replace(message, find, replace, -1)
		}

		// Users
		for _, u := range g.Members {
			find := fmt.Sprintf("<@%s>", u.User.ID)
			find2 := fmt.Sprintf("<@!%s>", u.User.ID)
			replace := fmt.Sprintf("@%s", u.User.Username)
			message = strings.Replace(message, find, replace, -1)
			message = strings.Replace(message, find2, replace, -1)
		}

		// Roles
		for _, r := range g.Roles {
			find := fmt.Sprintf("<@&%s>", r.ID)
			replace := fmt.Sprintf("@%s", r.Name)
			message = strings.Replace(message, find, replace, -1)
		}

		// Multiline
		lines := strings.Split(message, "\n")
		if len(lines) > 3 {
			url := uploadToPtpb(message)

			for _, line := range lines[:2] {
				incomingDiscord(m.Author.Username, channel, line)
			}
			incomingDiscord("[SYSTEM]", channel, fmt.Sprintf("full message from %s: %s", iAddAntiPing(m.Author.Username), url))
		} else {
			for _, line := range lines {
				incomingDiscord(m.Author.Username, channel, line)
			}
		}

		//incomingDiscord(m.Author.Username, channel, message)
	}
	for _, a := range m.Attachments {
		incomingDiscord(m.Author.Username, channel, a.ProxyURL)
	}
}

func uploadToPtpb(s string) string {
	resp, err := http.PostForm("https://ptpb.pw/",
		url.Values{"c": {s}, "p": {"1"}})
	defer resp.Body.Close()

	if err != nil {
		log.Errorf("Failed to upload to PTPB: %s", err)
		return "Failed to upload to PTPB"
	}
	if resp.StatusCode == http.StatusOK {
		return resp.Header.Get("Location")
	}

	log.Errorf("Failed to upload to PTPB: HTTP %d", resp.StatusCode)
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read body: %s", err)
	} else {
		log.Errorf(string(ret))
	}
	return fmt.Sprintf("Failed to upload to PTPB: HTTP %d", resp.StatusCode)

}

func dOutgoing(nick, channel, message string) {
	chanParts := strings.Split(channel, "#")
	guildID := dGuilds[chanParts[0]]
	chanID := dGuildChans[chanParts[0]][chanParts[1]]

	g, err := dSession.Guild(guildID)
	if err != nil {
		log.Errorf("Failed to get guild with ID %s: %s", guildID, err)
		return
	}

	// Channels
	for _, c := range g.Channels {
		if c.Type != "text" {
			continue
		}
		find := fmt.Sprintf("#%s", c.Name)
		replace := fmt.Sprintf("<#%s>", c.ID)
		message = strings.Replace(message, find, replace, -1)
	}

	// Users
	for _, u := range g.Members {
		find := fmt.Sprintf("@%s", u.User.Username)
		replace := fmt.Sprintf("<@%s>", u.User.ID)
		message = strings.Replace(message, find, replace, -1)
	}

	// Roles
	for _, r := range g.Roles {
		find := fmt.Sprintf("@%s", r.Name)
		replace := fmt.Sprintf("<@&%s>", r.ID)
		message = strings.Replace(message, find, replace, -1)
	}

	dSession.ChannelMessageSend(chanID, fmt.Sprintf("**<%s>** %s", nick, message))
}
