package server

import (
	"context"

	dlserver "github.com/rancher/dynamiclistener/server"
	steve "github.com/rancher/steve/pkg/server"
	"github.com/rancher/wrangler/v2/pkg/k8scheck"
	"github.com/rancher/wrangler/v2/pkg/ratelimit"
	"k8s.io/client-go/rest"

	"github.com/llmos-ai/llmos-controller/pkg/controller"
	"github.com/llmos-ai/llmos-controller/pkg/server/config"
	"github.com/llmos-ai/llmos-controller/pkg/server/ui"
)

type APIServer struct {
	ctx             context.Context
	kubeconfig      string
	httpListenPort  int
	httpsListenPort int
	threadiness     int
	namespace       string
	name            string

	mgmt        *config.Management
	steveServer *steve.Server
	restConfig  *rest.Config
}

// Options define the api server options
type Options struct {
	Context         context.Context
	KubeConfig      string
	HTTPListenPort  int
	HTTPSListenPort int
	Threadiness     int
	Namespace       string
	Name            string
	Debug           bool
	ProfileAddress  string
}

func NewServer(o Options) (*APIServer, error) {
	s := &APIServer{
		ctx:             o.Context,
		kubeconfig:      o.KubeConfig,
		httpListenPort:  o.HTTPListenPort,
		httpsListenPort: o.HTTPSListenPort,
		threadiness:     o.Threadiness,
		namespace:       o.Namespace,
		name:            o.Name,
	}

	clientConfig, err := GetConfig(s.kubeconfig)
	if err != nil {
		return nil, err
	}

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	restConfig.RateLimiter = ratelimit.None
	s.restConfig = restConfig

	err = k8scheck.Wait(s.ctx, *restConfig)
	if err != nil {
		return nil, err
	}

	serverOptions, err := s.setDefaults(restConfig)
	if err != nil {
		return nil, err
	}

	// register the controller
	if err = controller.Register(s.ctx, s.mgmt, s.threadiness); err != nil {
		return nil, err
	}

	// set up a new api server
	s.steveServer, err = steve.New(o.Context, restConfig, serverOptions)
	if err != nil {
		return nil, err
	}

	// configure the api ui
	ui.ConfigureAPIUI(s.steveServer.APIServer)

	return s, nil
}

func (s *APIServer) setDefaults(cfg *rest.Config) (*steve.Options, error) {
	var err error
	opts := &steve.Options{}

	// set up the management config
	s.mgmt, err = config.SetupManagement(s.ctx, cfg, s.namespace)
	if err != nil {
		return nil, err
	}

	// define the next handler after the mgmt is setup
	r := NewRouter(s.mgmt)
	opts.Next = r.Routes()

	return opts, nil
}

func (s *APIServer) ListenAndServe(opts *dlserver.ListenOpts) error {
	return s.steveServer.ListenAndServe(s.ctx, s.httpsListenPort, s.httpListenPort, opts)
}
