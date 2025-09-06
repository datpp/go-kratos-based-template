package bootstrap

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	// etcd config
	etcdKratos "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	etcdV3 "go.etcd.io/etcd/client/v3"
	GRPC "google.golang.org/grpc"

	// consul config
	consulKratos "github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/hashicorp/consul/api"

	// nacos config
	nacosKratos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	nacosClients "github.com/nacos-group/nacos-sdk-go/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// getConfigKey Get valid config
func getConfigKey(configKey string, useBackslash bool) string {
	if useBackslash {
		return strings.Replace(configKey, `.`, `/`, -1)
	} else {
		return configKey
	}
}

// NewRemoteConfigSource New Remote Config Source
func NewRemoteConfigSource(configType, configHost, configKey string) config.Source {
	switch configType {
	case "nacos":
		uri, _ := url.Parse(configHost)
		h := strings.Split(uri.Host, ":")
		addr := h[0]
		port, _ := strconv.Atoi(h[1])
		return NewNacosConfigSource(addr, uint64(port), configKey)
	case "consul":
		return NewConsulConfigSource(configHost, configKey)
	case "etcd":
		return NewEtcdConfigSource(configHost, configKey)
	case "apollo":
		return NewApolloConfigSource(configHost, configKey)
	}
	return nil
}

// NewNacosConfigSource New Remote Config Source - Nacos
func NewNacosConfigSource(configAddr string, configPort uint64, configKey string) config.Source {
	sc := []nacosConstant.ServerConfig{
		*nacosConstant.NewServerConfig(configAddr, configPort),
	}

	cc := nacosConstant.ClientConfig{
		TimeoutMs:            10 * 1000, // http请求超时时间，单位毫秒
		BeatInterval:         5 * 1000,  // 心跳间隔时间，单位毫秒
		UpdateThreadNum:      20,        // 更新服务的线程数
		LogLevel:             "debug",
		CacheDir:             "../../configs/cache", // 缓存目录
		LogDir:               "../../configs/log",   // 日志目录
		NotLoadCacheAtStart:  true,                  // 在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: true,                  // 当服务列表为空时是否更新本地缓存，true--更新,false--不更新
	}

	nacosClient, err := nacosClients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	return nacosKratos.NewConfigSource(nacosClient,
		nacosKratos.WithGroup(getConfigKey(configKey, false)),
		nacosKratos.WithDataID("bootstrap.yaml"),
	)
}

// NewEtcdConfigSource New Remote Config Source - Etcd
func NewEtcdConfigSource(configHost, configKey string) config.Source {
	etcdClient, err := etcdV3.New(etcdV3.Config{
		Endpoints:   []string{configHost},
		DialTimeout: time.Second, DialOptions: []GRPC.DialOption{GRPC.WithBlock()},
	})
	if err != nil {
		panic(err)
	}

	etcdSource, err := etcdKratos.New(etcdClient, etcdKratos.WithPath(getConfigKey(configKey, true)))
	if err != nil {
		panic(err)
	}

	return etcdSource
}

// NewApolloConfigSource New Remote Config Source - Apollo
func NewApolloConfigSource(_, _ string) config.Source {
	return nil
}

// NewConsulConfigSource New Remote Config Source - Consul
func NewConsulConfigSource(configHost, configKey string) config.Source {
	consulClient, err := api.NewClient(&api.Config{
		Address: configHost,
	})
	if err != nil {
		panic(err)
	}

	consulSource, err := consulKratos.New(consulClient,
		consulKratos.WithPath(getConfigKey(configKey, true)),
	)
	if err != nil {
		panic(err)
	}

	return consulSource
}

// NewFileConfigSource New File Config Source
func NewFileConfigSource(filePath string) config.Source {
	return file.NewSource(filePath)
}

// NewConfigProvider New Config Provider
func NewConfigProvider(configType, configHost, configPath, configKey string) config.Config {
	if configType == "" || configHost == "" {
		return config.New(
			config.WithSource(
				NewFileConfigSource(configPath),
			),
		)
	}
	
	return config.New(
		config.WithSource(
			NewFileConfigSource(configPath),
			NewRemoteConfigSource(configType, configHost, configKey),
		),
	)
}
