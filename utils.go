package utils

import (
	"errors"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/gin-gonic/gin"
)

//ErrRespWithJSON outputs error message to client
func ErrRespWithJSON(ctx *gin.Context, code int, message string) {
	if gin.Mode() == "debug" {
		resp := map[string]string{"code": strconv.Itoa(code), "error": message}
		log.Println(resp)
		ctx.JSON(code, resp)
		ctx.AbortWithStatus(code)
	} else {
		ctx.AbortWithStatus(code)
	}

}

//ExtractToken extracts token from header
func ExtractToken(header http.Header) (string, error) {
	authorization := header.Get("Authorization")
	if authorization == "" {
		return "", errors.New("Bad Request")
	}
	splits := strings.SplitN(authorization, " ", 2)
	if !(len(splits) == 2 && splits[0] == "Bearer") {
		return "", errors.New("Invalid authentication")
	}
	if splits[1] == "" {
		return "", errors.New("Invalid authentication")
	}
	return splits[1], nil
}

//GetURLPath generates URL format string
func GetURLPath(original string) (string, error) {
	regex, err := regexp.Compile(`[^a-zA-Z0-9\s]+`)
	if err != nil {
		return "", err
	}
	original = regex.ReplaceAllString(original, "")
	array := strings.Split(original, " ")
	var tmp []string
	for _, v := range array {
		if v != "" {
			tmp = append(tmp, v)
		}
	}
	return strings.Join(tmp, "-"), nil
}

//GetRandomString generates rendom string
func GetRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

//GetDistance gets distance in Metres by two coordination
func GetDistance(lat1, long1, lat2, long2 float64) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = long1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = long2 * math.Pi / 180
	r = 6378100
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
