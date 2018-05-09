package configuration

import (
	"fmt"
	"steve/structs/configuration"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type ConfigurationImpl struct {
	consulClient *consulapi.Client
}

func insertConfig(kvPair *consulapi.KVPair, prefix string, configs map[string]string) error {
	if strings.HasPrefix(kvPair.Key, prefix) {
		k := kvPair.Key[len(prefix):]
		configs[k] = string(kvPair.Value)
		return nil
	}
	return fmt.Errorf("error consul response, Key:%s, Need Prefix:%s", kvPair.Key, prefix)
}

func (ci *ConfigurationImpl) GetConfig(param *configuration.ConfigGetParam) (configs map[string]string, err error) {
	configs = map[string]string{}
	err = nil

	kv := ci.consulClient.KV()
	version := param.Version
	if version == "" {
		version, err = ci.GetConfigVer(param.Env)
		if err != nil {
			return
		}
	}
	prefix := fmt.Sprintf("%s/%s/", param.Env, param.Version)
	key := prefix + param.Key
	if param.Prefix {
		kvPairs, _, err := kv.List(key, nil)
		if err != nil {
			return nil, err
		}
		for _, kvPair := range kvPairs {
			if err = insertConfig(kvPair, prefix, configs); err != nil {
				return nil, err
			}
		}
	} else {
		kvPair, _, err := kv.Get(key, nil)
		if err != nil {
			return nil, err
		}
		if err = insertConfig(kvPair, prefix, configs); err != nil {
			return nil, err
		}
	}
	return configs, nil
}

func (ci *ConfigurationImpl) GetConfigVer(env configuration.Env) (string, error) {
	key := string(env) + "/version"
	kv := ci.consulClient.KV()
	kvPair, _, err := kv.Get(key, nil)

	if err != nil {
		return "", err
	}
	if kvPair != nil {
		return "", nil
	}
	return string(kvPair.Value), nil
}

func NewConfiguration() (*ConfigurationImpl, error) {
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &ConfigurationImpl{
		consulClient: client,
	}, nil
}
