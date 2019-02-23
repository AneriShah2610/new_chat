package helper

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strconv"
)

// Random String Generation
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Random integer Generation
func Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

// HashKey creation
func HashKeycreation(memberIdArray []int)(string){
	newHash := sha1.New()
	newHash.Write([]byte(strconv.Itoa(memberIdArray[0]) + strconv.Itoa(memberIdArray[1])))
	hashKey := hex.EncodeToString(newHash.Sum(nil))
	return hashKey
}