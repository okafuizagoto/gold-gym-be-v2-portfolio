package config

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
	Topics  struct {
		LocalToProd string `yaml:"local_to_prod"`
		ProdToLocal string `yaml:"prod_to_local"`
		HTTPInsert  string `yaml:"http_insert"`
		Sales       string `yaml:"sales"`
	} `yaml:"topics"`
}

type (
	// Config ...
	Config struct {
		Server        ServerConfig        `yaml:"server"`
		Database      DatabaseConfig      `yaml:"database"`
		API           APIConfig           `yaml:"api"`
		Credential    Credential          `yaml:"credential"`
		Firebase      FirebaseConfig      `yaml:"firebase"`
		Swagger       SwaggerConfig       `yaml:"swagger"`
		Redis         Redis               `yaml:"redis"`
		Kafka         KafkaConfig         `yaml:"kafka"`
		Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
		MongoDB       MongoConfig         `yaml:"mongodb"`
		Discord       DiscordConfig       `yaml:"discord"`
	}

	// DiscordConfig menampung webhook untuk notifikasi error (lihat
	// internal/delivery/http/middleware/discord_alert.go). Webhook kosong =
	// notifikasi nonaktif, aman untuk environment development.
	DiscordConfig struct {
		Webhook string `yaml:"webhook"`
	}

	// MongoConfig holds MongoDB connection settings
	MongoConfig struct {
		URI      string `yaml:"uri"`
		Database string `yaml:"database"`
	}

	// ServerConfig ...
	ServerConfig struct {
		Port     string `yaml:"port"`
		GrpcPort string `yaml:"grpc_port"`
		Env      string `yaml:"env"`
	}

	// DatabaseConfig ...
	DatabaseConfig struct {
		Master     string `yaml:"master"`
		Production string `yaml:"production"`
	}

	// APIConfig ...
	APIConfig struct {
		Auth string `yaml:"auth"`
	}

	SwaggerConfig struct {
		Host    string   `yaml:"host"`
		Schemes []string `yaml:"schemes"`
	}

	Credential struct {
		Id string `yaml:"id"`
		Pw string `yaml:"pw"`
		Ip string `yaml:"ip"`
	}

	// FirebaseConfig ...
	FirebaseConfig struct {
		ProjectID     string `yaml:"projectID"`
		DatabaseURL   string `yaml:"databaseURL"`
		StorageBucket string `yaml:"storageBucket"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
	}

	// ElasticsearchConfig ...
	ElasticsearchConfig struct {
		Addresses []string `yaml:"addresses"`
		Username  string   `yaml:"username"`
		Password  string   `yaml:"password"`
	}
)
