package registry

import (
	"fmt"
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
			s.Log.Errorf("consulConfigUpdate fail err =%v", err)
			return
		}
		c := string(b)
		fmt.Println("..", c)
		// c := make([]*ConsulKVRep, 0)
		// json.Unmarshal(b, &c)
		// if len(c) > 0 {
		// 	v := c[0].Value
		// 	newValue, _ := base64.StdEncoding.DecodeString(v)
		// 	s.RefreshServiceFun(string(newValue))
		// 	s.Log.Debugf("consulConfigUpdate kvKey=%v  value=%v ", c[0].Key, string(newValue))
		// }
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
