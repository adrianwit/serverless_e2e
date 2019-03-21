package dstransfer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/viant/toolbox/cred"
	"github.com/viant/toolbox/secret"
)

var systemManager *ssm.SSM
func getSystemManager() (*ssm.SSM, error) {
	if systemManager != nil {
		return systemManager, nil
	}
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	systemManager = ssm.New(sess)
	return systemManager, nil
}


func getParameters(name string) (*ssm.Parameter, error) {
	systemManager, err := getSystemManager()
	if err != nil {
		return nil, err
	}
	output, err := systemManager.GetParameter(&ssm.GetParameterInput{
		Name:aws.String(name),
		WithDecryption:aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	return output.Parameter, nil
}


func getCredConfigFromParam(name string) (*cred.Config, error){
	param, err := getParameters(name)
	if err != nil {
		return nil, err
	}
	secretService := secret.New("", false)
	return secretService.GetCredentials(*param.Value)
}