package brook

import (
	"net/url"
	"os"

	"github.com/mdp/qrterminal"
)

// QR generate and print QR code
func QR(server, password string) {
	s := server + " " + password
	s = "brook://" + url.PathEscape(s)
	qrterminal.Generate(s, qrterminal.H, os.Stdout)
}
