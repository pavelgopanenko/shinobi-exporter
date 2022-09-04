package config

import "github.com/urfave/cli/v2"

const (
	EndpointFlag = "shinobi.endpoint"
	TokenFlag    = "shinobi.token"
	GroupToken   = "shinobi.group"
	Insecure     = "shinobi.insecure"

	WebListenAddress     = "web.listen-address"
	DefaultListenAddress = ":9765"
)

func CLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: EndpointFlag, Usage: "shinobi endpoint address", Required: true},
		&cli.StringFlag{Name: TokenFlag, Usage: "shinobi client token", Required: true},
		&cli.StringSliceFlag{Name: GroupToken, Usage: "shinobi groups", Required: true},
		&cli.BoolFlag{Name: Insecure, Usage: "insecure communication", Value: false},

		&cli.StringFlag{Name: WebListenAddress, Usage: "listen network address", Value: DefaultListenAddress},
	}
}
