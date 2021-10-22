package registry

import (
	"net/http"

	"github.com/smallnest/rpcx/log"
	"github.com/soheilhy/cmux"
	"github.com/xurwxj/gtils/base"
)

// CmuxPlugin rpcx CMuxMatch plugin implements
type CmuxPlugin struct {
	Log      log.Logger
	Services *[]string
}

// MuxMatch main CMuxMatch func
func (s *CmuxPlugin) MuxMatch(m cmux.CMux) {
	http1Matcher := cmux.HTTP1HeaderFieldPrefix("Consul-Health-Check", "serviceCheck")
	http1aMatcher := cmux.HTTP1HeaderFieldPrefix("Consul-Health-Check", "[\"serviceCheck")
	http2Matcher := cmux.HTTP2HeaderFieldPrefix("Consul-Health-Check", "serviceCheck")
	http2aMatcher := cmux.HTTP2HeaderFieldPrefix("Consul-Health-Check", "[\"serviceCheck")
	listener := m.Match(http1Matcher, http1aMatcher, http2Matcher, http2aMatcher)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		serviceName := r.URL.Query().Get("service")
		// fmt.Println(serviceName)
		if base.FindInStringSlice(*s.Services, serviceName) {
			resByte, err := base.GetByteArrayFromInterface(map[string]interface{}{
				"status": "UP",
				"application": map[string]string{
					"status": "UP",
				},
			})
			if err == nil {
				_, err = w.Write(resByte)
				if err == nil {
					return
				}
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fail"))
	})
	httpS := &http.Server{
		Handler: mux,
	}
	go httpS.Serve(listener)
	err := m.Serve()
	if err != nil {
		s.Log.Errorf("MuxMatch Serve err: ", err)
	}
}
