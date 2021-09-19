package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/proxy"
)

const apiPrefix = "/"
const servicePort = 8080
const port = 8080
const namespace = "websocket-testcase-ns"
const service = "websocket-testcase"

func kubeConfigPath() string {
	envKubeConfigPath := os.Getenv("KUBECONFIG")
	if envKubeConfigPath != "" {
		return envKubeConfigPath
	}

	home := homedir.HomeDir()
	return filepath.Join(home, ".kube", "config")
}

func loadKubernetesConfiguration(kubeConfigPath string) clientcmd.ClientConfig {
	log.Printf("Using kube config %s", kubeConfigPath)
	configPathList := filepath.SplitList(kubeConfigPath)
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{}
	if len(configPathList) <= 1 {
		configLoadingRules.ExplicitPath = kubeConfigPath
	} else {
		configLoadingRules.Precedence = configPathList
	}
	contextName := ""
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		configLoadingRules,
		&clientcmd.ConfigOverrides{
			CurrentContext: contextName,
		},
	)
}

func main() {
	kubernetesConfig := loadKubernetesConfiguration(kubeConfigPath())
	restClientConfig, err := kubernetesConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	filter := &proxy.FilterServer{
		AcceptPaths:   proxy.MakeRegexpArrayOrDie(proxy.DefaultPathAcceptRE),
		RejectPaths:   proxy.MakeRegexpArrayOrDie(proxy.DefaultPathRejectRE),
		AcceptHosts:   proxy.MakeRegexpArrayOrDie(proxy.DefaultHostAcceptRE),
		RejectMethods: proxy.MakeRegexpArrayOrDie(proxy.DefaultMethodRejectRE),
	}

	proxyHandler, err := proxy.NewProxyHandler(apiPrefix, filter, restClientConfig, time.Second*2)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle(apiPrefix, proxyHandler)
	mux.Handle("/example/", getRerouteHttpHandlerExampleAPI(proxyHandler, namespace, service))

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", int(port)))
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Handler: mux,
	}

	log.Printf("Starting proxy. namespace: [%v], service name: [%s], port: [%v]", namespace, service, port)
	server.Serve(l)
}

func getExampleApiServerProxiedHostAndPath(namespace string, service string) string {
	return fmt.Sprintf("/api/v1/namespaces/%s/services/%s:%d/proxy/", namespace, service, servicePort)
}

func getRerouteHttpHandlerExampleAPI(proxyHandler http.Handler, namespace string, service string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, "/example/", getExampleApiServerProxiedHostAndPath(namespace, service), 1)
		proxyHandler.ServeHTTP(w, r)
	})
}
