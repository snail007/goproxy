package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha256"
	"errors"

	lbuf "github.com/snail007/goproxy/core/lib/buf"
	"github.com/Yawning/chacha20"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/cast5"
)

const leakyBufSize = 2048
const maxNBuf = 2048

var leakyBuf = lbuf.NewLeakyBuf(maxNBuf, leakyBufSize)
var errEmptyPassword = errors.New("proxy key")

func md5sum(d []byte) []byte {
	h := md5.New()
	h.Write(d)
	return h.Sum(nil)
}

func evpBytesToKey(password string, keyLen int) (key []byte) {
	const md5Len = 16
	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, md5sum([]byte(password)))

	// Repeatedly call md5 until bytes generated is enough.
	// Each call to md5 uses data: prev md5 sum + password.
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], md5sum(d))
	}
	return m[:keyLen]
}

type DecOrEnc int

const (
	Decrypt DecOrEnc = iota
	Encrypt
)

func newStream(block cipher.Block, err error, key, iv []byte,
	doe DecOrEnc) (cipher.Stream, error) {
	if err != nil {
		return nil, err
	}
	if doe == Encrypt {
		return cipher.NewCFBEncrypter(block, iv), nil
	} else {
		return cipher.NewCFBDecrypter(block, iv), nil
	}
}

func newAESCFBStream(key, iv []byte, doe DecOrEnc) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	return newStream(block, err, key, iv, doe)
}

func newAESCTRStream(key, iv []byte, doe DecOrEnc) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCTR(block, iv), nil
}

func newDESStream(key, iv []byte, doe DecOrEnc) (cipher.Stream, error) {
	block, err := des.NewCipher(key)
	return newStream(block, err, key, iv, doe)
}

func newBlowFishStream(key, iv []byte, doe DecOrEnc) (cipher.Stream, error) {
	block, err := blowfish.NewCipher(key)
	return newStream(block, err, key, iv, doe)
}

func newCast5Stream(key, iv []byte, doe DecOrEnc) (cipher.Stream, error) {
	block, err := cast5.NewCipher(key)
	return newStream(block, err, key, iv, doe)
}

func newRC4MD5Stream(key, iv []byte, _ DecOrEnc) (cipher.Stream, error) {
	h := md5.New()
	h.Write(key)
	h.Write(iv)
	rc4key := h.Sum(nil)

	return rc4.NewCipher(rc4key)
}

func newChaCha20Stream(key, iv []byte, _ DecOrEnc) (cipher.Stream, error) {
	return chacha20.NewCipher(key, iv)
}

func newChaCha20IETFStream(key, iv []byte, _ DecOrEnc) (cipher.Stream, error) {
	return chacha20.NewCipher(key, iv)
}

type cipherInfo struct {
	keyLen    int
	ivLen     int
	newStream func(key, iv []byte, doe DecOrEnc) (cipher.Stream, error)
}

var cipherMethod = map[string]*cipherInfo{
	"aes-128-cfb":   {16, 16, newAESCFBStream},
	"aes-192-cfb":   {24, 16, newAESCFBStream},
	"aes-256-cfb":   {32, 16, newAESCFBStream},
	"aes-128-ctr":   {16, 16, newAESCTRStream},
	"aes-192-ctr":   {24, 16, newAESCTRStream},
	"aes-256-ctr":   {32, 16, newAESCTRStream},
	"des-cfb":       {8, 8, newDESStream},
	"bf-cfb":        {16, 8, newBlowFishStream},
	"cast5-cfb":     {16, 8, newCast5Stream},
	"rc4-md5":       {16, 16, newRC4MD5Stream},
	"rc4-md5-6":     {16, 6, newRC4MD5Stream},
	"chacha20":      {32, 8, newChaCha20Stream},
	"chacha20-ietf": {32, 12, newChaCha20IETFStream},
}

func GetCipherMethods() (keys []string) {
	keys = []string{}
	for k := range cipherMethod {
		keys = append(keys, k)
	}
	return
}
func CheckCipherMethod(method string) error {
	if method == "" {
		method = "aes-256-cfb"
	}
	_, ok := cipherMethod[method]
	if !ok {
		return errors.New("Unsupported encryption method: " + method)
	}
	return nil
}

type Cipher struct {
	WriteStream cipher.Stream
	ReadStream  cipher.Stream
	key         []byte
	info        *cipherInfo
}

func NewCipher(method, password string) (c *Cipher, err error) {
	if password == "" {
		return nil, errEmptyPassword
	}
	mi, ok := cipherMethod[method]
	if !ok {
		return nil, errors.New("Unsupported encryption method: " + method)
	}
	key := evpBytesToKey(password, mi.keyLen)
	c = &Cipher{key: key, info: mi}
	if err != nil {
		return nil, err
	}
	//hash(key) -> read IV
	riv := sha256.New().Sum(c.key)[:c.info.ivLen]
	c.ReadStream, err = c.info.newStream(c.key, riv, Decrypt)
	if err != nil {
		return nil, err
	} //hash(read IV) -> write IV
	wiv := sha256.New().Sum(riv)[:c.info.ivLen]
	c.WriteStream, err = c.info.newStream(c.key, wiv, Encrypt)
	if err != nil {
		return nil, err
	}
	return c, nil
}
