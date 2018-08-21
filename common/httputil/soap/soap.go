package soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Body represents a SOAP message body
type Body struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Message string   `xml:",innerxml"`
}

// Envelope represents a SOAP message envelope
type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Xsi     string   `xml:"xmlns xsi,attr"`
	Xsd     string   `xml:"xmlns xsd,attr"`
	Body    Body
}

// DoRequest sends a SOAP request to the designated URL. msg is the payload to
// be put into the SOAP body.
func DoRequest(url string, msg interface{}) (string, error) {
	msgdata, err := xml.Marshal(msg)
	if err != nil {
		return "", err
	}

	env := &Envelope{
		Xsi:  "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:  "http://www.w3.org/2001/XMLSchema",
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: Body{Message: string(msgdata)},
	}
	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(env); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("SOAPAction", "\"\"")

	// bb, _ := httputil.DumpRequestOut(req, true)
	// fmt.Println(string(bb))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %v\n%v", resp.Status, string(result))
	}

	retval := Envelope{
		Body: Body{Message: ""},
	}
	if err := xml.Unmarshal(result, &retval); err != nil {
		return "", err
	}
	return retval.Body.Message, nil
}
