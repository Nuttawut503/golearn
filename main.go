package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	signedStr := "sadsadsa"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Minute * 5).Unix(),
		"data": map[string]interface{}{"name": "Josh"},
	})
	fmt.Println(token)
	ss, err := token.SignedString([]byte(signedStr))
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(ss)
	v, err := jwt.Parse(ss, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected method (%v)", t.Header["alg"])
		}
		return []byte(signedStr), nil
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
	claims := v.Claims.(jwt.MapClaims)
	fmt.Println(claims["data"].(map[string]interface{})["name"].(string))
	fmt.Println(os.Getenv("Test"))
}
