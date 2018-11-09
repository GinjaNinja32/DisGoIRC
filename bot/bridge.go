package bot

import (
	"strings"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"

	"github.com/GinjaNinja32/DisGoIRC/format"
)

// Config requires the required config to connect to IRC/Discord and the mapping between them
type Config struct {
	IRC     IRCConfig         `json:"irc"`
	Discord DiscordConfig     `json:"discord"`
	Mapping map[string]string `json:"mapping"`
}

var (
	conf            Config
	inverseMapping  map[string]string
	modifiedMapping map[string]string
)

// Init starts the bridge with the given config
func Init(c Config) {
	conf = c
	inverseMapping = map[string]string{}
	modifiedMapping = map[string]string{}
	for k, v := range conf.Mapping {
		ircChannelPassword := strings.Split(k, " ") // "#channel password" -> ["#channel", "password"]
		ircChannel := ircChannelPassword[0]
		inverseMapping[v] = ircChannel
		modifiedMapping[ircChannel] = v
	}
	dInit()
	iInit()
}

// hasCommand checks for the existence of the configured command characters at the start of a message
func hasCommand(message, commandChars string) bool {
	firstRune, _ := utf8.DecodeRuneInString(message)
	return firstRune != 0 && strings.ContainsRune(commandChars, firstRune)
}

// incomingIRC is called on every message from a mapped IRC channel and posts it to the configured Discord channel
func incomingIRC(nick, channel, message string) {
	log.Infof("IRC %s <%s> %s", channel, nick, message)

	discordChan, ok := modifiedMapping[channel]
	if !ok {
		return
	}

	log.Debugf("Mapping IRC:%s to DIS:%s", channel, discordChan)

	fs := format.ParseIRC(message)

	if hasCommand(message, conf.IRC.CommandChars) {
		dOutgoing(nick, discordChan, format.FormattedString{{Text: "Command sent by " + nick}}, true)
		dOutgoing(nick, discordChan, fs, true)
		return
	}

	dOutgoing(nick, discordChan, fs, false)
}

// incomingDiscord is called on every message from a mapped Discord channel and posts it to the configured IRC channel
func incomingDiscord(nick, channel, message string) {
	log.Infof("DIS %s <%s> %s", channel, nick, message)

	ircChan, ok := inverseMapping[channel]
	if !ok {
		return
	}

	log.Debugf("Mapping DIS:%s to IRC:%s", channel, ircChan)

	fs := format.ParseDiscord(message)

	if hasCommand(message, conf.Discord.CommandChars) {
		iOutgoing(nick, ircChan, format.FormattedString{{Text: "Command sent by " + nick}}, true)
		iOutgoing(nick, ircChan, fs, true)
		return
	}

	iOutgoing(nick, ircChan, fs, false)
}
