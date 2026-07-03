package balancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/nomen06/load-balancer/internal/proxy"
)

// using atomic here because it is much much faster for simple increment operatons rather than using mutexes
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}
func (b *Backend) Isalive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

type Serverpool struct {
	backends []*Backend
	current  uint64
}

func (s *Serverpool) AddBackend(targetURL string) error {
	purl, err := url.Parse(targetURL)
	if err != nil {
		return err
	}
	p, err := proxy.NewProxy(targetURL)
	if err != nil {
		return err
	}
	s.backends = append(s.backends, &Backend{
		URL:          purl,
		Alive:        true,
		ReverseProxy: p,
	})
	return nil
}

func (s *Serverpool) NextInd() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

func (s *Serverpool) Nextavl() *Backend {
	// if len(s.backends) == 0 {
	// 	return nil
	// }
	next := s.NextInd()
	l := len(s.backends)
	// return s.backends[next]

	// checking for the alive ones
	for i := 0; i < l; i++ {
		ind := (next + i) % l
		if s.backends[ind].Isalive() {
			if i != 0 {
				atomic.StoreUint64(&s.current, uint64(ind))
			}
			return s.backends[ind]
		}
	}
	return nil
}
func (s *Serverpool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//for testing
	// if r.URL.Path == "/panic" {
	// 	panic("PANICCCC")
	// }
	//....working(tested)
	peer := s.Nextavl()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "503", http.StatusServiceUnavailable)
}
