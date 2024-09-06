package myconfig

var (
	GConfig Config
)

type Config struct {
	Project  string
	PrintLog bool
	Smtp     SmtpConfig
	Redis    RedisConfig
	Binance  CexConfig
}

type SmtpConfig struct {
	SmtpHost       string
	SmtpPort       string
	SenderEmail    string
	SenderPassword string
	Receivers      string
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

type CexConfig struct {
	ApiKey    string
	SecretKey string
}

func SetProject(project string) {
	GConfig.Project = project
}

func init() {
	//3985658674 ff6052024
	GConfig = Config{
		Project:  "go-report",
		PrintLog: true,
		Smtp:     SmtpConfig{SmtpHost: "smtp.qq.com", SmtpPort: "465", SenderEmail: "912858811@qq.com", SenderPassword: "qdbjvcktwrxabccg", Receivers: "mlks_2008@hotmail.com"},
		Redis:    RedisConfig{Host: "localhost:6379", Password: "", DB: 0},
		//ryy main
		//Binance:  CexConfig{ApiKey: "s7celdxF8wcfhcdtmJDrfAwfplhrWOhJmwAiLkwPCIKjyA8KvNjwD66gm1fJBmER", SecretKey: "xoEuc0N8Dl2WmkfHKMp0oxFvfRwcjjjfzFOhFVzjG79ohEsHhzTr52RrR1ENSMNr"},
		//wbd sub btcbusd_virtual@izkm2k6tnoemail.com
		Binance: CexConfig{ApiKey: "JUZWoFOfVM7abdnHYHgSevlGIfi6XyrZCWZc8YhKHDCor5g1An55CUCvrhS6aGx9", SecretKey: "t7SPjnakqH6F0bWe2dzh1TAREHMuzY0OS4sk7YV0Brfzyl1noExAV6Bnwd95F7si"},
	}
}
