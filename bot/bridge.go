package bot

import (
	log "github.com/Sirupsen/logrus"

	"github.com/GinjaNinja32/DisGoIRC/format"
)

// Config requires the required config to connect to IRC/Discord and the mapping between them
type Config struct {
	IRC     IRCConfig         `json:"irc"`
	Discord DiscordConfig     `json:"discord"`
	Mapping map[string]string `json:"mapping"`
}

var (
	conf           Config
	inverseMapping map[string]string
)

// Init starts the bridge with the given config
func Init(c Config) {
	conf = c
	inverseMapping = map[string]string{}
	for k, v := range conf.Mapping {
		inverseMapping[v] = k
	}
	dInit()
	iInit()
}

func incomingIRC(nick, channel, message string) {
	log.Infof("IRC %s <%s> %s", channel, nick, message)

	discordChan, ok := conf.Mapping[channel]
	if !ok {
		return
	}

	log.Debugf("Mapping IRC:%s to DIS:%s", channel, discordChan)

	fs := format.ParseIRC(message)

	dOutgoing(nick, discordChan, fs)
}

func incomingDiscord(nick, channel, message string) {
	log.Infof("DIS %s <%s> %s", channel, nick, message)

	ircChan, ok := inverseMapping[channel]
	if !ok {
		return
	}

	log.Debugf("Mapping DIS:%s to IRC:%s", channel, ircChan)

	fs := format.ParseDiscord(message)

	iOutgoing(nick, ircChan, fs)
}
