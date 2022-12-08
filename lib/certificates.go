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
	rsaPath := "/etc/openvpn/easy-rsa"
	pass := false
	if passphrase != "" {
		pass = true
	}
	//	//	varsPath := models.GlobalCfg.OVConfigPath + "easy-rsa/vars"
	//	cmd := exec.Command("/bin/bash", "-c",
	//		fmt.Sprintf(
	//			//			"source %s &&"+
	//			"export KEY_NAME=%s &&"+
	//				"%s/easyrsa --batch build-client-full %s nopass", name, rsaPath, name))
	//	cmd.Dir = models.GlobalCfg.OVConfigPath
	//	output, err := cmd.CombinedOutput()
	path := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
	certs, err := ReadCerts(path)
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
		}
	}
	if !exists && !pass {
		cmd := exec.Command("/bin/bash", "-c",
			fmt.Sprintf(
				//			    "source %s &&"+
				"export KEY_NAME=%s &&"+
					"%s/easyrsa --batch build-client-full %s nopass",
				name, rsaPath, name))
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
				//			    "source %s &&"+
				"export KEY_NAME=%s &&"+
					"export PASSPHRASE=%s &&"+
					"%s/easyrsa --passin=pass:$PASSPHRASE --pass out=pass:$PASSPHRASE build-client-full %s",
				name, passphrase, rsaPath, name))
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
	path := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
	certs, err := ReadCerts(path)
	if err != nil {
		beego.Error(err)
	}
	Dump(certs)
	for _, v := range certs {
		if v.Details.Name == name {
			rsaPath := "/etc/openvpn/easy-rsa/"
			//			varsPath := models.GlobalCfg.OVConfigPath + "keys/vars"

			cmd := exec.Command("/bin/bash", "-c",
				fmt.Sprintf(
					//					"source %s &&"+
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
	return nil //do nothing for now
}

func RemoveCertificate(name string, serial string) error {
	path := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
	certs, err := ReadCerts(path)
	if err != nil {
		beego.Error(err)
	}
	Dump(certs)
	for _, v := range certs {
		if v.Details.Name == name {
			keyDb := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
			/*file, err := os.Open(keyDb)
			    	if err != nil {
							beego.Error(err)
							return err
			    	}*/
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/certs_by_serial/" + serial + ".pem")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/issued/" + name + ".crt")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/private/" + name + ".key")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/" + name + ".ovpn")
			_ = os.Remove(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/" + name + ".conf")
			lines, err := readLines(keyDb)
			if err != nil {
				beego.Error(err)
				return err
			}
			newkeyDb := ""
			for _, line := range lines {
				if !checkSubstrings(line, name, "\t"+serial) {
					newkeyDb += line + "\n"
				}
			}
			err = ioutil.WriteFile(keyDb, []byte(newkeyDb), 0644)
			if err != nil {
				beego.Error(err)
				return err
			}
			return nil
		}
	}
	return nil //do nothing for now
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
