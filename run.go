package brook

import "github.com/txthinking/brook/sysproxy"

// RunServer used to make a new Server and start to listen
func RunServer(address, password string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	s, err := NewServer(address, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}

// RunClient used to make a new Client and start a socks5 proxy to listen
func RunClient(address, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewClient(address, ip, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServe()
}

// RunClientAsHTTP used to make a new Client and start a http proxy to listen
func RunClientAsHTTP(address, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewClient(address, ip, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServeHTTP()
}

// RunTunnel used to start a tunnel
func RunTunnel(address, to, server, password string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	c, err := NewTunnel(address, to, server, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return err
	}
	return c.ListenAndServe()
}

// RunStreamServer used to make a new StreamServer and start to listen
func RunStreamServer(address, password string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	s, err := NewStreamServer(address, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}

// RunStreamClient used to make a new StreamClient and start a socks5 proxy to listen
func RunStreamClient(address, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewStreamClient(address, ip, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServe()
}

// RunClientAsHTTP used to make a new StreamClient and start a http proxy to listen
func RunStreamClientAsHTTP(address, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewStreamClient(address, ip, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServeHTTP()
}

// RunSSServer used to make a new Server and start to listen
func RunSSServer(address, password string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	s, err := NewSSServer(address, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}

// RunSSClient used to make a new Client and start a socks5 proxy to listen
func RunSSClient(address, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewSSClient(address, ip, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServe()
}

// RunSSClientAsHTTP used to make a new Client and start a http proxy to listen
func RunSSClientAsHTTP(address, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewSSClient(address, ip, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServeHTTP()
}

// RunRelay used to make a new Relay and start to listen
func RunRelay(address, remote string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	r, err := NewRelay(address, remote, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return err
	}
	return r.ListenAndServe()
}

// RunSocks5Server used to make a new Socks5Server and start a raw socks5 proxy to listen
func RunSocks5Server(address, ip, username, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	s, err := NewSocks5Server(address, ip, username, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}

// RunSocks5ToHTTP used to make a new Socks5ToHTTP and start a http proxy to listen
func RunSocks5ToHTTP(address, socks5 string, timeout, deadline int) error {
	s, err := NewSocks5ToHTTP(address, socks5, timeout, deadline)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}

// RunSystemProxy used to set/remove system proxy
func RunSystemProxy(remove bool, pac string) error {
	if remove {
		if err := sysproxy.TurnOffSystemProxy(); err != nil {
			return err
		}
		return nil
	}
	if err := sysproxy.TurnOnSystemProxy(pac); err != nil {
		return err
	}
	return nil
}
