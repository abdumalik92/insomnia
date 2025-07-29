package conf

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	AppName string `env:"APP_NAME"`
	AppENV  string `env:"APP_ENV"`

	DefaultUserDomain       string `env:"DEFAULT_USER_DOMAIN"`
	DefaultUserWorkstations string `env:"DEFAULT_USER_WORKSTATIONS"`
	ImmutableUserAttributes bool   `env:"IMMUTABLE_USER_ATTRIBUTES"`
	AllowedUserAttributes   string `env:"ALLOWED_USER_ATTRIBUTES"`

	KeycloakBaseUrl  string `env:"KEYCLOAK_BASE_URL"`
	KeycloakRealm    string `env:"KEYCLOAK_REALM"`
	KeycloakUserName string `env:"KEYCLOAK_USERNAME"`
	KeycloakPassword string `env:"KEYCLOAK_PASSWORD"`

	LDAPServer       string `env:"LDAP_SERVER"`
	LDAPBindUser     string `env:"LDAP_BIND_USER"`
	LDAPBindPassword string `env:"LDAP_BIND_PASSWORD"`
	LDAPBaseDN       string `env:"LDAP_BASE_DN"`

	SentryDSN   string `env:"SENTRY_DSN"`
	LogFilePath string `env:"LOG_FILE_PATH"`

	JaegerServiceName      string  `env:"JAEGER_SERVICE_NAME"`
	JaegerSamplerType      string  `env:"JAEGER_SAMPLER_TYPE"`
	JaegerSamplerParam     float64 `env:"JAEGER_SAMPLER_PARAM"`
	JaegerAgentHostPort    string  `env:"JAEGER_AGENT_HOST_PORT"`
	JaegerLogsEnabled      bool    `env:"JAEGER_LOGS_ENABLED"`
	JaegerReporterLogSpans bool    `env:"JAEGER_REPORTER_LOG_SPANS"`

	NewRabbitMQUser       string `env:"NEW_RABBITMQ_USER"`
	NewRabbitMQPassword   string `env:"NEW_RABBITMQ_PASSWORD"`
	NewRabbitMQHost       string `env:"NEW_RABBITMQ_HOST"`
	NewRabbitMQPort       string `env:"NEW_RABBITMQ_PORT"`
	NewRabbitMQUseTLS     bool   `env:"NEW_RABBITMQ_USE_TLS"`
	NewRabbitMQCACert     string `env:"NEW_RABBITMQ_CA_CERT"`
	NewRabbitMQClientCert string `env:"NEW_RABBITMQ_CLIENT_CERT"`
	NewRabbitMQClientKey  string `env:"NEW_RABBITMQ_CLIENT_KEY"`

	NewRabbitMQPrimaryVHost    string `env:"NEW_RABBITMQ_PRIMARY_VHOST"`
	NewRabbitMQPrimaryExchange string `env:"NEW_RABBITMQ_PRIMARY_EXCHANGE"`

	RabbitMQDeadLetterExchange string `env:"RABBITMQ_DEAD_LETTER_EXCHANGE"`

	KeycloakRabbitMQUser     string `env:"KEYCLOAK_RABBITMQ_USER"`
	KeycloakRabbitMQPassword string `env:"KEYCLOAK_RABBITMQ_PASSWORD"`
	KeycloakRabbitMQHost     string `env:"KEYCLOAK_RABBITMQ_HOST"`
	KeycloakRabbitMQPort     string `env:"KEYCLOAK_RABBITMQ_PORT"`

	KeycloakRabbitMQVHost           string `env:"KEYCLOAK_RABBITMQ_VHOST"`
	KeycloakRabbitMQPrimaryExchange string `env:"KEYCLOAK_RABBITMQ_PRIMARY_EXCHANGE"`

	TestRabbitMQUser       string `env:"TEST_RABBITMQ_USER"`
	TestRabbitMQPassword   string `env:"TEST_RABBITMQ_PASSWORD"`
	TestRabbitMQHost       string `env:"TEST_RABBITMQ_HOST"`
	TestRabbitMQPort       string `env:"TEST_RABBITMQ_PORT"`
	TestRabbitMQUseTLS     bool   `env:"TEST_RABBITMQ_USE_TLS"`
	TestRabbitMQCACert     string `env:"TEST_RABBITMQ_CA_CERT"`
	TestRabbitMQClientCert string `env:"TEST_RABBITMQ_CLIENT_CERT"`
	TestRabbitMQClientKey  string `env:"TEST_RABBITMQ_CLIENT_KEY"`

	TestRabbitMQVHost           string `env:"TEST_RABBITMQ_VHOST"`
	TestRabbitMQPrimaryExchange string `env:"TEST_RABBITMQ_PRIMARY_EXCHANGE"`
}

func NewConfig(v *viper.Viper) (*Config, error) {
	encoder := func(c *mapstructure.DecoderConfig) {
		c.TagName = "env"
	}
	c := new(Config)
	if err := v.Unmarshal(c, encoder); err != nil {
		return nil, err
	}
	return c, nil
}
