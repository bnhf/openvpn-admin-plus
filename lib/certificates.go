package lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/bnhf/pivpn-tap-web-ui/models"
)

//Cert
//https://groups.google.com/d/msg/mailing.openssl.users/gMRbePiuwV0/wTASgPhuPzkJ
type Cert struct {
	EntryType   string
	Expiration  string
	ExpirationT time.Time
	Revocation  string
	RevocationT time.Time
	Serial      string
	FileName    string
	Details     *Details
}

type Details struct {
	Name         string
	CN           string
	Country      string
	Organisation string
	Email        string
}

func ReadCerts(path string) ([]*Cert, error) {
	certs := make([]*Cert, 0)
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return certs, err
	}
	lines := strings.Split(trim(string(text)), "\n")
	for _, line := range lines {
		fields := strings.Split(trim(line), "\t")
		if len(fields) != 6 {
			return certs,
				fmt.Errorf("incorrect number of lines in line: \n%s\n. Expected %d, found %d",
					line, 6, len(fields))
		}
		expT, _ := time.Parse("060102150405Z", fields[1])
		revT, _ := time.Parse("060102150405Z", fields[2])
		c := &Cert{
			EntryType:   fields[0],
			Expiration:  fields[1],
			ExpirationT: expT,
			Revocation:  fields[2],
			RevocationT: revT,
			Serial:      fields[3],
			FileName:    fields[4],
			Details:     parseDetails(fields[5]),
		}
		certs = append(certs, c)
	}

	return certs, nil
}

func parseDetails(d string) *Details {
	details := &Details{}
	lines := strings.Split(trim(string(d)), "/")
	for _, line := range lines {
		if strings.Contains(line, "") {
			fields := strings.Split(trim(line), "=")
			switch fields[0] {
			case "name":
				details.CN = fields[1]
			case "CN":
				details.Name = fields[1]
			case "C":
				details.Country = fields[1]
			case "O":
				details.Organisation = fields[1]
			case "emailAddress":
				details.Email = fields[1]
			default:
				beego.Warn(fmt.Sprintf("Undefined entry: %s", line))
			}
		}
	}
	return details
}

func trim(s string) string {
	return strings.Trim(strings.Trim(s, "\r\n"), "\n")
}

func CreateCertificate(name string, passphrase string) error {
	rsaPath := models.GlobalCfg.OVConfigPath + "easy-rsa"
	rsaIndex := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
	pass := false
	if passphrase != "" {
		pass = true
	}
	certs, err := ReadCerts(rsaIndex)
	if err != nil {
		//		beego.Debug(string(output))
		beego.Error(err)
		//		return err
	}
	Dump(certs)
	exists := false
	for _, v := range certs {
		if v.Details.Name == name {
			exists = true
			return err
		}
	}
	if !exists && !pass {
		cmd := exec.Command("/bin/bash", "-c",
			fmt.Sprintf(
				"%s/easyrsa --batch build-client-full %s nopass",
				rsaPath, name))
		cmd.Dir = models.GlobalCfg.OVConfigPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			beego.Debug(string(output))
			beego.Error(err)
			return err
		}
		return nil
	}
	if !exists && pass {
		cmd := exec.Command("/bin/bash", "-c",
			fmt.Sprintf(
				"%s/easyrsa --passout=pass:%s build-client-full %s",
				rsaPath, passphrase, name))
		cmd.Dir = models.GlobalCfg.OVConfigPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			beego.Debug(string(output))
			beego.Error(err)
			return err
		}
		return nil
	}
	return nil
}

func RevokeCertificate(name string, serial string) error {
	rsaPath := models.GlobalCfg.OVConfigPath + "easy-rsa"
	rsaIndex := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
	certs, err := ReadCerts(rsaIndex)
	if err != nil {
		beego.Error(err)
	}
	Dump(certs)
	for _, v := range certs {
		if v.Details.Name == name {
			cmd := exec.Command("/bin/bash", "-c",
				fmt.Sprintf(
					"%s/easyrsa --batch revoke %s &&"+
						"%s/easyrsa gen-crl &&"+
						"cp %s/pki/crl.pem %s/..",
					rsaPath, name, rsaPath, rsaPath, rsaPath))
			cmd.Dir = models.GlobalCfg.OVConfigPath
			output, err2 := cmd.CombinedOutput()
			if err2 != nil {
				beego.Debug(string(output))
				beego.Error(err2)
				return err2
			}
			return nil
		}
	}
	return nil
}

func RemoveCertificate(name string, serial string) error {
	rsaIndex := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
	certs, err := ReadCerts(rsaIndex)
	if err != nil {
		beego.Error(err)
	}
	Dump(certs)
	for _, v := range certs {
		if v.Details.Name == name {
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/certs_by_serial/" + serial + ".pem")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/issued/" + name + ".crt")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/private/" + name + ".key")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/reqs/" + name + ".req")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/" + name + ".ovpn")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/" + name + ".conf")
			lines, err := readLines(rsaIndex)
			if err != nil {
				beego.Error(err)
				return err
			}
			newrsaIndex := ""
			for _, line := range lines {
				if !checkSubstrings(line, name, "\t"+serial) {
					newrsaIndex += line + "\n"
				}
			}
			err = ioutil.WriteFile(rsaIndex, []byte(newrsaIndex), 0644)
			if err != nil {
				beego.Error(err)
				return err
			}
			return nil
		}
	}
	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func checkSubstrings(str string, subs ...string) bool {
	matches := 0
	isCompleteMatch := true
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			matches += 1
		} else {
			isCompleteMatch = false
		}
	}
	return isCompleteMatch
}
