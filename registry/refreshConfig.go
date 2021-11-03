package registry

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/smallnest/rpcx/log"
	"github.com/soheilhy/cmux"
	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/viper"
)

type CmuxPluginConfig struct {
	CmuxPlugin
	Url string
}

type ConsulKVRep struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (s *CmuxPluginConfig) MuxMatch(m cmux.CMux) {
	http1Matcher := cmux.HTTP1HeaderFieldPrefix("Consul-Update", "param")
	http1aMatcher := cmux.HTTP1HeaderFieldPrefix("Consul-Update", "[\"param")
	http2Matcher := cmux.HTTP2HeaderFieldPrefix("Consul-Update", "param")
	http2aMatcher := cmux.HTTP2HeaderFieldPrefix("Consul-Update", "[\"param")

	listener := m.Match(http1Matcher, http1aMatcher, http2Matcher, http2aMatcher)
	mux := http.NewServeMux()
	mux.HandleFunc(s.Url, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("consulConfigUpdate fail err =%v", err)
			return
		}
		c := make([]*ConsulKVRep, 0)
		json.Unmarshal(b, &c)
		if len(c) > 0 {
			v := c[0].Value
			newValue, _ := base64.StdEncoding.DecodeString(v)
			MergeConfig(string(newValue))
			log.Debugf("consulConfigUpdate kvKey=%v  value=%v ", c[0].Key, string(newValue))
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

func MergeConfig(data string) {
	if data == "" {
		return
	}
	base.CheckPathExistOrCreate("nacos")
	err := os.WriteFile("nacos/remote.json", []byte(data), os.ModePerm)
	if err != nil {
		log.Errorf("mergeConfig:WriteFile config=%v err=%v ", data, err)
	}
	err = viper.MergeConfig(bytes.NewReader([]byte(data)))
	if err != nil {
		log.Errorf("mergeConfig config=%v  err=%v ", data, err)
	} else {
		log.Infof("mergeConfig done success ")
	}
}
