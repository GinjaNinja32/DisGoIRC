package bot

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	IRC     IRCConfig         "irc"
	Discord DiscordConfig     "discord"
	Mapping map[string]string "mapping"
}

var (
	conf           Config
	inverseMapping map[string]string
)

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
	log.Debugf("IRC %s <%s> %s", channel, nick, message)

	discordChan, ok := conf.Mapping[channel]
	if !ok {
		return
	}

	log.Debugf("Mapping IRC:%s to DIS:%s", channel, discordChan)

	message = fmtIrcToDiscord(message)

	dOutgoing(nick, discordChan, message)
}

func incomingDiscord(nick, channel, message string) {
	log.Debugf("DIS %s <%s> %s", channel, nick, message)

	ircChan, ok := inverseMapping[channel]
	if !ok {
		return
	}

	log.Debugf("Mapping DIS:%s to IRC:%s", channel, ircChan)

	iOutgoing(nick, ircChan, message)
}

var specialIrc = regexp.MustCompile("|[0-9]{0,2}")

func fmtReplaceInPairs(msg, find, replace string) string {
	r := regexp.MustCompile(find)
	active := false
	msg = r.ReplaceAllStringFunc(msg, func(a string) string {
		active = !active
		return replace
	})
	if active {
		msg = msg + replace
	}
	return msg
}

func fmtIrcToDiscord(msg string) string {
	msg = specialIrc.ReplaceAllString(msg, "")
	msg = fmtReplaceInPairs(msg, "", "**")
	msg = fmtReplaceInPairs(msg, "", "__")
	msg = fmtReplaceInPairs(msg, "", "*")
	return msg
}
