upstream api { server host.docker.internal:{{ .API_PORT }}; }

server {
  location / {
    proxy_pass http://api;
  }
}
