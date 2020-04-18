# Usage

```
$> cfgenerator -interpreter=gotpl volumes/secrets volumes/config < config.conf.tpl
upstream api { server host.docker.internal:1337; }

server {
  location / {
    proxy_pass http://api;
  }
}
```
