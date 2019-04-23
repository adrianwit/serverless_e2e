package split

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/viant/toolbox/cred"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/option"
)

func getSecret(ctx context.Context, key string, cipherBase64Text string) (*cred.Config, error) {
	kmsService, err := cloudkms.NewService(ctx, option.WithScopes(cloudkms.CloudPlatformScope, cloudkms.CloudkmsScope))
	if err != nil {
		return nil, err
	}
	service := cloudkms.NewProjectsLocationsKeyRingsCryptoKeysService(kmsService)
	response, err := service.Decrypt(key, &cloudkms.DecryptRequest{Ciphertext: cipherBase64Text}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	plainConfig, err := base64.StdEncoding.DecodeString(string(response.Plaintext))
	if err != nil {
		return nil, err
	}
	credConfig := &cred.Config{}
	err = json.NewDecoder(bytes.NewReader(plainConfig)).Decode(credConfig)
	if err != nil {
		err = fmt.Errorf("failed to ecode: %v, due to:%v", plainConfig, err)
	}
	return credConfig, err
}
