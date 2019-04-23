package cloud_function

import (
	"fmt"
	"strings"
)

type Resource string

func (r Resource) InstanceID() (string, error) {
	fragments := strings.Split(string(r), "/")
	if len(fragments) < 4 {
		return "", fmt.Errorf("invalid resource %v", r)
	}
	return fragments[3], nil

}

func (r Resource) Key() (string, error) {
	fragments := strings.Split(string(r), "/")
	if len(fragments) < 4 {
		return "", fmt.Errorf("invalid resource %v", r)
	}
	return fragments[len(fragments)-2], nil
}

func (r Resource) DatabaseURL() (string, error) {
	instanceID, err := r.InstanceID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s.firebaseio.com", instanceID), nil
}

func (r Resource) RefPath() (string, error) {
	key, err := r.Key()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("posts/%v", key), nil
}
