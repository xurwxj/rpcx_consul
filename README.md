# rpcx_consul
consul plugin for rpcx

```
"watches":[
        {
            "type": "keyprefix",
            "prefix": "dev/eds",
            "handler_type": "http",
            "token":"eb2b7d55-3c37-1888-3dab-60fab07cbddc",
            "http_handler_config": {
              "path": "http://10.20.30.40:6974/keyprefix",
              "method": "POST",
              "header": { "Consul-Update": ["param"] },
              "timeout": "10s",
              "tls_skip_verify": false
            }
          }
    ]
```