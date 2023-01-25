# SGProxy - Simple GO Proxy

Basic HTTP/S proxy. Created to add HTTP Auth to a request from a client that doesn't support supplying
auth in URL, for example VScode's Juypyter Notebook Server browser.

# Usage
```
Usage of sgproxy:
  -bind string
        Address to listen on, default '127.0.0.1' (default "127.0.0.1")
  -port int
        Port to listen on, default 8443 (default 8000)
  -url string
        Default URL to redirect to if 'sgproxy-url' query param not present (default "https://www.google.com")
  -username string
        Default HTTP Auth username to redirect to if 'sgproxy-username' query param not present
  -password string
        Default HTTP Auth password to redirect to if 'sgproxy-password' query param not present
  -tls-cert string
        Listen using HTTPS, path to TLS certificate to use
  -tls-key string
        Listen using HTTPS, path to TLS key to use
```

# Examples
```bash
# Proxy to URL, adding HTTP auth
sgproxy -username 'admin' -password 'secretpassword' -url 'https://19b8-35-227-105-49.eu.ngrok.io/k/117266883'

# Proxy to URL, listen using HTTPS
#   First generate TLS certificate and key
openssl genrsa -out tls.key 4096
openssl req -new -x509 -days 1826 -key tls.key -out tls.crt -subj '/C=US/ST=Oregon/L=Portland/O=GoProxy/OU=GoProxy/CN=GoProxy'
#  Then run goproxy
sgproxy -port 8443 -tls-key tls.key -tls-cert tls.crt

# Send request to proxyt, but overwrite runtime-default url
# Used to proxy multiple URLs at once, or be different users
curl 'http://127.0.0.1:8000/?sgproxy-url=https://bing.com'
curl 'http://127.0.0.1:8000/?sgproxy-username=other&sgproxy-password=person'
```
