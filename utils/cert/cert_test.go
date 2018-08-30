package cert

import (
	"os"
	"testing"
)

func TestCaGen(t *testing.T) {
	err := CreateCaToFile("ca", "test", 365)
	if err != nil {
		t.Fatal(err)
		return
	}
	ca, key, err := ParseCertAndKey("ca.crt", "ca.key")
	if err != nil {
		t.Fatal(err)
		return
	}
	if ca.Subject.Organization[0] != "test" {
		t.Fatalf("Organization %s not match test", ca.Subject.Organization[0])
		return
	}
	err = key.Validate()
	if err != nil {
		t.Fatal(err)
		return
	}
	os.Remove("ca.crt")
	os.Remove("ca.key")

}
func TestSign(t *testing.T) {
	err := CreateCaToFile("ca", "test", 365)
	if err != nil {
		t.Fatal(err)
		return
	}
	ca, key, err := ParseCertAndKey("ca.crt", "ca.key")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = CreateCaToFile("ca", "test", 365)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = CreateSignCertToFile(ca, key, "server.com", 365, "server")
	if err != nil {
		t.Fatal(err)
		return
	}
	servercrt, serverkey, err := ParseCertAndKey("server.crt", "server.key")
	if err != nil {
		t.Fatal(err)
		return
	}
	if servercrt.Subject.CommonName != "server.com" {
		t.Fatalf("CommonName %s not match server.com", ca.Subject.CommonName)
		return
	}
	err = serverkey.Validate()
	if err != nil {
		t.Fatal(err)
		return
	}
	os.Remove("ca.crt")
	os.Remove("ca.key")
	os.Remove("server.crt")
	os.Remove("server.key")
}
