### Setting Development Environment and Running

1. Create a soft link to link the GOPATH to the actual files: `ln -s /<GOPATH>/src/github.com/visenze/goproxy /path/to/grabber/goproxy`.
2. For the imports from third parties, run `go get ...`.
3. Under `goproxy` directory, run `go build`.
4. To start the service, run `./proxy http -p <local ip> -T tcp --log proxy.log`, `-p` is for the local ip, `-T` is for the parent proxy type, `--log` specifies the log file.

### Design Pattern

1. [The simple version](../../docs/images/proxy_design_pattern(simple).png)
2. [Work flow](../../docs/images/Proxy_Flow.jpg)

### Workflow

#### Before Requests

Before any http/https requests, initialising the service will first create the local connection with the local ip being the params input.

#### For Each Request

1. The service extracts the request's url
2. The service checks the pattern cache and pattern table to look for matched patterns for the url, and return the mapped `proxyName`. The pattern cache is cleared for every 5 mins.
3. If there is no `proxyName` mapped, it means `external proxy provider` is not required, and the request will be made directly.
4. With the `proxyName`, the service looks for proxy cache and proxy table for matched `Proxy`, proxy information like `endpoint`, `port`, `proxyType`, `apiEndpoint` are included in the `Proxy` structure.
5. If the `proxyType` is "static", the service will dial the ip of the proxy, and returns the response.
6. If the `proxyType` is "dynamic", the service will first call external api as provided by `apiEndpoint` and get 10 temporary proxy ips.
7. The service connects to one of the ips randomly.
8. The proxy cache is cleared for every 5 mins.

### Detailed documentation

[link to detailed document](https://docs.google.com/document/d/1w6Jqn3qfW6YWYabCAYBqQz9_Pt3ErWacpnpxcO_ilsQ/edit?pli=1#)
[link to detailed setting for original proxy](https://github.com/snail007/goproxy)

### Remarks

Edited files are in `utils` and `http` folders
