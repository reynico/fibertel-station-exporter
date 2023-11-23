package collector

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var csrfToken string

type FibertelStation struct {
	URL      string
	Username string
	Password string
	client   *http.Client
}

type LoginResponseSalts struct {
	Error     string `json:"error"`
	Salt      string `json:"salt"`
	SaltWebUI string `json:"saltwebui"`
}

type LoginResponse struct {
	Error   string             `json:"error"`
	Message string             `json:"message"`
	Data    *LoginResponseData `json:"data"`
}

type LoginResponseData struct {
	Interface       string `json:"intf"`
	User            string `json:"user"`
	Uid             string `json:"uid"`
	DefaultPassword string `json:"Dpd"`
	RemoteAddress   string `json:"remoteAddr"`
	UserAgent       string `json:"userAgent"`
	HttpReferer     string `json:"httpReferer"`
}

type LogoutResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ModemStatusResponse struct {
	Error   string           `json:"error"`
	Message string           `json:"message"`
	Data    *ModemStatusData `json:"data"`
}

type ModemStatusData struct {
	OfdmUpstreamData   []*OfdmUpstreamData        `json:"exUSTbl"`
	OfdmDownstreamData []*OfdmDownstreamData      `json:"exDSTbl"`
	Downstream         []*DocsisDownstreamChannel `json:"DSTbl"`
	Upstream           []*DocsisUpstreamChannel   `json:"USTbl"`
}

type OfdmDownstreamData struct {
	Id                   string `json:"__id"`
	ChannelIdOfdm        string `json:"ChannelID"`
	StartFrequency       string `json:"StartFrequency"`
	PLCFrequency         string `json:"PLCFrequency"`
	CentralFrequencyOfdm string `json:"CentralFrequency"`
	Bandwidth            string `json:"BandWidth"`
	PowerOfdm            string `json:"PowerLevel"`
	SnrOfdm              string `json:"SNRLevel"`
	FftOfdm              string `json:"FFT"`
	LockedOfdm           string `json:"LockStatus"`
	ChannelType          string `json:"ChannelType"`
}

type OfdmUpstreamData struct {
	Id                   string `json:"__id"`
	ChannelIdOfdm        string `json:"ChannelID"`
	StartFrequency       string `json:"StartFrequency"`
	PLCFrequency         string `json:"PLCFrequency"`
	CentralFrequencyOfdm string `json:"CentralFrequency"`
	Bandwidth            string `json:"BandWidth"`
	PowerOfdm            string `json:"PowerLevel"`
	FftOfdm              string `json:"FFT"`
	LockedOfdm           string `json:"LockStatus"`
	ChannelType          string `json:"ChannelType"`
}

type DocsisDownstreamChannel struct {
	Id               string `json:"__id"`
	ChannelId        string `json:"ChannelID"`
	CentralFrequency string `json:"Frequency"`
	Power            string `json:"PowerLevel"`
	Snr              string `json:"SNRLevel"`
	Modulation       string `json:"Modulation"`
	Locked           string `json:"LockStatus"`
	ChannelType      string `json:"ChannelType"`
}

type DocsisUpstreamChannel struct {
	Id               string `json:"__id"`
	ChannelIdUp      string `json:"ChannelID"`
	CentralFrequency string `json:"Frequency"`
	Power            string `json:"PowerLevel"`
	ChannelType      string `json:"ChannelType"`
	SymbolRate       string `json:"SymbolRate"`
	LockStatus       string `json:"LockStatus"`
}

func NewFibertelStation(stationUrl, username, password string) *FibertelStation {
	cookieJar, err := cookiejar.New(nil)
	parsedUrl, err := url.Parse(stationUrl)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cookieJar.SetCookies(parsedUrl, []*http.Cookie{
		{
			Name:  "Cwd",
			Value: "No",
		},
	})
	if err != nil {
		panic(err)
	}
	return &FibertelStation{
		URL:      stationUrl,
		Password: password,
		Username: username,
		client: &http.Client{
			Jar:       cookieJar,
			Timeout:   time.Second * 20, // getting DOCSIS status can be slow!
			Transport: tr,
		},
	}
}

func (v *FibertelStation) Login() (*LoginResponse, error) {
	_, err := v.doRequest("GET", v.URL, "")
	if err != nil {
		return nil, err
	}
	loginResponseSalts, err := v.getLoginSalts()
	if err != nil {
		return nil, err
	}

	derivedPassword := GetLoginPassword(v.Password, loginResponseSalts.Salt, loginResponseSalts.SaltWebUI)
	data := url.Values{}
	data.Set("username", v.Username)
	data.Set("password", derivedPassword)

	responseBody, err := v.doRequest("POST", v.URL+"/api/v1/session/login", data.Encode())
	if err != nil {
		return nil, err
	}
	loginResponse := &LoginResponse{}
	err = json.Unmarshal(responseBody, loginResponse)
	if loginResponse.Error != "ok" {
		return nil, fmt.Errorf("got non error=ok message from fibertel station")
	}

	// This is a dummy request, somehow this is required in order to make the posterior GETs
	responseMenu, err := v.doRequest("GET", v.URL+"/api/v1/session/menu", "")
	if err != nil {
		return nil, err
	}

	log.Debugf("Response menu: %s\n", responseMenu)

	return loginResponse, nil
}

func (v *FibertelStation) Logout() (*LogoutResponse, error) {
	responseBody, err := v.doRequest("POST", v.URL+"/api/v1/session/logout", "")
	if err != nil {
		return nil, err
	}
	logoutResponse := &LogoutResponse{}
	err = json.Unmarshal(responseBody, logoutResponse)
	if err != nil {
		return nil, err
	}
	if logoutResponse.Error != "ok" {
		return nil, fmt.Errorf("Got non error=ok message from fibertel station")
	}
	return logoutResponse, nil
}

func (v *FibertelStation) GetModemStatus() (*ModemStatusResponse, error) {
	responseBody, err := v.doRequest("GET", v.URL+"/api/v1/modem/exUSTbl,exDSTbl,USTbl,DSTbl?_="+strconv.FormatInt(makeTimestamp(), 10), "")
	if err != nil {
		return nil, err
	}
	log.Debugf("Docsis response body: %s\n", responseBody)
	modemStatusResponse := &ModemStatusResponse{}
	return modemStatusResponse, json.Unmarshal(responseBody, modemStatusResponse)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (v *FibertelStation) getLoginSalts() (*LoginResponseSalts, error) {
	data := url.Values{}
	data.Set("username", v.Username)
	data.Set("password", "seeksalthash")
	data.Set("logout", "true")
	responseBody, err := v.doRequest("POST", v.URL+"/api/v1/session/login", data.Encode())
	if err != nil {
		return nil, err
	}
	loginResponseSalts := &LoginResponseSalts{}
	err = json.Unmarshal(responseBody, loginResponseSalts)
	if err != nil {
		return nil, err
	}
	if loginResponseSalts.Error != "ok" {
		return nil, fmt.Errorf("Got non error=ok message from fibertel station")
	}
	return loginResponseSalts, nil
}

func (v *FibertelStation) doRequest(method, url, body string) ([]byte, error) {
	requestBody := strings.NewReader(body)
	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		log.Errorf("error building request: %s", err.Error())
		return nil, err
	}
	if method == "POST" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	request.Header.Set("Referer", "http://192.168.0.1")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("X-Csrf-Token", csrfToken)
	response, err := v.client.Do(request)
	if err != nil {
		log.Errorf("error performing request: %s", err.Error())
		return nil, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	csrfTokenRegex := regexp.MustCompile(`auth=([^;]+)`)
	csrfTokenMatches := csrfTokenRegex.FindStringSubmatch(response.Header.Get("Set-Cookie"))
	if len(csrfTokenMatches) >= 2 {
		csrfToken = csrfTokenMatches[1]
	}

	return io.ReadAll(response.Body)
}

// GetLoginPassword derives the password using the given salts
func GetLoginPassword(password, salt, saltWebUI string) string {
	return DoPbkdf2NotCoded(DoPbkdf2NotCoded(password, salt), saltWebUI)
}

func DoPbkdf2NotCoded(key, salt string) string {
	temp := pbkdf2.Key([]byte(key), []byte(salt), 0x3e8, 0x80, sha256.New)
	return hex.EncodeToString(temp[:16])
}
