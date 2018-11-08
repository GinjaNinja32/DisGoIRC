package bot

import (
	"strings"

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
	conf              Config
	inverseMapping    map[string]string
	modifiedMapping   map[string]string
	commandCharacters []string
)

// Init starts the bridge with the given config
func Init(c Config) {
	conf = c
	inverseMapping = map[string]string{}
	modifiedMapping = map[string]string{}
	commandCharacters = strings.Split(conf.IRC.CommandChars, "")
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
func hasCommand(message string) bool {
	hasCommand := false
	splitMessage := strings.Fields(message)
	for _, element := range commandCharacters {
		//Assuming that the command character will always be the first character in the message,
		//as well as there never being any null incoming messages.
		if strings.Contains(strings.Split(splitMessage[0], "")[0], string(element)) {
			hasCommand = true
		}
	}
	return hasCommand
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

	if hasCommand(message) {
		dOutgoing(nick, discordChan, format.ParseIRC("Command sent by "+nick), true)
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

	if hasCommand(message) {
		iOutgoing(nick, ircChan, format.ParseDiscord("Command sent by "+nick), true)
		iOutgoing(nick, ircChan, fs, true)
		return
	}

	iOutgoing(nick, ircChan, fs, false)
}
