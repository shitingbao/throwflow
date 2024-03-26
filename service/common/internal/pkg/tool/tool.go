package tool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/speps/go-hashids/v2"
	mrand "math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetToken() string {
	timestamp := strconv.Itoa(time.Now().Nanosecond())
	hash := md5.New()
	hash.Write([]byte(timestamp))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func TimeToString(layout string, time time.Time) string {
	return time.Format(layout)
}

func GetMd5(data string) string {
	hash := md5.Sum([]byte(data))

	return hex.EncodeToString(hash[:])
}

func Decimal(num float64, places int32) float64 {
	num, _ = decimal.NewFromFloat(num).Round(places).Float64()

	return num
}

func GetShortCode() (string, error) {
	mrand.Seed(time.Now().UnixNano())

	hdata := hashids.NewData()
	hdata.Salt = "mengma"

	hid, err := hashids.NewWithData(hdata)

	if err != nil {
		return "", err
	}

	hash, err := hid.Encode([]int{mrand.Intn(100), mrand.Intn(100), mrand.Intn(100)})

	if err != nil {
		return "", err
	}

	return hash, nil
}

func StructToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(s).Elem()
	t := reflect.TypeOf(s).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		tag := field.Tag.Get("json")

		if strings.Contains(tag, ",omitempty") && reflect.DeepEqual(value, reflect.Zero(field.Type).Interface()) {
			continue
		}

		key := strings.ToLower(field.Name[:1]) + field.Name[1:]

		result[key] = value
	}

	return result
}

func SortMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
