package stealth

import (
	"crypto/rand"
	"crypto/tls"
	"math/big"
)

// TLSFingerprint masks TLS fingerprints to match real browsers
type TLSFingerprint struct {
	browserProfiles map[string]*TLSProfile
}

// TLSProfile defines TLS configuration for a browser
type TLSProfile struct {
	cipherSuites    []uint16
	supportedCurves []tls.CurveID
	supportedPoints []uint8
	extensions      []uint16
}

// NewTLSFingerprint creates a new TLS fingerprint masker
func NewTLSFingerprint() *TLSFingerprint {
	tf := &TLSFingerprint{
		browserProfiles: make(map[string]*TLSProfile),
	}

	// Chrome TLS profile
	tf.browserProfiles["chrome"] = &TLSProfile{
		cipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		supportedCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		supportedPoints: []uint8{0}, // uncompressed
		extensions: []uint16{
			0x0000, // server_name
			0x0005, // status_request
			0x000a, // supported_groups
			0x000b, // ec_point_formats
			0x000d, // signature_algorithms
			0x0010, // application_layer_protocol_negotiation
			0x0012, // signed_certificate_timestamp
			0x0017, // extended_master_secret
			0x002b, // supported_versions
			0x002d, // psk_key_exchange_modes
			0x0033, // key_share
		},
	}

	// Safari TLS profile
	tf.browserProfiles["safari"] = &TLSProfile{
		cipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		},
		supportedCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		supportedPoints: []uint8{0},
		extensions: []uint16{
			0x0000, // server_name
			0x0005, // status_request
			0x000a, // supported_groups
			0x000b, // ec_point_formats
			0x000d, // signature_algorithms
			0x0010, // application_layer_protocol_negotiation
			0x0017, // extended_master_secret
			0x002b, // supported_versions
		},
	}

	return tf
}

// GetTLSConfig returns a TLS config matching the specified browser
func (tf *TLSFingerprint) GetTLSConfig(browserType string) *tls.Config {
	profile, ok := tf.browserProfiles[browserType]
	if !ok {
		profile = tf.browserProfiles["chrome"] // Default to Chrome
	}

	return &tls.Config{
		CipherSuites:             profile.cipherSuites,
		PreferServerCipherSuites: false,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13,
		CurvePreferences:         profile.supportedCurves,
	}
}

// GetRandomBrowser returns a random browser type
func (tf *TLSFingerprint) GetRandomBrowser() string {
	browsers := []string{"chrome", "safari"}
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(browsers))))
	return browsers[int(n.Int64())]
}
