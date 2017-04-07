package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/IBM/ubiquity/resources"
	"github.com/gorilla/mux"
)

func ExtractErrorResponse(response *http.Response) error {
	errorResponse := resources.GenericResponse{}
	err := UnmarshalResponse(response, &errorResponse)
	if err != nil {
		return err
	}
	return fmt.Errorf("%s", errorResponse.Err)
}

func FormatURL(url string, entries ...string) string {
	base := url
	if !strings.HasSuffix(url, "/") {
		base = fmt.Sprintf("%s/", url)
	}
	suffix := ""
	for _, entry := range entries {
		suffix = path.Join(suffix, entry)
	}
	return fmt.Sprintf("%s%s", base, suffix)
}

func HttpExecute(httpClient *http.Client, logger *log.Logger, requestType string, requestURL string, rawPayload interface{}) (*http.Response, error) {
	payload, err := json.MarshalIndent(rawPayload, "", " ")
	if err != nil {
		logger.Printf("Internal error marshalling params %#v", err)
		return nil, fmt.Errorf("Internal error marshalling params")
	}

	request, err := http.NewRequest(requestType, requestURL, bytes.NewBuffer(payload))
	if err != nil {
		logger.Printf("Error in creating request %#v", err)
		return nil, fmt.Errorf("Error in creating request")
	}

	return httpClient.Do(request)
}

func WriteResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	fmt.Fprintf(w, string(data))
}

func Unmarshal(r *http.Request, object interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}

func UnmarshalResponse(r *http.Response, object interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}
func UnmarshalDataFromRequest(r *http.Request, object interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, object)
	if err != nil {
		return err
	}

	return nil
}

func ExtractVarsFromRequest(r *http.Request, varName string) string {
	return mux.Vars(r)[varName]
}
