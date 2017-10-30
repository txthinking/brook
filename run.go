package brook

// RunServer used to make a new Server and start to listen
func RunServer(address, password string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	s, err := NewServer(address, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}

// RunClient used to make a new Client and start a socks5 proxy to listen
func RunClient(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewClient(address, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServe(nil)
}

// RunClientAsHTTP used to make a new Client and start a http proxy to listen
func RunClientAsHTTP(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewClient(address, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServeHTTP(nil)
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
func RunStreamClient(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewStreamClient(address, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServe(nil)
}

// RunClientAsHTTP used to make a new StreamClient and start a http proxy to listen
func RunStreamClientAsHTTP(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewStreamClient(address, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServeHTTP(nil)
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
func RunSSClient(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewSSClient(address, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServe(nil)
}

// RunSSClientAsHTTP used to make a new Client and start a http proxy to listen
func RunSSClientAsHTTP(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	c, err := NewSSClient(address, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return err
	}
	return c.ListenAndServeHTTP(nil)
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
func RunSocks5Server(address, username, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) error {
	s, err := NewSocks5Server(address, username, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
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
	return s.ListenAndServe(nil)
}
