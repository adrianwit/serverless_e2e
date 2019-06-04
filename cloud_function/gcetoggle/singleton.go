package gcetoggle


var srv Service

func GetService() (Service, error) {
	if srv != nil {
		return srv, nil
	}
	var err error
	srv, err = New()
	return srv, err
}

