# Decompress Request
Plugin to decompress compressed requests in gzip

Headers:

- gzip Flag: ```x-sensedia-gzip: true || false```
- Inform Content-Type: ```x-sensedia-content-type```

## Configuration
The following declaration (given here in YAML) defines a plugin:
```YAML
experimental:
  plugins:
    traefik-plugin-decompress:
      moduleName: "github.com/Sensedia/traefik-plugin-decompress"
      version: "v1.0.5"
```
Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the http.middlewares section:
```YAML
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
    name: my-traefik-plugin-decompress
    namespace: my-namespace
spec:
    plugin:
        traefik-plugin-decompress:
            responseHeader: "200"

```