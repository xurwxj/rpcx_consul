package registry

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/smallnest/rpcx/log"
	"github.com/soheilhy/cmux"
	"github.com/xurwxj/gtils/base"
)

type CmuxConfigPlugin struct {
	Log             log.Logger
	URL             string
	ConfigPath      string
	MergeConfigFunc func(data string)
}

type ConsulKVRep struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (s *CmuxConfigPlugin) MuxMatch(m cmux.CMux) {
	http1Matcher := cmux.HTTP1HeaderFieldPrefix("Consul-Update", "param")
	http1aMatcher := cmux.HTTP1HeaderFieldPrefix("Consul-Update", "[\"param")
	http2Matcher := cmux.HTTP2HeaderFieldPrefix("Consul-Update", "param")
	http2aMatcher := cmux.HTTP2HeaderFieldPrefix("Consul-Update", "[\"param")

	listener := m.Match(http1Matcher, http1aMatcher, http2Matcher, http2aMatcher)
	mux := http.NewServeMux()
	mux.HandleFunc(s.URL, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.Log.Errorf("consulConfigUpdate fail err =%v", err)
			return
		}
		c := make([]*ConsulKVRep, 0)
		json.Unmarshal(b, &c)
		if len(c) > 0 {
			v := c[0].Value
			newValue, _ := base64.StdEncoding.DecodeString(v)
			s.MergeConfigFunc(string(newValue))
			s.Log.Debugf("consulConfigUpdate kvKey=%v  value=%v ", c[0].Key, string(newValue))
		}
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
