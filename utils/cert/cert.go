package cert

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func CreateSignCertToFile(rootCa *x509.Certificate, rootKey *rsa.PrivateKey, domainOrIP string, expireDays int, name string) (err error) {
	cert, key, err := CreateSignCert(rootCa, rootKey, domainOrIP, expireDays)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(name+".crt", cert, 0755)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(name+".key", key, 0755)
	return
}
func CreateSignCert(rootCa *x509.Certificate, rootKey *rsa.PrivateKey, domainOrIP string, expireDays int) (certBytes []byte, keyBytes []byte, err error) {
	cer := &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63()), //证书序列号
		Subject: pkix.Name{
			Country:            []string{getCountry()},
			Organization:       []string{domainOrIP},
			OrganizationalUnit: []string{domainOrIP},
			CommonName:         domainOrIP,
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour * 24 * time.Duration(expireDays)),
		BasicConstraintsValid: true,
		IsCA:        false,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		//KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		EmailAddresses: []string{},
		IPAddresses:    []net.IP{},
	}
	if ip := net.ParseIP(domainOrIP); ip != nil {
		cer.IPAddresses = append(cer.IPAddresses, ip)
	} else {
		cer.DNSNames = append(cer.DNSNames, domainOrIP)
	}

	// cer.IPAddresses = append(cer.IPAddresses, alternateIPs...)
	// cer.DNSNames = append(cer.DNSNames, alternateDNS...)

	//生成公钥私钥对
	priKey, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		return
	}
	certBytes, err = x509.CreateCertificate(cryptorand.Reader, cer, rootCa, &priKey.PublicKey, rootKey)
	if err != nil {
		return
	}

	//编码证书文件和私钥文件
	caPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	certBytes = pem.EncodeToMemory(caPem)

	buf := x509.MarshalPKCS1PrivateKey(priKey)
	keyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: buf,
	}
	keyBytes = pem.EncodeToMemory(keyPem)
	return
}
func CreateCaToFile(name, domainOrIP string, expireDays int) (err error) {
	ca, key, err := CreateCa(domainOrIP, expireDays)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(name+".crt", ca, 0755)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(name+".key", key, 0755)
	return
}
func CreateCa(organization string, expireDays int) (certBytes []byte, keyBytes []byte, err error) {
	priv, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         organization,
			Organization:       []string{organization},
			OrganizationalUnit: []string{organization},
			Country:            []string{getCountry()},
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour * 24 * time.Duration(expireDays)),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}

	derBytes, err := x509.CreateCertificate(cryptorand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	// Generate cert
	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, err
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, nil, err
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}
func ParseCertAndKeyBytes(certPemFileByes, keyFileBytes []byte) (cert *x509.Certificate, privateKey *rsa.PrivateKey, err error) {
	//解析根证书
	cert, err = ParseCertBytes(certPemFileByes)
	if err != nil {
		return
	}
	//解析私钥
	privateKey, err = ParseKeyBytes(keyFileBytes)
	return
}
func ParseCertAndKey(certPemFile, keyFile string) (cert *x509.Certificate, privateKey *rsa.PrivateKey, err error) {
	//解析根证书
	cert, err = ParseCert(certPemFile)
	if err != nil {
		return
	}
	//解析私钥
	privateKey, err = ParseKey(keyFile)
	return
}
func ParseCert(certPemFile string) (cert *x509.Certificate, err error) {
	//解析证书
	certFile_, err := ioutil.ReadFile(certPemFile)
	if err != nil {
		return
	}
	cert, err = ParseCertBytes(certFile_)
	return
}
func ParseKey(keyFile string) (key *rsa.PrivateKey, err error) {
	//解析证书
	keyFile_, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return
	}
	key, err = ParseKeyBytes(keyFile_)
	return
}
func ParseCertBytes(certPemFileBytes []byte) (cert *x509.Certificate, err error) {
	caBlock, _ := pem.Decode(certPemFileBytes)
	cert, err = x509.ParseCertificate(caBlock.Bytes)
	return
}
func ParseKeyBytes(keyFileBytes []byte) (praKey *rsa.PrivateKey, err error) {
	keyBlock, _ := pem.Decode(keyFileBytes)
	praKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	return
}
func getCountry() string {
	CList := []string{"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO", "AR", "AT", "AU", "AZ", "BB", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BL", "BM", "BN", "BO", "BR", "BS", "BW", "BY", "BZ", "CA", "CF", "CG", "CH", "CK", "CL", "CM", "CN", "CO", "CR", "CS", "CU", "CY", "CZ", "DE", "DJ", "DK", "DO", "DZ", "EC", "EE", "EG", "ES", "ET", "FI", "FJ", "FR", "GA", "GB", "GD", "GE", "GF", "GH", "GI", "GM", "GN", "GR", "GT", "GU", "GY", "HK", "HN", "HT", "HU", "ID", "IE", "IL", "IN", "IQ", "IR", "IS", "IT", "JM", "JO", "JP", "KE", "KG", "KH", "KP", "KR", "KT", "KW", "KZ", "LA", "LB", "LC", "LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY", "MA", "MC", "MD", "MG", "ML", "MM", "MN", "MO", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NE", "NG", "NI", "NL", "NO", "NP", "NR", "NZ", "OM", "PA", "PE", "PF", "PG", "PH", "PK", "PL", "PR", "PT", "PY", "QA", "RO", "RU", "SA", "SB", "SC", "SD", "SE", "SG", "SI", "SK", "SL", "SM", "SN", "SO", "SR", "ST", "SV", "SY", "SZ", "TD", "TG", "TH", "TJ", "TM", "TN", "TO", "TR", "TT", "TW", "TZ", "UA", "UG", "US", "UY", "UZ", "VC", "VE", "VN", "YE", "YU", "ZA", "ZM", "ZR", "ZW"}
	return CList[int(randInt(4))%len(CList)]
}
func randInt(strLen int) int64 {
	codes := "123456789"
	codeLen := len(codes)
	data := make([]byte, strLen)
	rand.Seed(time.Now().UnixNano() + rand.Int63() + rand.Int63() + rand.Int63() + rand.Int63())
	for i := 0; i < strLen; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}
	i, _ := strconv.ParseInt(string(data), 10, 64)
	return i
}
