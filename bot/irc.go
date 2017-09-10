package bot

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	irc "github.com/thoj/go-ircevent"

	"github.com/GinjaNinja32/DisGoIRC/format"
)

// IRCConfig represents the required configuration to connect to IRC
type IRCConfig struct {
	Nick string `json:"nick"`
	User string `json:"user"`
	Pass string `json:"pass"`

	SSL    bool   `json:"ssl"`
	Server string `json:"server"`
}

var (
	iSession *irc.Connection
)

func iInit() {
	c := conf.IRC
	iSession = irc.IRC(c.Nick, c.User)

	iSession.UseTLS = c.SSL
	iSession.TLSConfig = &tls.Config{InsecureSkipVerify: true} // don't verify SSL certs
	iSession.Password = c.Pass
	iSession.AddCallback("PRIVMSG", iPrivmsg)
	iSession.AddCallback("CTCP_ACTION", iAction)

	err := iSession.Connect(c.Server)
	if err != nil {
		log.Fatalf("Failed to initialise IRC session: %s", err)
	}

	iSession.AddCallback("001", iSetupSession)

	log.Infof("Connected to IRC")
}

func iSetupSession(e *irc.Event) {
	for c := range conf.Mapping {
		iSession.Join(c)
	}
}

func iPrivmsg(e *irc.Event) {
	incomingIRC(e.Nick, strings.ToLower(e.Arguments[0]), e.Message())
}
func iAction(e *irc.Event) {
	incomingIRC(e.Nick, strings.ToLower(e.Arguments[0]), fmt.Sprintf("_%s_", e.Message()))
}

var outgoingNickRegex = regexp.MustCompile(`\b[a-zA-Z0-9]`)

func iAddAntiPing(s string) string {
	// add a \uFEFF character to avoid pinging the user
	return outgoingNickRegex.ReplaceAllString(s, "$0\ufeff")
}

func iOutgoing(nick, channel string, message format.FormattedString) {
	nick = iAddAntiPing(nick)
	iSession.Privmsg(channel, fmt.Sprintf("<%s> %s", nick, message.RenderIRC()))
}
