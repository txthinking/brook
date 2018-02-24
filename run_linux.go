package brook

// RunTproxy used to start a tproxy
func RunTproxy(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	c, err := NewTproxy(address, server, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return err
	}
	return c.ListenAndServe()
}
