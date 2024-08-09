package parser

import (
	"github.com/go-kratos/kratos/v2/log"
	"sync"
)

type Parser struct {
	sync.RWMutex
	resource IDataResource
	configs  map[string]IDataConfig
}

func New() *Parser {
	return &Parser{
		RWMutex: sync.RWMutex{},
		configs: make(map[string]IDataConfig),
	}
}

func NewBuilder(res IDataResource) *Parser {
	return &Parser{
		configs:  make(map[string]IDataConfig),
		resource: res,
	}
}

func (p *Parser) Register(configs ...IDataConfig) {
	if len(configs) < 1 {
		return
	}

	for _, cfg := range configs {
		if cfg != nil {
			p.configs[cfg.Name()] = cfg
		}
	}
}

func (p *Parser) Init() {
	for _, cfg := range p.configs {
		data, err := p.resource.ReadBytes(cfg.Name())
		if err != nil {
			panic(err)
		}
		p.onchange(cfg, data)
		p.resource.Init(cfg)
	}

	p.resource.OnChange(func(configName string, data []byte) {
		cfg := p.GetIConfig(configName)
		if cfg != nil {
			p.onchange(cfg, data)
		}
	})
}

func (p *Parser) onchange(cfg IDataConfig, data []byte) {
	p.Lock()
	defer p.Unlock()

	err := cfg.OnReload(data)
	if err != nil {
		log.Errorf("[Config.Parser] cfg[%s] reload data error:%v", cfg.Name(), err)
		return
	}
}

func (p *Parser) GetIConfig(name string) IDataConfig {
	p.Lock()
	defer p.Unlock()

	return p.configs[name]
}
