package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mozilla/tls-observatory/constants"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	pkg        = flag.String("pkg", "rules", "package name to be added to the output file")
	outputFile = flag.String("outputFile", "tls_config.go", "name of the output file")
)

// TLSConfURL url where Mozilla publishes the TLS ciphers recommendations
const TLSConfURL = "https://statics.tls.security.mozilla.org/server-side-tls-conf.json"

// ServerSideTLSJson contains all the available configurations and the version of the current document.
type ServerSideTLSJson struct {
	Configurations map[string]Configuration `json:"configurations"`
	Version        float64                  `json:"version"`
}

// Configuration represents configurations levels declared by the Mozilla server-side-tls
// see https://wiki.mozilla.org/Security/Server_Side_TLS
type Configuration struct {
	OpenSSLCiphersuites   []string `json:"openssl_ciphersuites"`
	OpenSSLCiphers        []string `json:"openssl_ciphers"`
	TLSVersions           []string `json:"tls_versions"`
	TLSCurves             []string `json:"tls_curves"`
	CertificateTypes      []string `json:"certificate_types"`
	CertificateCurves     []string `json:"certificate_curves"`
	CertificateSignatures []string `json:"certificate_signatures"`
	RsaKeySize            float64  `json:"rsa_key_size"`
	DHParamSize           float64  `json:"dh_param_size"`
	ECDHParamSize         float64  `json:"ecdh_param_size"`
	HstsMinAge            float64  `json:"hsts_min_age"`
	OldestClients         []string `json:"oldest_clients"`
	OCSPStaple            bool     `json:"ocsp_staple"`
	ServerPreferredOrder  bool     `json:"server_preferred_order"`
	MaxCertLifespan       float64  `json:"maximum_certificate_lifespan"`
}

type goCipherConfiguration struct {
	Name       string
	Ciphers    []string
	MinVersion string
	MaxVersion string
}

type goTLSConfiguration struct {
	cipherConfigs []goCipherConfiguration
}

// getTLSConfFromURL retrieves the json containing the TLS configurations from the specified URL.
func getTLSConfFromURL(url string) (*ServerSideTLSJson, error) {
	r, err := http.Get(url) //#nosec G107
	if err != nil {
		return nil, err
	}
	defer r.Body.Close() //#nosec G307

	var sstls ServerSideTLSJson
	err = json.NewDecoder(r.Body).Decode(&sstls)
	if err != nil {
		return nil, err
	}

	return &sstls, nil
}

func getGoCipherConfig(name string, sstls ServerSideTLSJson) (goCipherConfiguration, error) {
	caser := cases.Title(language.English)
	cipherConf := goCipherConfiguration{Name: caser.String(name)}
	conf, ok := sstls.Configurations[name]
	if !ok {
		return cipherConf, fmt.Errorf("TLS configuration '%s' not found", name)
	}

	// These ciphers are already defined in IANA format
	cipherConf.Ciphers = append(cipherConf.Ciphers, conf.OpenSSLCiphersuites...)

	for _, cipherName := range conf.OpenSSLCiphers {
		cipherSuite, ok := constants.CipherSuites[cipherName]
		if !ok {
			log.Printf("'%s' cipher is not available in crypto/tls package\n", cipherName)
		}
		if len(cipherSuite.IANAName) > 0 {
			cipherConf.Ciphers = append(cipherConf.Ciphers, cipherSuite.IANAName)
			if len(cipherSuite.NSSName) > 0 && cipherSuite.NSSName != cipherSuite.IANAName {
				cipherConf.Ciphers = append(cipherConf.Ciphers, cipherSuite.NSSName)
			}
		}
	}

	versions := mapTLSVersions(conf.TLSVersions)
	if len(versions) > 0 {
		cipherConf.MinVersion = fmt.Sprintf("0x%04x", versions[0])
		cipherConf.MaxVersion = fmt.Sprintf("0x%04x", versions[len(versions)-1])
	} else {
		return cipherConf, fmt.Errorf("no TLS versions found for configuration '%s'", name)
	}
	return cipherConf, nil
}

func getGoTLSConf() (goTLSConfiguration, error) {
	sstls, err := getTLSConfFromURL(TLSConfURL)
	if err != nil || sstls == nil {
		msg := fmt.Sprintf("Could not load the Server Side TLS configuration from Mozilla's website. Check the URL: %s. Error: %v\n",
			TLSConfURL, err)
		panic(msg)
	}

	tlsConfig := goTLSConfiguration{}

	modern, err := getGoCipherConfig("modern", *sstls)
	if err != nil {
		return tlsConfig, err
	}
	tlsConfig.cipherConfigs = append(tlsConfig.cipherConfigs, modern)

	intermediate, err := getGoCipherConfig("intermediate", *sstls)
	if err != nil {
		return tlsConfig, err
	}
	tlsConfig.cipherConfigs = append(tlsConfig.cipherConfigs, intermediate)

	old, err := getGoCipherConfig("old", *sstls)
	if err != nil {
		return tlsConfig, err
	}
	tlsConfig.cipherConfigs = append(tlsConfig.cipherConfigs, old)

	return tlsConfig, nil
}

func getCurrentDir() (string, error) {
	dir := "."
	if args := flag.Args(); len(args) == 1 {
		dir = args[0]
	} else if len(args) > 1 {
		return "", errors.New("only one directory at a time")
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func main() {
	dir, err := getCurrentDir()
	if err != nil {
		log.Fatalln(err)
	}
	tlsConfig, err := getGoTLSConf()
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	err = generatedHeaderTmpl.Execute(&buf, *pkg)
	if err != nil {
		log.Fatalf("Failed to generate the header: %v", err)
	}
	for _, cipherConfig := range tlsConfig.cipherConfigs {
		err := generatedRuleTmpl.Execute(&buf, cipherConfig)
		if err != nil {
			log.Fatalf("Failed to generated the cipher config: %v", err)
		}
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("warnings: Failed to format the code: %v", err)
		src = buf.Bytes()
	}

	outputPath := filepath.Join(dir, *outputFile)
	if err := os.WriteFile(outputPath, src, 0o644); err != nil /*#nosec G306*/ {
		log.Fatalf("Writing output: %s", err)
	}
}
