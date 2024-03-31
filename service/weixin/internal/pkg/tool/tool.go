package tool

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func TimeToString(layout string, time time.Time) string {
	return time.Format(layout)
}

func StringToTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}

func Marshal(i interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	return buffer.Bytes(), err
}

func GetToken() string {
	timestamp := strconv.Itoa(time.Now().Nanosecond())
	hash := md5.New()
	hash.Write([]byte(timestamp))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetRandCode(data string) string {
	b := make([]byte, 4)
	_, err := rand.Read(b)

	hash := sha3.New224()
	hash.Write([]byte(data))

	bytes := hash.Sum(nil)

	if err != nil {
		return hex.EncodeToString(bytes)
	} else {
		return hex.EncodeToString(bytes) + fmt.Sprintf("%x", b)
	}
}

func Decimal(num float64, places int32) float64 {
	num, _ = decimal.NewFromFloat(num).Round(places).Float64()

	return num
}

func GetMonthStartTimeAndEndTime(data string) (startTime, endTime time.Time) {
	monthTime, err := StringToTime("2006-01", data)

	if err != nil {
		return
	}

	year, month, _ := monthTime.Date()

	startTime = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	nextMonthTime := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	endTime = nextMonthTime.Add(-time.Second)

	return
}

func FormatPhone(phone string) string {
	if len(phone) < 10 {
		return phone
	}

	return phone[:3] + "*" + phone[7:]
}

func GetMd5(data string) string {
	hash := md5.Sum([]byte(data))

	return hex.EncodeToString(hash[:])
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
