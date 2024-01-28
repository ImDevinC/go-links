package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

const defaultCert = `-----BEGIN CERTIFICATE-----
MIIDAzCCAeugAwIBAgIUH8uj5uCIpLBnBiUcAL7VdYXqvF8wDQYJKoZIhvcNAQEL
BQAwETEPMA0GA1UEAwwGZ29saW5rMB4XDTI0MDEyNjE1MzEzMVoXDTI1MDEyNTE1
MzEzMVowETEPMA0GA1UEAwwGZ29saW5rMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAwwveGXmb6Cm4GI2ht56Qxi9Cue9AwOZR9qx4ehsF6sRu55PxlTqL
RYUTGRIJ4OlyOsMIIeJq3ka7RGxJz3EUdrEOK+q1PgagUAQnhhUKBoDBr1u3UHmK
49U1D+U5hzck7jDlc0yU22z25v/SpGfBvLMhuS6f2bxFU6zupoiqR5e4yzgVlPeF
BMl0w39MMZtCLusEsWYZln/nwSkCvh/6odoHYteauCL3GHIznI55eVQ88x8SAe5m
dSkvJOAgp3dziivN2iqdyDc9535mMtYiQay9NSBdX7ORHrLz+9GvbuCS1o0LJ8X8
0eJVoIlozstBzg9v89Lohkm8TxV+7/XWXwIDAQABo1MwUTAdBgNVHQ4EFgQUgvtn
YfwPGO15zbgpWhm8yMsq/CUwHwYDVR0jBBgwFoAUgvtnYfwPGO15zbgpWhm8yMsq
/CUwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAIOBAsXj+fXOQ
NiLIcoVBJiKlXzPRTZD1ENTfm2UiSl3Rcq5EALXnphLR9TSHSVD9suA6s/JtSCj9
TIChJjHzqEBoaQKggJajblL3WYsx1kwlWB49ZbLIs4q5/GVMlBUyq+74hYAcPwxf
18g1ZUeVmicNXaEZ1yA4P4yIuhnt14UeV6+UdyywC90XulLnPcSZVekWibu+m+N2
yMy7P0WORk/ab5xBOmQekEhNRM9xtNmXr7S8KewQDYWAyVpeufgEvTflUfePEKbq
8+ZtjnREmYEvjBI2xtKCMb4VVw8IiowrQmWZXN+cPQYaTXXulx4IFrqFB1YIFVwf
6JcJJLzdig==
-----END CERTIFICATE-----`

const defaultKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDDC94ZeZvoKbgY
jaG3npDGL0K570DA5lH2rHh6GwXqxG7nk/GVOotFhRMZEgng6XI6wwgh4mreRrtE
bEnPcRR2sQ4r6rU+BqBQBCeGFQoGgMGvW7dQeYrj1TUP5TmHNyTuMOVzTJTbbPbm
/9KkZ8G8syG5Lp/ZvEVTrO6miKpHl7jLOBWU94UEyXTDf0wxm0Iu6wSxZhmWf+fB
KQK+H/qh2gdi15q4IvcYcjOcjnl5VDzzHxIB7mZ1KS8k4CCnd3OKK83aKp3INz3n
fmYy1iJBrL01IF1fs5EesvP70a9u4JLWjQsnxfzR4lWgiWjOy0HOD2/z0uiGSbxP
FX7v9dZfAgMBAAECggEAHMaydOW5N8535791nYaaa3LrkE0ZK5yPeSpG+BTmkZV7
m1T7bn3dsFsKz0cXCATJIpxFpeImzDZ5SIYFecKAN2a15YrSRJ1rp4KFZkXSXdU/
PiT07b2Q3T5Gftrd2vIq062JHLKuv5Ta9IfHxxO7xKBVGGIlmVUqkUbRSuac4MTi
fsBCRDV3HwpW2pLY3qUlsWQHOGxnenQ9G1wu7ArfY8XWuFFewxMu9OhscHiPpaBb
zrg7WbyfiUJrdIYmaLZAtKo9QIS66jQSH0Mw7DMwlZ1sVsHGGR62+twLHU/c531Y
NgxHxIYDl1JdgkDz/ngsZ00+VPcHsE4mLrlwIuMoSQKBgQD0VhHxrxTPoiO9Riwp
kss1NuEkPQ2psu4ARsQLTV4N+vG0Xe08UwB6ZNNssaNXbJWqNleNsWowRqvIwnTF
PHyhM2H9lEx6Dan9VDSW+7jz1XS0EDWFQvifzZHF6eIpYzTBvq7x/5uVTMM6sIl4
8szpPnfn99cVPHcDAV1yvdQuIwKBgQDMW3JrDdh0FMtsAfqxQxJWqjSX7HqKje0n
RTIZFgCAy/EgvIbXigSKoDP2UCoC4dN6YIjfNW6L99gCdtEN+mlFHcJ7V9h9FpWB
FuTgaiGHd8nyvYz2e8JGIcoMAxU/KLWKoRTh/D9DxJzs57oV8pmUizTItqiy31bo
TXAP5NbUlQKBgGHwtlSol7D7D1Rpcn5fpzD3hJvgFT/2x0w5EZBuPMth0c048UWD
B+gHzm/9bLo4fm2yRro3aZdcXLOmruP48QQ08oyRC27JV2CChmoXEPY8lAExliKK
y9pSrqIktFFewOEArGO40AaytHcsGI7w1I6SScIkKIUMra/4thquWQT9AoGBAIKC
TFouJ3RK63b49J9MVGPgo2H69m+SIEiaGlqHAJ9An6fmfr8cN7ZIhabin1Hj4uke
yYqzVvwwtlUsiGpC6APp85BOE8YfLC+a7WScovke+Wv6vhGUDAg6AA0X0vPZDceR
BAMm00h2QjnR67ekjYyeMoGUlbxWgewtuEmOPdzZAoGACzibWNEpZqYVPt2t1r25
i3LKPqDoEvCUYxFRvFmM3UzFF4lRVf8sijsuBkhiuGNwi0uHjsClQMYWtGRprV4j
Yba4Ys1ImIQypryr9CFZOP2ycr2PaajpwnaX01PxX0BnqeM3HDv6dWy4qAAEq8N/
cCkYBl4BcxLn1UdBhgCAY5w=
-----END PRIVATE KEY-----`

type Config struct {
	StaticPath string `env:"STATIC_PATH,default=/"`
	SSO        SSOConfig
	StoreType  StoreType `env:"STORE_TYPE,default=memory"`
	Mongo      MongoConfig
}

type SSOConfig struct {
	SamlCert     []byte `env:"SSO_SAML_CERT"`
	SamlKey      []byte `env:"SSO_SAML_KEY"`
	MetadataFile string `env:"SSO_METADATA_FILE"`
	EntityID     string `env:"SSO_ENTITY_ID"`
	CallbackURL  string `env:"SSO_CALLBACK_URL"`
	Require      bool   `env:"SSO_REQUIRE,default=false"`
}

type MongoConfig struct {
	Username     string `env:"MONGO_USERNAME"`
	Password     string `env:"MONGO_PASSWORD"`
	Host         string `env:"MONGO_HOST"`
	DatabaseName string `env:"MONGO_DB_NAME"`
}

type StoreType string

const (
	StoreTypeMemory StoreType = "memory"
	StoreTypeFile   StoreType = "file"
	StoreTypeMongo  StoreType = "mongo"
)

func FromEnv(ctx context.Context) (Config, error) {
	cfg := Config{}
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.SSO.SamlCert == nil {
		cfg.SSO.SamlCert = []byte(defaultCert)
	}
	if cfg.SSO.SamlKey == nil {
		cfg.SSO.SamlKey = []byte(defaultKey)
	}
	return cfg, nil
}
