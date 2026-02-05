# GinWrapper

# Usage
To use GinWrapper you need to have GOLang installed on your system. You can download and install GOLang from the official website: [https://go.dev/](https://go.dev/).

Once you have GOLang installed, you can now use GinWrapper in your project:

```go
import (
	"github.com/IzomSoftware/GinWrapper/" 
)
```

then in your main.go you should define your HttpsServer:

```go
var (
	HttpsServer httpscore.HttpsServer
)
```

Now, you need to setup logger and configuration as well as responses:

```go
func main() {
	logger.SetupLogger("Your Project Name")

  // Adjust as needed
	configuration.DefaultConfig =
		configuration.Holder{
			Debug: false,
			HTTPSServer: configuration.HTTPSServer{
				Enabled:      true,
				Address:      "0.0.0.0",
				Port:         2009,
				APIUserAgent: "LiteGuard Client 1.0/b (Software)",
				TlsConfiguration: configuration.HttpsTlsConfiguration{
					Enable:   false,
					CertFile: "cert.pem",
					KeyFile:  "key.pem",
				},
			},
		}
  // setup configuration
	configuration.SetupConfig("config.toml")

  // add responses 
	httpscore.Responses["index"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/", "/index.html"},
	}
	httpscore.Responses["projects"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "projects.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/projects", "/projects.html"},
	}
	httpscore.Responses["projects"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "members.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/members", "/members.html"},
	}
	httpscore.Responses["technologies"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "technologies.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/technologies", "/technologies.html"},
	}
	httpscore.Responses["colleagues"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "colleagues.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/colleagues", "/colleagues.html"},
	}

  // first argument is templateDir and second one is assetsDir
	HttpsServer.ListenAndServe("assets/templates/*", "/assets")
}
```

## Contribution Guidelines 🤝

Feel free to contribute to the development of our project. we will notice it.
