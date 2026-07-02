package cert

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"

	"xpnsec.com/certificate-tool/v2/pkg/crypto"
)

// Types for ASN.1 SAN serialization.
type (
	// SubjectAltName is a struct that can be marshaled as ASN.1
	// into the SAN field in an x.509 certificate.
	//
	// See RFC 3280: https://www.ietf.org/rfc/rfc3280.txt
	//
	// T is the ASN.1 encodeable struct corresponding to an otherName
	// item of the GeneralNames sequence.
	SubjectAltName[T any] struct {
		OtherName otherName[T] `asn1:"tag:0"`
	}

	otherName[T any] struct {
		OID   asn1.ObjectIdentifier
		Value T `asn1:"tag:0"`
	}

	UPN struct {
		Value string `asn1:"utf8"`
	}

	ADSid struct {
		// Value is the bytes representation of the user's SID string,
		// e.g. []byte("S-1-5-21-1329593140-2634913955-1900852804-500")
		Value []byte // Gets encoded as an asn1 octet string
	}
)

var SubjectAltNameExtensionOID = asn1.ObjectIdentifier{2, 5, 29, 17}
var UPNOtherNameOID = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 20, 2, 3}
var EnhancedKeyUsageExtensionOID = asn1.ObjectIdentifier{2, 5, 29, 37}
var ClientAuthenticationOID = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 2}
var SmartcardLogonOID = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 20, 2, 2}

func GenerateWindowsCSR(username, domain string) ([]byte, []byte, error) {
	key, err := crypto.GenerateRSAPrivateKey(2048)
	if err != nil {
		return nil, nil, err
	}

	keyDER := x509.MarshalPKCS1PrivateKey(key)

	san := pkix.Extension{Id: SubjectAltNameExtensionOID}
	san.Value, err = asn1.Marshal(
		SubjectAltName[UPN]{
			OtherName: otherName[UPN]{
				OID: UPNOtherNameOID,
				Value: UPN{
					Value: fmt.Sprintf("%s@%s", username, domain),
				},
			},
		},
	)

	var EnhancedKeyUsageExtension = pkix.Extension{
		Id: EnhancedKeyUsageExtensionOID,
		Value: func() []byte {
			val, err := asn1.Marshal([]asn1.ObjectIdentifier{
				ClientAuthenticationOID,
				SmartcardLogonOID,
			})
			if err != nil {
				panic(err)
			}
			return val
		}(),
	}

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: username},
		// We have to pass SAN and ExtKeyUsage as raw extensions because
		// crypto/x509 doesn't support what we need:
		// - x509.ExtKeyUsage doesn't have the Smartcard Logon variant
		// - x509.CertificateRequest doesn't have OtherName SAN fields (which
		//   is a type of SAN distinct from DNSNames, EmailAddresses, IPAddresses
		//   and URIs)
		ExtraExtensions: []pkix.Extension{
			EnhancedKeyUsageExtension,
			san,
		},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, key)
	if err != nil {
		return nil, nil, err
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

	return csrPEM, keyDER, nil
}
