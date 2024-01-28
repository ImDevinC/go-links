package config

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	ctx := context.Background()
	cfg, err := FromEnv(ctx)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	expected := Config{
		StaticPath: "/",
		StoreType:  StoreTypeMemory,
		SSO: SSOConfig{
			SamlCert:     []byte(defaultCert),
			SamlKey:      []byte(defaultKey),
			MetadataFile: "",
			EntityID:     "",
			CallbackURL:  "",
			Require:      false,
		},
	}
	diff := cmp.Diff(cfg, expected)
	if !assert.Equal(t, "", diff) {
		t.Fail()
	}
}

func TestSamlConfig(t *testing.T) {
	os.Setenv("SSO_SAML_CERT", "testCert")
	os.Setenv("SSO_SAML_KEY", "testKey")
	os.Setenv("SSO_METADATA_FILE", "testFile")
	os.Setenv("SSO_ENTITY_ID", "testEntity")
	os.Setenv("SSO_CALLBACK_URL", "testURL")
	os.Setenv("SSO_REQUIRE", "true")
	ctx := context.Background()
	cfg, err := FromEnv(ctx)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	expected := Config{
		StaticPath: "/",
		StoreType:  StoreTypeMemory,
		SSO: SSOConfig{
			SamlCert:     []byte("testCert"),
			SamlKey:      []byte("testKey"),
			MetadataFile: "testFile",
			EntityID:     "testEntity",
			CallbackURL:  "testURL",
			Require:      true,
		},
	}
	diff := cmp.Diff(cfg, expected)
	if !assert.Equal(t, "", diff) {
		t.Fail()
	}
}

func TestStoreTypes(t *testing.T) {
	cases := []struct {
		Name              string
		StoreTypeInput    string
		ExpectedStoreType StoreType
	}{
		{
			Name:              "Validate memoryType",
			StoreTypeInput:    "memory",
			ExpectedStoreType: StoreTypeMemory,
		},
		{
			Name:              "Validate fileType",
			StoreTypeInput:    "file",
			ExpectedStoreType: StoreTypeFile,
		},
		{
			Name:              "Validate fileType",
			StoreTypeInput:    "unknown",
			ExpectedStoreType: StoreTypeMemory,
		},
	}

	ctx := context.Background()
	for _, tc := range cases {
		os.Setenv("STORE_TYPE", tc.StoreTypeInput)
		cfg, err := FromEnv(ctx)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		expected := Config{
			StaticPath: "/",
			StoreType:  StoreType(tc.StoreTypeInput),
			SSO: SSOConfig{
				SamlCert:     []byte(defaultCert),
				SamlKey:      []byte(defaultKey),
				MetadataFile: "",
				EntityID:     "",
				CallbackURL:  "",
				Require:      false,
			},
		}
		diff := cmp.Diff(cfg, expected)
		if !assert.Equal(t, "", diff) {
			t.Fail()
		}
	}

}
