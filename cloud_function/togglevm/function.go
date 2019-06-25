package togglevm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/api/compute/v1"
	"io"
	"net/http"
	"strings"
)

//ToggleVmFn HTTP cloud function entry point to start top box
func ToggleVmFn(writer http.ResponseWriter, httpRequest *http.Request) {
	request, err := newRequest(httpRequest.Body)

	if err == nil {
		err = handleEvent(request.Action, request, request.Selector, writer)
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func newRequest(reader io.Reader) (*InstanceRequest, error) {
	decoder := json.NewDecoder(reader)
	request := &InstanceRequest{}
	return request, decoder.Decode(&request)
}



func handleEvent(URI string, request *InstanceRequest, selector string, writer http.ResponseWriter) (error) {
	service, err := GetService()
	if err != nil {
		return  err
	}

	var response *InstanceResponse
	switch strings.ToLower(URI) {
	case "start":
		response = service.Start(&StartRequest{*request})

	case "stop":
		response = service.Stop(&StopRequest{*request})

	default:
		return fmt.Errorf("unsupported action: %v, valid: start,stop", URI)
	}
	if response == nil {
		return fmt.Errorf("response was nil")
	}



	if response.Status == StatusError  {
		return fmt.Errorf("%v", response.Error)
	}


	var netInterface *compute.NetworkInterface
	if len(response.Instance.NetworkInterfaces) > 0 {
		netInterface = response.Instance.NetworkInterfaces[0]
	}
	switch strings.ToLower(selector) {
			case "ip":
				if netInterface == nil {
					return fmt.Errorf("netInterface was nil")
				}
				_, err := writer.Write([]byte(netInterface.NetworkIP))
				return err
				
			case "natip":
				if netInterface == nil {
					return fmt.Errorf("netInterface was nil")
				}
			if len(netInterface.AccessConfigs) > 0 {
				_, err := writer.Write([]byte(netInterface.AccessConfigs[0].NatIP))
				return err
			}		
	}
	buffer := new(bytes.Buffer)
	if err = json.NewEncoder(buffer).Encode(response);err != nil {
		return err
	}
	_, err = writer.Write(buffer.Bytes())
	return err
}
