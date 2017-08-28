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

// DiscordConfig represents the required config to connect to Discord
type DiscordConfig struct {
	Token        string `json:"token"`
	UseNicknames bool   `json:"use_nicknames"`
}

var (
	dBotID      string
	dSession    *discord.Session
	dGuilds     = map[string]string{}
	dGuildChans = map[string]map[string]string{}

	dMsgQueue = make(chan func())
)

const (
	msgTypeText = "text"
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
		var chans []*discord.Channel
		chans, err = dSession.GuildChannels(g.ID)
		if err != nil {
			log.Fatalf("Failed to get channels for %s: %s", g.Name, err)
		}

		dGuilds[g.Name] = g.ID
		dGuildChans[g.Name] = map[string]string{}
		for _, c := range chans {
			if c.Type == msgTypeText {
				dGuildChans[g.Name][c.Name] = c.ID
			}
		}
	}

	dSession.AddHandler(dMessageCreate)

	err = dSession.Open()
	if err != nil {
		log.Fatalf("Failed to connect to Discord: %s", err)
	}

	go func() {
		for f := range dMsgQueue {
			f()
		}
	}()

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
	authorName := getDisplayNameForUser(m.Author, g.Members)

	if m.Content != "" {
		message := convertMentionsForIRC(g, m)

		dispatchMessageToIRC(authorName, channel, message)
	}
	for _, a := range m.Attachments {
		incomingDiscord(authorName, channel, a.ProxyURL)
	}
}

func convertMentionsForIRC(g *discord.Guild, m *discord.MessageCreate) string {
	message := m.Content

	// Channels
	for _, c := range g.Channels {
		if c.Type != msgTypeText {
			continue
		}
		find := fmt.Sprintf("<#%s>", c.ID)
		replace := fmt.Sprintf("#%s", c.Name)
		message = strings.Replace(message, find, replace, -1)
	}

	// Users
	for _, u := range g.Members {
		display := getDisplayNameForMember(u)
		if display == "" {
			log.Errorf("%s/%q/%q had an invalid display name", u.User.ID, u.User.Username, u.Nick)
			continue
		}
		find := fmt.Sprintf("<@%s>", u.User.ID)
		find2 := fmt.Sprintf("<@!%s>", u.User.ID)
		replace := fmt.Sprintf("@%s", display)
		message = strings.Replace(message, find, replace, -1)
		message = strings.Replace(message, find2, replace, -1)
	}

	// Roles
	for _, r := range g.Roles {
		find := fmt.Sprintf("<@&%s>", r.ID)
		replace := fmt.Sprintf("@%s", r.Name)
		message = strings.Replace(message, find, replace, -1)
	}

	return message
}

func dispatchMessageToIRC(authorName, channel, message string) {
	// Multiline
	lines := strings.Split(message, "\n")
	lines, forceClip := clipLinesForIRC(lines)
	if len(lines) > 3 || forceClip {
		url := uploadToPtpb(message)

		n := 2
		if len(lines) < 2 {
			n = len(lines)
		}

		for _, line := range lines[:n] {
			incomingDiscord(authorName, channel, line)
		}
		incomingDiscord("[SYSTEM]", channel, fmt.Sprintf("full message from %s: %s", iAddAntiPing(authorName), url))
	} else {
		for _, line := range lines {
			incomingDiscord(authorName, channel, line)
		}
	}
}

func clipLinesForIRC(s []string) ([]string, bool) {
	ret := []string{}
	anyLineForceClip := false

	for _, line := range s {
		if len(line) < 300 {
			ret = append(ret, line)
		} else {
			words := strings.Split(line, " ")
			for len(words) != 0 {
				l := words[0]
				words = words[1:]
				for len(words) != 0 && len(l)+len(words[0]) < 300 {
					l = l + " " + words[0]
					words = words[1:]
				}

				anyLineForceClip = anyLineForceClip || len(l) > 300
				ret = append(ret, l)
			}
		}
	}

	return ret, anyLineForceClip
}

func uploadToPtpb(s string) string {
	resp, err := http.PostForm("https://ptpb.pw/",
		url.Values{"c": {s}, "p": {"1"}})

	if err != nil {
		log.Errorf("Failed to upload to PTPB: %s", err)
		return "Failed to upload to PTPB"
	}
	defer func() {
		_ = resp.Body.Close()
	}()
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
		if c.Type != msgTypeText {
			continue
		}
		find := fmt.Sprintf("#%s", c.Name)
		replace := fmt.Sprintf("<#%s>", c.ID)
		message = strings.Replace(message, find, replace, -1)
	}

	// Users
	for _, u := range g.Members {
		display := getDisplayNameForMember(u)
		if display == "" {
			log.Errorf("%s/%q/%q had an invalid display name", u.User.ID, u.User.Username, u.Nick)
			continue
		}
		find := fmt.Sprintf("@%s", display)
		replace := fmt.Sprintf("<@%s>", u.User.ID)
		message = strings.Replace(message, find, replace, -1)
	}

	// Roles
	for _, r := range g.Roles {
		find := fmt.Sprintf("@%s", r.Name)
		replace := fmt.Sprintf("<@&%s>", r.ID)
		message = strings.Replace(message, find, replace, -1)
	}

	dMsgQueue <- func() {
		_, err := dSession.ChannelMessageSend(chanID, fmt.Sprintf("**<%s>** %s", nick, message))
		if err != nil {
			log.Errorf("Failed to send message to %s: <%s> %s", chanID, nick, message)
		}
	}
}

func getDisplayNameForMember(member *discord.Member) string {
	if conf.Discord.UseNicknames && member.Nick != "" {
		return member.Nick
	}

	return member.User.Username
}

func getDisplayNameForUser(user *discord.User, members []*discord.Member) string {
	if conf.Discord.UseNicknames {
		for _, m := range members {
			if m.User.ID == user.ID {
				return getDisplayNameForMember(m)
			}
		}
	}

	return user.Username
}
