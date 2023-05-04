package pkg_util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	pkg_error "go.cognotif/pkg/error"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/nacl/box"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func BindRequestAndValidate(c echo.Context, request any) error {
	if err := c.Bind(request); err != nil {
		return fmt.Errorf(pkg_error.BAD_REQUEST)
	}

	if err := (&echo.DefaultBinder{}).BindHeaders(c, request); err != nil {
		return err
	}

	if err := (&echo.DefaultBinder{}).BindPathParams(c, request); err != nil {
		return err
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	return nil
}

func Atoi(a string) int {
	d, _ := strconv.Atoi(a)
	return d
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateNonce() [24]byte {
	var nonce [24]byte
	_, _ = io.ReadFull(rand.Reader, nonce[:])
	return nonce
}

func ExtractNonce(ciphertext []byte) ([]byte, [24]byte) {
	var nonce [24]byte
	_, rest := copy(nonce[:], ciphertext[:24]), ciphertext[24:]
	return rest, nonce
}

func EkstractKeypairCognotif() ([32]byte, [32]byte) {

	key, _ := base64.StdEncoding.DecodeString(os.Getenv("KEYPAIR"))

	var pub, priv [32]byte
	_ = copy(pub[:], key[:32])
	_ = copy(priv[:], key[32:])
	return pub, priv
}

func SealBox(plaintext []byte) (ciphertext []byte) {
	nonce := GenerateNonce()
	pub, priv := EkstractKeypairCognotif()
	return box.Seal(nonce[:], plaintext, &nonce, &pub, &priv)
}

func OpenBox(chipertext []byte) (plaintext []byte, ok bool) {
	rest, nonce := ExtractNonce(chipertext)
	pub, priv := EkstractKeypairCognotif()
	return box.Open(nil, rest, &nonce, &pub, &priv)
}
