# go-links
go-links is a custom url-shortener service that allows users to easily manage their own links. There are lots of options, but I wanted to build my own to make sure it supported my specific use cases.

Once installed, you will want to your DNS Search Suffix to include the domain you installed to. IE: If you're FQDN is `go.mysite.com`, you would want to set your DNS Search Suffix to `mysite.com`. Now when you type `http://go` into your browser, you should be directed to your site.

## SSL
Unless you want to sign all of your own certificates with your own root CA which then has to be trusted on all your devices, your SSL cert will most likely have to use an FQDN (`go.mysite.com`). To help with this, the service redirects all traffic from `http://go` to the value of the `FQDN` environment variable. This will allow you to use a certificate covered endpoint for things like SAML.

## Setup
### Kubernetes
```shell
helm repo add golinks https://imdevinc.github.io/go-links
helm repo update
helm install golinks -n golinks --create-namespace --set config.fqdn=go.mysite.com
```

