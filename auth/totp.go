package auth

import (
	"fmt"

	"github.com/pquerna/otp/totp"
)

// TOTP anahtarı oluşturma
func GenerateTOTPKey() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Berf",
		AccountName: "user@example.com",
	})
	if err != nil {
		return "", err
	}
	return key.URL(), nil
}

// TOTP kodunu doğrulama
func ValidateTOTPCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

func main() {
	// TOTP anahtarı oluşturma
	key, err := GenerateTOTPKey()
	if err != nil {
		fmt.Println("Error generating TOTP key:", err)
		return
	}
	fmt.Println("TOTP Key URL:", key)

	// TOTP kodunu doğrulama (Authenticator uygulamasında üretilen kodu kullanın)
	isValid := ValidateTOTPCode("YOUR_SECRET", "123456")
	fmt.Println("Is code valid?", isValid)
}
