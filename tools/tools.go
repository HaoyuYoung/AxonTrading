package tools

import (
	"AxonTrading/base"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

func UnifiedSymbol(exchange, symbol string) string {
	//Default: Token1_Token2

	tokens := strings.Split(symbol, "_")

	switch exchange {
	default:
		return symbol
	case base.BINANCE, base.BITGET, base.BYBIT, base.MEXCV3:
		return tokens[0] + tokens[1]
	case base.BITMART, base.GATEIO, base.MEXC:
		return symbol
	case base.KUCOIN, base.OKEX:
		return tokens[0] + "-" + tokens[1]
	}
}

func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func TimeToUnix(e time.Time) int64 {
	local, _ := time.LoadLocation("Local")
	timeUnix, _ := time.ParseInLocation("2006-01-02 15:04:05", e.Format("2006-01-02 15:04:05"), local)
	return timeUnix.UnixNano() / 1e6
}

func TimeToUnix2(e time.Time) int64 {
	local, _ := time.LoadLocation("Asia/Shanghai")
	timeUnix, _ := time.ParseInLocation("2006-01-02 15:04:05", e.Format("2006-01-02 15:04:05"), local)
	return timeUnix.UnixNano()
}

func UnixToTime(e string) (dataTime time.Time, err error) {
	data, err := strconv.ParseInt(e, 10, 64)
	dataTime = time.Unix(data/1000, 0)
	return
}

// RandFloat 随机 Float
func RandFloat(min, max float64) float64 {
	rand.Seed(time.Now().Unix())
	ran := rand.Float64()
	return ran*min + (1-ran)*max
}

// RandInt 随机整数
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// RandIntWithSeed 随机整数
func RandIntWithSeed(seed, min, max int) int {
	rand.Seed(int64(seed))
	return rand.Intn(max-min) + min
}

func RandFloat64WithSeed(seed int, min, max float64) float64 {
	rand.Seed(int64(seed))
	return min + rand.Float64()*(max-min)
}

// FormatFloatCeil 舍弃的尾数不为0，强制进位
func FormatFloatCeil(num float64, decimal int) (float64, error) {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Ceil(num*d)/d, 'f', -1, 64)
	return strconv.ParseFloat(res, 64)
}

// FormatFloatFloor 强制舍弃尾数
func FormatFloatFloor(num float64, decimal int) (float64, error) {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Floor(num*d)/d, 'f', -1, 64)
	return strconv.ParseFloat(res, 64)
}

// GetMaxFloat64 获取 Float64 列表中的最大值
func GetMaxFloat64(l []float64) (max float64) {
	max = l[0]
	for _, v := range l {
		if v > max {
			max = v
		}
	}
	return
}

// GetMinFloat64 获取 Float64 列表中的最小值
func GetMinFloat64(l []float64) (min float64) {
	min = l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return
}

func splitString(r rune) bool {
	return r == '_' || r == '-'
}

// SplitStringChar 字符串按照多个字符分割
func SplitStringChar(s string) []string {
	a := strings.FieldsFunc(s, splitString)
	return a
}

// GetDecimalPlaces 从 float64 中获取小数点的位置
func GetDecimalPlaces(f float64) int {
	numstr := fmt.Sprint(f)
	tmp := strings.Split(numstr, ".")
	if len(tmp) <= 1 {
		return 0
	}
	return len(tmp[1])
}

// GetDecimalPlacesStr 从 float64 中获取小数点的位置
func GetDecimalPlacesStr(str string) int {
	tmp := strings.Split(str, ".")
	if len(tmp) <= 1 {
		return 0
	}
	return len(tmp[1])
}

func S2M(i interface{}) map[string]string {
	m := make(map[string]string)
	j, _ := json.Marshal(i)
	err := json.Unmarshal(j, &m)
	fmt.Println(err)

	return m
}

func ObjectToJson(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json 转换错误 》 ", err)
		return ""
	} else {
		return string(bytes)
	}
}

func arrayInGroupsOf(arr []int, num int64) [][]int {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]int{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]int, 0)
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}

func ReplaceStringChar(symbol string) string {
	return strings.NewReplacer("_", "", "-", "").Replace(symbol)
}
func GetSizeWithU(price, total float64, apres int) string {
	sizeF := total / price
	return strconv.FormatFloat(sizeF, 'f', apres, 64)
}
func FormatSymbol(exchange, symbol string) string {
	symbols := SplitStringChar(symbol)
	if len(symbols) == 0 {
		return ""
	}
	if len(symbols) == 1 && strings.Contains(symbol, "USDT") {
		tokens := strings.Split(symbol, "USDT")
		symbols[0] = tokens[0]
		symbols = append(symbols, "USDT")
	} else {
		return symbol
	}

	newSymbol := func(s string) string {
		return symbols[0] + s + symbols[1]
	}

	switch exchange {
	case base.KUCOIN, base.OKEX, base.PROBIT:
		return newSymbol("-")
	case base.BINANCE, base.MEXC, base.GATEIO, base.MEXCV3, base.BITMART:
		return newSymbol("_")
	case base.BYBIT, base.BITRUE:
		return newSymbol("")
	default:
		return symbol
	}
}

func SplitSymbol(symbol string) []string {
	symbols := SplitStringChar(symbol)

	if len(symbols) == 0 {
		return nil
	}

	if len(symbols) == 1 && strings.Contains(symbol, "USDT") {
		tokens := strings.Split(symbol, "USDT")
		symbols[0] = tokens[0]
		symbols = append(symbols, "USDT")
		return symbols
	}

	if len(symbols) == 2 {
		return symbols
	}

	return nil
}

func FormatSide(exchange, side string) string {
	if exchange == "" || side == "" {
		return ""
	}

	switch exchange {
	case base.GATEIO:
		if side == base.ASK {
			return "sell"
		} else if side == base.BID {
			return "buy"
		}

	case base.BITRUE:
		if side == base.ASK {
			return "SELL"
		} else if side == base.BID {
			return "BUY"
		}
	default:
		return side
	}

	return side
}

func ParseSide(exchange, side string) string {
	if exchange == "" || side == "" {
		return ""
	}

	switch exchange {
	case base.BITRUE:
		if side == "SELL" {
			return base.ASK
		} else if side == "BUY" {
			return base.BID
		}
	default:
		return side
	}

	return side
}

func Map2UrlQuery(mapParams map[string]string) string {
	var strParams string
	mapSort := []string{}
	for key := range mapParams {
		mapSort = append(mapSort, key)
	}
	sort.Strings(mapSort)

	for _, key := range mapSort {
		strParams += (key + "=" + mapParams[key] + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}
