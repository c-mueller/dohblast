package doh

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"net/http"
)

var HTTPClient = http.DefaultClient

var (
	InvalidStatusErr = errors.New("Invalid HTTP Status code!")
)

func QueryDoH(endpoint string, question dns.Msg) (*dns.Msg, error) {
	d, err := question.Pack()
	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("%s?ct=application/dns-message&dns=%s", endpoint, base64.RawURLEncoding.EncodeToString(d))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/dns-message, application/dns-udpwireformat, application/json")
	req.Header.Set("User-Agent", "doh-blast 1.0.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, InvalidStatusErr
	}

	reply, err := ioutil.ReadAll(resp.Body)
	m := new(dns.Msg)
	err = m.Unpack(reply)
	if err != nil {
		return nil, err
	}

	return m, nil
}
