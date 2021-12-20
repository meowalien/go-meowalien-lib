package password

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/random"
	"go.uber.org/zap/buffer"
	"golang.org/x/crypto/bcrypt"
)

const RecommendSaltLength = 24

var bffPool = buffer.NewPool()

// HashPassword will hash the given password with salt and Pepper, then return the hashed password and salt.
func HashPassword(password string, pepper string, saltLength int) (string, string, error) {
	if password == "" {
		return "", "", fmt.Errorf("the password should not be empty")
	}
	if pepper == "" {
		return "", "", fmt.Errorf("the pepper should not be empty")
	}
	if saltLength <= 0 {
		return "", "", fmt.Errorf("the salt length should not be less than or equal to 0")
	}
	randomSalt := random.RandomString(saltLength)


	bf := bffPool.Get()
	defer bf.Free()
	defer bf.Reset()

	var err error
	_ , err = bf.WriteString(password)
	if err != nil{
		return "" , "" , fmt.Errorf("error when WriteString - password :%w" , err)
	}
	_ , err = bf.WriteString(pepper)
	if err != nil{
		return "" , "" , fmt.Errorf("error when WriteString - pepper :%w" , err)
	}
	_ , err = bf.WriteString(randomSalt)
	if err != nil{
		return "" , "" , fmt.Errorf("error when WriteString - randomSalt :%w" , err)
	}
	bf.Free()
	bytes, err := bcrypt.GenerateFromPassword(bf.Bytes(), bcrypt.DefaultCost)
	return string(bytes), randomSalt, err
}

// CheckPasswordHash check if the given password, hashedPassword, and salt match.
func CheckPasswordHash(password, hash, pepper, salt string, passwordEncryption bool) bool {
	if password == "" || ((pepper == "" || salt == "") && passwordEncryption) {
		return false
	}
	if !passwordEncryption {

		return password == hash
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), append([]byte(password), pepper+salt...))
	return err == nil
}
