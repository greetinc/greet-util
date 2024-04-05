package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/gomail.v2"
)

// func NewMailer() *gomail.Dialer {
// 	// Replace the following information with your Mailtrap.io credentials
// 	username := "7de3a28724e886"
// 	password := "353081a2c62514"
// 	host := "smtp.mailtrap.io"
// 	port := 587

// 	// Create a new dialer with Mailtrap.io settings
// 	dialer := gomail.NewDialer(host, port, username, password)

// 	return dialer
// }

// func GetClientURL() string {
// 	// Fetch the clientURL from environment variable or configuration
// 	clientURL := os.Getenv("CLIENT_URL")
// 	if clientURL == "" {
// 		// Provide a default value or handle the missing configuration accordingly
// 		clientURL = "http://localhost:2345"
// 	}
// 	return clientURL
// }

// func SendResetEmail(email, resetLink string) error {
// 	message := gomail.NewMessage()

// 	message.SetHeader("From", "aseprayana95@gmail.com")
// 	message.SetHeader("To", email)
// 	message.SetHeader("Subject", "Password Reset")
// 	message.SetBody("text/html", fmt.Sprintf("Click <a href='%s'>here</a> to reset your password.", resetLink))

// 	dialer := NewMailer()

// 	// Send the email
// 	if err := dialer.DialAndSend(message); err != nil {
// 		return err
// 	}

// 	return nil
// }

func IsValidEmail(email string) bool {
	// This is a basic email validation regex, it may not cover all cases
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	return err == nil && matched
}

func IsDuplicateEntryError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code.Name() == "unique_violation"
	}
	return false
}
func Mailtrap(to, otp string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "aseprayana95@gmail.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", "Test Email")
	//if using otp kode
	mailer.SetBody("text/html", fmt.Sprintf("Your verification code is: <strong>%s</strong>", otp))
	//click button at email and verify
	// mailer.SetBody("text/html", fmt.Sprintf("Hello, this is a test email from "+
	// "Mailtrap: <a href='http://localhost:2345/verify/%s'>Verify Account</a>Your verification code is: <strong>%s</strong>",
	// verificationToken, otp))

	dialer := gomail.NewDialer("smtp.mailtrap.io", 587, "126a08e0c5ff69", "9f8b22657e0257")
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Use this only for development, not secure for production

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}

func Negotiate(c echo.Context, code int, i interface{}) error {
	mediaType := c.QueryParam("mediaType")

	switch mediaType {
	case "xml":
		return c.XML(code, i)
	case "json":
		return c.JSON(code, i)
	default:
		return c.JSON(code, i)
	}
}

// file type support
func ValidateFileType(filePath string) error {
	// Periksa tipe file menggunakan library filetype
	buf := make([]byte, 512)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Read(buf)
	if err != nil {
		return err
	}

	// Menggunakan Match dengan jenis file yang diizinkan
	kind, unknown := filetype.Match(buf)
	if unknown != nil {
		return errors.New("unknown file type")
	}

	// Kumpulan jenis file yang diizinkan (misal: "image/jpeg", "image/png")
	allowedTypes := []string{
		"image/jpeg",
		"image/png",
		"image/jpg",

		// tambahkan jenis file lain yang diizinkan
	}

	// Validasi jenis file yang diizinkan
	for _, allowedType := range allowedTypes {
		if kind.MIME.Value == allowedType {
			return nil // File adalah jenis yang diizinkan
		}
	}

	// Jenis file tidak diizinkan
	return errors.New("invalid file type")
}

func GetVerificationLink(token string) string {
	baseURL := "http://localhost" // Ganti dengan URL aplikasi Anda
	return fmt.Sprintf("%s/verify/%s", baseURL, token)
}

func GenerateRandomString() string {
	randomUUID := uuid.New()
	return randomUUID.String()
}

func FormatIDR(value float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%s", currency.IDR.Amount(value))
}

func CalculateJumlahBunga(principal float64, rate float64, time int) float64 {
	// Formula: Bunga = Principal * Rate * Time
	bunga := principal * rate * float64(time)
	return bunga
}

// saveFile
func SaveFile(fileHeader *multipart.FileHeader, destination string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("gagal membuat file: %v", err)
	}
	defer func() {
		cerr := dst.Close()
		if err == nil && cerr != nil {
			err = fmt.Errorf("error saat menutup file: %v", cerr)
		}
	}()

	_, err = io.Copy(dst, src)
	if err != nil {
		// Hapus file yang ter-copy sebagian jika terjadi kesalahan
		_ = os.Remove(destination)
		return fmt.Errorf("gagal menyalin file: %v", err)
	}

	return nil
}

// Encrypt a string using AES-GCM
const aesKey = "your-secret-key-32-bytes" // Pastikan untuk mengganti kunci ini dengan yang kuat dan rahasia

func EncryptFileName(originalFileName string) (string, error) {
	key, err := hex.DecodeString(aesKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encryptedFileName := gcm.Seal(nonce, nonce, []byte(originalFileName), nil)
	return hex.EncodeToString(encryptedFileName), nil
}

// Convert degrees to radians
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
