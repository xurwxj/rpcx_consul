package registry

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/smallnest/rpcx/log"
	"github.com/soheilhy/cmux"
	"github.com/xurwxj/gtils/base"
)

type CmuxServicePlugin struct {
	Log               log.Logger
	URL               string
	RefreshServiceFun func(services map[string][]string)
}
type ServicesList struct {
	Services map[string][]string
}

func (s *CmuxServicePlugin) MuxMatch(m cmux.CMux) {
	http1Matcher := cmux.HTTP1HeaderFieldPrefix("Consul-Service", "services")
	http1aMatcher := cmux.HTTP1HeaderFieldPrefix("Consul-Service", "[\"services")
	http2Matcher := cmux.HTTP2HeaderFieldPrefix("Consul-Service", "services")
	http2aMatcher := cmux.HTTP2HeaderFieldPrefix("Consul-Service", "[\"services")

	listener := m.Match(http1Matcher, http1aMatcher, http2Matcher, http2aMatcher)
	mux := http.NewServeMux()
	mux.HandleFunc(s.URL, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.Log.Errorf("consulServiceUpdate fail err =%v", err)
			return
		}
		services := make(map[string][]string)
		if err := json.Unmarshal(b, &services); err != nil {
			s.Log.Errorf("consulServiceUpdate fail err =%v", err)
			return

		}

		services = dealServiceData(services)
		s.Log.Debugf("consulServiceUpdate services=%v services=%v ", string(b), services)
		s.RefreshServiceFun(services)
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
	})

	httpS := &http.Server{
		Handler: mux,
	}
	go httpS.Serve(listener)
}

func dealServiceData(serviceMap map[string][]string) map[string][]string {
	newServiceMap := make(map[string][]string)
	for key, nodes := range serviceMap {
		for _, node := range nodes {
			if services, ok := newServiceMap[node]; ok {
				services = append(services, key)
				newServiceMap[node] = services
				continue
			}
			newServiceMap[node] = []string{key}

		}
	}
	return newServiceMap
}
