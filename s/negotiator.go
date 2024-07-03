package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"regexp"
	"strings"

	"github.com/greetinc/greet-auth-srv/entity"

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
// 		clientURL = "http://localhost:8080"
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

var CountryCodes = map[string]string{
	"Afghanistan":                      "+93",
	"Albania":                          "+355",
	"Algeria":                          "+213",
	"American Samoa":                   "+1-684",
	"Andorra":                          "+376",
	"Angola":                           "+244",
	"Anguilla":                         "+1-264",
	"Antarctica":                       "+672",
	"Antigua and Barbuda":              "+1-268",
	"Argentina":                        "+54",
	"Armenia":                          "+374",
	"Aruba":                            "+297",
	"Australia":                        "+61",
	"Austria":                          "+43",
	"Azerbaijan":                       "+994",
	"Bahamas":                          "+1-242",
	"Bahrain":                          "+973",
	"Bangladesh":                       "+880",
	"Barbados":                         "+1-246",
	"Belarus":                          "+375",
	"Belgium":                          "+32",
	"Belize":                           "+501",
	"Benin":                            "+229",
	"Bermuda":                          "+1-441",
	"Bhutan":                           "+975",
	"Bolivia":                          "+591",
	"Bosnia and Herzegovina":           "+387",
	"Botswana":                         "+267",
	"Brazil":                           "+55",
	"British Indian Ocean Territory":   "+246",
	"British Virgin Islands":           "+1-284",
	"Brunei":                           "+673",
	"Bulgaria":                         "+359",
	"Burkina Faso":                     "+226",
	"Burundi":                          "+257",
	"Cambodia":                         "+855",
	"Cameroon":                         "+237",
	"Canada":                           "+1",
	"Cape Verde":                       "+238",
	"Cayman Islands":                   "+1-345",
	"Central African Republic":         "+236",
	"Chad":                             "+235",
	"Chile":                            "+56",
	"China":                            "+86",
	"Christmas Island":                 "+61",
	"Cocos Islands":                    "+61",
	"Colombia":                         "+57",
	"Comoros":                          "+269",
	"Cook Islands":                     "+682",
	"Costa Rica":                       "+506",
	"Croatia":                          "+385",
	"Cuba":                             "+53",
	"Curacao":                          "+599",
	"Cyprus":                           "+357",
	"Czech Republic":                   "+420",
	"Democratic Republic of the Congo": "+243",
	"Denmark":                          "+45",
	"Djibouti":                         "+253",
	"Dominica":                         "+1-767",
	"Dominican Republic":               "+1-809, 1-829, 1-849",
	"East Timor":                       "+670",
	"Ecuador":                          "+593",
	"Egypt":                            "+20",
	"El Salvador":                      "+503",
	"Equatorial Guinea":                "+240",
	"Eritrea":                          "+291",
	"Estonia":                          "+372",
	"Ethiopia":                         "+251",
	"Falkland Islands":                 "+500",
	"Faroe Islands":                    "+298",
	"Fiji":                             "+679",
	"Finland":                          "+358",
	"France":                           "+33",
	"French Polynesia":                 "+689",
	"Gabon":                            "+241",
	"Gambia":                           "+220",
	"Georgia":                          "+995",
	"Germany":                          "+49",
	"Ghana":                            "+233",
	"Gibraltar":                        "+350",
	"Greece":                           "+30",
	"Greenland":                        "+299",
	"Grenada":                          "+1-473",
	"Guam":                             "+1-671",
	"Guatemala":                        "+502",
	"Guernsey":                         "+44-1481",
	"Guinea":                           "+224",
	"Guinea-Bissau":                    "+245",
	"Guyana":                           "+592",
	"Haiti":                            "+509",
	"Honduras":                         "+504",
	"Hong Kong":                        "+852",
	"Hungary":                          "+36",
	"Iceland":                          "+354",
	"India":                            "+91",
	"Indonesia":                        "+62",
	"Iran":                             "+98",
	"Iraq":                             "+964",
	"Ireland":                          "+353",
	"Isle of Man":                      "+44-1624",
	"Israel":                           "+972",
	"Italy":                            "+39",
	"Ivory Coast":                      "+225",
	"Jamaica":                          "+1-876",
	"Japan":                            "+81",
	"Jersey":                           "+44-1534",
	"Jordan":                           "+962",
	"Kazakhstan":                       "+7",
	"Kenya":                            "+254",
	"Kiribati":                         "+686",
	"Kosovo":                           "+383",
	"Kuwait":                           "+965",
	"Kyrgyzstan":                       "+996",
	"Laos":                             "+856",
	"Latvia":                           "+371",
	"Lebanon":                          "+961",
	"Lesotho":                          "+266",
	"Liberia":                          "+231",
	"Libya":                            "+218",
	"Liechtenstein":                    "+423",
	"Lithuania":                        "+370",
	"Luxembourg":                       "+352",
	"Macau":                            "+853",
	"Macedonia":                        "+389",
	"Madagascar":                       "+261",
	"Malawi":                           "+265",
	"Malaysia":                         "+60",
	"Maldives":                         "+960",
	"Mali":                             "+223",
	"Malta":                            "+356",
	"Marshall Islands":                 "+692",
	"Mauritania":                       "+222",
	"Mauritius":                        "+230",
	"Mayotte":                          "+262",
	"Mexico":                           "+52",
	"Micronesia":                       "+691",
	"Moldova":                          "+373",
	"Monaco":                           "+377",
	"Mongolia":                         "+976",
	"Montenegro":                       "+382",
	"Montserrat":                       "+1-664",
	"Morocco":                          "+212",
	"Mozambique":                       "+258",
	"Myanmar":                          "+95",
	"Namibia":                          "+264",
	"Nauru":                            "+674",
	"Nepal":                            "+977",
	"Netherlands":                      "+31",
	"Netherlands Antilles":             "+599",
	"New Caledonia":                    "+687",
	"New Zealand":                      "+64",
	"Nicaragua":                        "+505",
	"Niger":                            "+227",
	"Nigeria":                          "+234",
	"Niue":                             "+683",
	"North Korea":                      "+850",
	"Northern Mariana Islands":         "+1-670",
	"Norway":                           "+47",
	"Oman":                             "+968",
	"Pakistan":                         "+92",
	"Palau":                            "+680",
	"Palestine":                        "+970",
	"Panama":                           "+507",
	"Papua New Guinea":                 "+675",
	"Paraguay":                         "+595",
	"Peru":                             "+51",
	"Philippines":                      "+63",
	"Pitcairn":                         "+64",
	"Poland":                           "+48",
	"Portugal":                         "+351",
	"Puerto Rico":                      "+1-787, 1-939",
	"Qatar":                            "+974",
	"Republic of the Congo":            "+242",
	"Reunion":                          "+262",
	"Romania":                          "+40",
	"Russia":                           "+7",
	"Rwanda":                           "+250",
	"Saint Barthelemy":                 "+590",
	"Saint Helena":                     "+290",
	"Saint Kitts and Nevis":            "+1-869",
	"Saint Lucia":                      "+1-758",
	"Saint Martin":                     "+590",
	"Saint Pierre and Miquelon":        "+508",
	"Saint Vincent and the Grenadines": "+1-784",
	"Samoa":                            "+685",
	"San Marino":                       "+378",
	"Sao Tome and Principe":            "+239",
	"Saudi Arabia":                     "+966",
	"Senegal":                          "+221",
	"Serbia":                           "+381",
	"Seychelles":                       "+248",
	"Sierra Leone":                     "+232",
	"Singapore":                        "+65",
	"Sint Maarten":                     "+1-721",
	"Slovakia":                         "+421",
	"Slovenia":                         "+386",
	"Solomon Islands":                  "+677",
	"Somalia":                          "+252",
	"South Africa":                     "+27",
	"South Korea":                      "+82",
	"South Sudan":                      "+211",
	"Spain":                            "+34",
	"Sri Lanka":                        "+94",
	"Sudan":                            "+249",
	"Suriname":                         "+597",
	"Svalbard and Jan Mayen":           "+47",
	"Swaziland":                        "+268",
	"Sweden":                           "+46",
	"Switzerland":                      "+41",
	"Syria":                            "+963",
	"Taiwan":                           "+886",
	"Tajikistan":                       "+992",
	"Tanzania":                         "+255",
	"Thailand":                         "+66",
	"Togo":                             "+228",
	"Tokelau":                          "+690",
	"Tonga":                            "+676",
	"Trinidad and Tobago":              "+1-868",
	"Tunisia":                          "+216",
	"Turkey":                           "+90",
	"Turkmenistan":                     "+993",
	"Turks and Caicos Islands":         "+1-649",
	"Tuvalu":                           "+688",
	"U.S. Virgin Islands":              "+1-340",
	"Uganda":                           "+256",
	"Ukraine":                          "+380",
	"United Arab Emirates":             "+971",
	"United Kingdom":                   "+44",
	"United States":                    "+1",
	"Uruguay":                          "+598",
	"Uzbekistan":                       "+998",
	"Vanuatu":                          "+678",
	"Vatican":                          "+379",
	"Venezuela":                        "+58",
	"Vietnam":                          "+84",
	"Wallis and Futuna":                "+681",
	"Western Sahara":                   "+212",
	"Yemen":                            "+967",
	"Zambia":                           "+260",
	"Zimbabwe":                         "+263",
}

func GetCountryCode(country string) (string, error) {
	code, ok := CountryCodes[country]
	if !ok {
		return "", fmt.Errorf("country code not found for %s", country)
	}
	return code, nil
}

// FormatPhoneNumber memformat nomor telepon berdasarkan kode negara
func FormatPhoneNumber(phone, country string) (string, error) {
	// Hilangkan semua spasi dan tanda penghubung dari nomor telepon
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, ",", "")
	phone = strings.ReplaceAll(phone, ".", "")
	phone = strings.ReplaceAll(phone, "e", "")

	// Ambil kode negara dari peta CountryCodes
	code, exists := CountryCodes[country]
	if !exists {
		return "", errors.New("invalid country")
	}

	// Jika nomor telepon sudah dimulai dengan kode negara, kembalikan nomor tersebut
	if strings.HasPrefix(phone, code) {
		return phone, nil
	}

	// Jika nomor telepon dimulai dengan '0', ganti '0' dengan kode negara
	if strings.HasPrefix(phone, code) {
		// Jika nomor telepon dimulai dengan '0' setelah kode negara
		if strings.HasPrefix(phone, code+"0") {
			return code + phone[len(code)+1:], nil
		}
		return phone, nil
	}

	// Jika nomor telepon dimulai dengan '+', kembalikan error
	if strings.HasPrefix(phone, "+") {
		return "", errors.New("invalid phone number format")
	}

	// Tambahkan kode negara di depan nomor telepon
	return code + phone, nil
}

func IsValidEmail(email string) bool {
	// This is a basic email validation regex, it may not cover all cases
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	return err == nil && matched
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
	// "Mailtrap: <a href='http://localhost:8080/verify/%s'>Verify Account</a>Your verification code is: <strong>%s</strong>",
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

func IsDuplicateEntryError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code.Name() == "unique_violation"
	}
	return false
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

func Haversine(coord1, coord2 entity.RadiusRange) float64 {
	const earthRadius = 6371000 // Earth radius in meters

	// Convert degrees to radians
	lat1 := degToRad(coord1.Latitude)
	lon1 := degToRad(coord1.Longitude)
	lat2 := degToRad(coord2.Latitude)
	lon2 := degToRad(coord2.Longitude)

	// Calculate differences
	dlat := lat2 - lat1
	dlon := lon2 - lon1

	// Haversine formula
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in meters
	distance := earthRadius * c
	return distance
}

// Convert degrees to radians
func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

var key = []byte("!@#$%123_pLaTfOm____aPp_123^&*()")

func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)
	blockSize := block.BlockSize()
	padding := blockSize - len(plainTextBytes)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	plainTextBytes = append(plainTextBytes, padText...)

	cipherText := make([]byte, len(plainTextBytes))
	mode := cipher.NewCBCEncrypter(block, key[:blockSize])
	mode.CryptBlocks(cipherText, plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(cipherText string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(cipherTextBytes)%blockSize != 0 {
		return "", errors.New("invalid cipher text length")
	}

	plainText := make([]byte, len(cipherTextBytes))
	mode := cipher.NewCBCDecrypter(block, key[:blockSize])
	mode.CryptBlocks(plainText, cipherTextBytes)

	padding := int(plainText[len(plainText)-1])
	if padding > blockSize || padding == 0 {
		return "", errors.New("invalid padding size")
	}

	return string(plainText[:len(plainText)-padding]), nil
}

func TruncateFullName(name string, maxLen int) string {
	if len(name) > maxLen {
		return name[:maxLen] + "..."
	}
	return name
}
