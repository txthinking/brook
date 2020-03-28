module github.com/txthinking/brook

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1

replace github.com/russross/blackfriday/v2 => gopkg.in/russross/blackfriday.v2 v2.0.1

go 1.14

require (
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.2
	github.com/mdp/qrterminal v1.0.1
	github.com/miekg/dns v1.1.29
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/tdewolff/minify v2.3.6+incompatible
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/txthinking/encrypt v0.0.0-20200324035805-5d1a78415440
	github.com/txthinking/gotun2socks v0.0.0-20180829122610-35016fdae05e
	github.com/txthinking/runnergroup v0.0.0-20200327135940-540a793bb997
	github.com/txthinking/socks5 v0.0.0-20200327133705-caf148ab5e9d
	github.com/txthinking/x v0.0.0-20200322173929-c86052c68359
	github.com/urfave/cli v1.22.3
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	gopkg.in/russross/blackfriday.v2 v2.0.0-00010101000000-000000000000 // indirect
)
