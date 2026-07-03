package balancer

import (
	"log"
	"net"
	"net/url"
	"time"
)

// simply checking the health on all backends by running a ping
// also updating their status
func (s *Serverpool) Healthcheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendalive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

// checking if the port is open
func isBackendalive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
