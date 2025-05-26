package helper

import (
	"encoding/base64"
	"fmt"
	"github.com/simabdi/auth-service/model"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"unicode"
)

func JsonResponse(code int, message string, success bool, error string, data interface{}) model.Response {
	meta := model.Meta{
		Code:    code,
		Status:  success,
		Message: message,
		Error:   error,
	}

	response := model.Response{
		Meta: meta,
		Data: data,
	}

	return response
}

func Std64Encode(plainText string) string {
	return base64.StdEncoding.EncodeToString([]byte(plainText))
}

func Std64Decode(encoded string) string {
	decodedByte, _ := base64.StdEncoding.DecodeString(encoded)
	return string(decodedByte)
}

func GenerateTransactionCode(latestNumber uint) string {
	now := time.Now()
	year := now.Year() % 100
	month := int(now.Month())
	return fmt.Sprintf("BIL%02d%02d%05d", year, month, latestNumber)
}

func GenerateMedicalNumber(clinicID, latestNumber uint) string {
	now := time.Now()
	year := now.Year() % 100
	return fmt.Sprintf("%02d%04d%03d", year, clinicID, latestNumber)
}

func GeneratePurchaseCode(queueNumber uint) string {
	t := time.Now()
	return fmt.Sprintf("POC%04d%02d%02d%05d",
		t.Year(),
		int(t.Month()),
		t.Day(),
		queueNumber%100000,
	)
}

func GenerateCooperativeMemberNumber(cooperativeNumber uint, queueNumber uint) string {
	t := time.Now()
	return fmt.Sprintf("%04d%02d%02d.%04d.%04d",
		t.Year(),
		int(t.Month()),
		t.Day(),
		cooperativeNumber%10000,
		queueNumber%10000,
	)
}

func IsDigitsOnly(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return r < '0' || r > '9'
	}) == -1
}

func TitleCase(s string) string {
	return strings.Join(MapWords(s, func(word string) string {
		if len(word) == 0 {
			return word
		}
		r := []rune(word)
		r[0] = unicode.ToUpper(r[0])
		for i := 1; i < len(r); i++ {
			r[i] = unicode.ToLower(r[i])
		}
		return string(r)
	}), " ")
}

func MapWords(s string, mapper func(string) string) []string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = mapper(word)
	}
	return words
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
