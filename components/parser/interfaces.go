package parser

type (
	IDataConfig interface {
		Name() string
		OnReload(data []byte) error
	}

	IDataResource interface {
		Name() string
		Init(dataConfig IDataConfig)
		ReadBytes(configName string) (data []byte, error error)
		OnChange(fn ConfigChangeFn)
		Stop()
	}

	ConfigChangeFn func(configName string, data []byte)
)
