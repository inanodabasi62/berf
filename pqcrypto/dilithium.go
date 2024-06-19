package pqcrypto

import (
	"crypto/rand"
	"fmt"

	"github.com/cloudflare/circl/dilithium"
)

// Anahtar çifti oluşturma
func GenerateDilithiumKeypair() ([]byte, []byte, error) {
	publicKey, privateKey, err := dilithium.NewKeys(rand.Reader)
	return publicKey[:], privateKey[:], err
}

// Mesaj imzalama
func SignMessage(privateKey, message []byte) ([]byte, error) {
	return dilithium.Sign(nil, message, privateKey)
}

// İmzayı doğrulama
func VerifySignature(publicKey, message, signature []byte) bool {
	return dilithium.Verify(message, signature, publicKey) == nil
}

func main() {
	pubKey, privKey, err := GenerateDilithiumKeypair()
	if err != nil {
		fmt.Println("Error generating keypair:", err)
		return
	}
	fmt.Println("Public Key:", pubKey)
	fmt.Println("Private Key:", privKey)

	message := []byte("Hello, Berf!")
	signature, err := SignMessage(privKey, message)
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}
	fmt.Println("Signature:", signature)

	isValid := VerifySignature(pubKey, message, signature)
	fmt.Println("Is signature valid?", isValid)
}
