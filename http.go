package brook

// HTTPMiddleman is a middleman who can intercept and handle request
type HTTPMiddleman interface {
	// HandleHTTPProxy does not need to close conn,
	// if handled is true that means the request has been handled, whatever err
	func HandleHTTPProxy(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error)
}
