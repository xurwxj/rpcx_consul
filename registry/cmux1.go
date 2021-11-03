package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/soheilhy/cmux"
	"github.com/xurwxj/gtils/base"
)

type CmuxPlugin1 struct {
	CmuxPlugin
}

type ConsulKVRep struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// MuxMatch main CMuxMatch func
func (s *CmuxPlugin1) MuxMatch(m cmux.CMux) {
	http1Matcher := cmux.HTTP1HeaderFieldPrefix("Consul-Update", "param")
	http1aMatcher := cmux.HTTP1HeaderFieldPrefix("Consul-Update", "[\"param")
	http2Matcher := cmux.HTTP2HeaderFieldPrefix("Consul-Update", "param")
	http2aMatcher := cmux.HTTP2HeaderFieldPrefix("Consul-Update", "[\"param")

	listener := m.Match(http1Matcher, http1aMatcher, http2Matcher, http2aMatcher)
	mux := http.NewServeMux()
	mux.HandleFunc("/watch", func(w http.ResponseWriter, r *http.Request) {
		// a := r.URL.Query()
		b, err := ioutil.ReadAll(r.Body)
		c := make([]*ConsulKVRep, 0)
		json.Unmarshal(b, &c)
		if len(c) > 0 {
			v := c[0].Value
			newValue, _ := base64.StdEncoding.DecodeString(v)
			var entryInfos map[string]interface{}
			json.Unmarshal([]byte(newValue), &entryInfos)
			fmt.Println(",,,,", entryInfos, string(newValue))
		}

		// fmt.Println("....", err, err1, b, v, enc)
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

	mux.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		// a := r.URL.Query()
		b, err := ioutil.ReadAll(r.Body)
		fmt.Println("....service...", string(b), err)
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
