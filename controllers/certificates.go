package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/bnhf/go-openvpn/client/config"
	"github.com/bnhf/pivpn-tap-web-ui/lib"
	"github.com/bnhf/pivpn-tap-web-ui/models"
)

type NewCertParams struct {
	Name string `form:"Name" valid:"Required;"`
}

type CertificatesController struct {
	BaseController
}

func (c *CertificatesController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Certificates",
	}
}

// @router /certificates/single-config/:key [get]
func (c *CertificatesController) DownloadSingleConfig() {
	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.ovpn", name)

	c.Ctx.Output.Header("Content-Type", "text/plain")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	keysPath := models.GlobalCfg.OVConfigPath + "easy-rsa/pki"
	if cfgPath, err := saveClientSingleConfig(name, keysPath); err == nil {
		c.Ctx.Output.Download(cfgPath, filename)
	}

}

// @router /certificates/:key [get]
func (c *CertificatesController) Download() {
	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.zip", name)

	c.Ctx.Output.Header("Content-Type", "application/zip")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	zw := zip.NewWriter(c.Controller.Ctx.ResponseWriter)

	keysPath := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/"
	if cfgPath, err := saveClientConfig(name); err == nil {
		addFileToZip(zw, cfgPath)
	}
	addFileToZip(zw, keysPath+"ca.crt")
	addFileToZip(zw, keysPath+"issued/"+name+".crt")
	addFileToZip(zw, keysPath+"private/"+name+".key")
	addFileToZip(zw, keysPath+"ta.key")

	if err := zw.Close(); err != nil {
		beego.Error(err)
	}
}

func addFileToZip(zw *zip.Writer, path string) error {
	header := &zip.FileHeader{
		Name:         filepath.Base(path),
		Method:       zip.Store,
		ModifiedTime: uint16(time.Now().UnixNano()),
		ModifiedDate: uint16(time.Now().UnixNano()),
	}
	fi, err := os.Open(path)
	if err != nil {
		beego.Error(err)
		return err
	}

	fw, err := zw.CreateHeader(header)
	if err != nil {
		beego.Error(err)
		return err
	}

	if _, err = io.Copy(fw, fi); err != nil {
		beego.Error(err)
		return err
	}

	return fi.Close()
}

// @router /certificates [get]
func (c *CertificatesController) Get() {
	c.TplName = "certificates.html"
	c.showCerts()
}

func (c *CertificatesController) showCerts() {
	path := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/index.txt"
	certs, err := lib.ReadCerts(path)
	if err != nil {
		beego.Error(err)
	}
	lib.Dump(certs)
	c.Data["certificates"] = &certs
}

// @router /certificates [post]
func (c *CertificatesController) Post() {
	c.TplName = "certificates.html"
	flash := beego.NewFlash()

	cParams := NewCertParams{}
	if err := c.ParseForm(&cParams); err != nil {
		beego.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if vMap := validateCertParams(cParams); vMap != nil {
			c.Data["validation"] = vMap
		} else {
			if err := lib.CreateCertificate(cParams.Name); err != nil {
				beego.Error(err)
				flash.Error(err.Error())
				flash.Store(&c.Controller)
			}
		}
	}
	c.showCerts()
}

func validateCertParams(cert NewCertParams) map[string]map[string]string {
	valid := validation.Validation{}
	b, err := valid.Valid(&cert)
	if err != nil {
		beego.Error(err)
		return nil
	}
	if !b {
		return lib.CreateValidationMap(valid)
	}
	return nil
}

func saveClientConfig(name string) (string, error) {
	cfg := config.New()
	cfg.ServerAddress = models.GlobalCfg.ServerAddress
	cfg.Cert = name + ".crt"
	cfg.Key = name + ".key"
	cfg.Ta = "ta.key"
	serverConfig := models.OVConfig{Profile: "default"}
	serverConfig.Read("Profile")
	cfg.Dev = serverConfig.Dev
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher
	cfg.Keysize = serverConfig.Keysize
	cfg.ExtraClientOptions = serverConfig.ExtraClientOptions
	cfg.PiVPNServer = serverConfig.PiVPNServer

	destPath := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/" + name + ".conf"
	if err := config.SaveToFile("conf/openvpn-client-config.tpl",
		cfg, destPath); err != nil {
		beego.Error(err)
		return "", err
	}

	return destPath, nil
}

func saveClientSingleConfig(name string, pathString string) (string, error) {
	certPath := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/"
	cfg := config.New()
	cfg.ServerAddress = models.GlobalCfg.ServerAddress
	cfg.Cert = readCert(certPath + "issued/" + name + ".crt")
	cfg.Key = readCert(certPath + "private/" + name + ".key")
	cfg.Ca = readCert(certPath + "ca.crt")
	cfg.Ta = readCert(certPath + "ta.key")
	serverConfig := models.OVConfig{Profile: "default"}
	serverConfig.Read("Profile")
	cfg.ExtraClientOptions = serverConfig.ExtraClientOptions
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher
	cfg.Keysize = serverConfig.Keysize
	cfg.PiVPNServer = serverConfig.PiVPNServer

	destPath := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/" + name + ".ovpn"
	if err := config.SaveToFile("conf/openvpn-client-config.ovpn.tpl",
		cfg, destPath); err != nil {
		beego.Error(err)
		return "", err
	}

	return destPath, nil
}

func readCert(path string) string {
	buff, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		beego.Error(err)
		return ""
	}
	return string(buff)
}
