package config

// html/template changed to text/template
import (
	"bytes"
	"io/ioutil"
	"text/template"
)

// Don't think these defaults are ever used -- see models/models.go
var defaultConfig = Config{
	Dev:                 "tap0",
	Port:                1194,
	Proto:               "udp",
	DNSServerOne:        "8.8.8.8",
	DNSServerTwo:        "8.8.4.4",
	Server:              "server-bridge 192.168.1.250 255.255.255.0 192.168.1.2 192.168.1.5",
	Cipher:              "AES-256-CBC",
	Keysize:             256,
	Auth:                "SHA256",
	Dh:                  "dh2048.pem",
	Keepalive:           "10 120",
	IfconfigPoolPersist: "ipp.txt",
	CCEncryption:        "/etc/openvpn/easy-rsa/pki/ta.key",
	ExtraServerOptions:  "# client-config-dir /etc/openvpn/ccd",
	ExtraClientOptions:  "",
}

//Config model
type Config struct {
	Dev   string
	Port  int
	Proto string

	Ca   string
	Cert string
	Key  string

	Cipher  string
	Keysize int
	Auth    string
	Dh      string

	DNSServerOne        string
	DNSServerTwo        string
	Server              string
	IfconfigPoolPersist string
	Keepalive           string
	CCEncryption        string

	ExtraServerOptions string
	ExtraClientOptions string

	Management string

	PiVPNServer string
}

//New returns config object with default values
func New() Config {
	return defaultConfig
}

//GetText injects config values into template
func GetText(tpl string, c Config) (string, error) {
	t := template.New("config")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	t.Execute(buf, c)
	return buf.String(), nil
}

//SaveToFile reads teamplate and writes result to destination file
func SaveToFile(tplPath string, c Config, destPath string) error {
	template, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	str, err := GetText(string(template), c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destPath, []byte(str), 0644)
}
