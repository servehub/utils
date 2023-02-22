package utils

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	if pos < 0 {
		pos = 0
	}
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func Contains(elm string, list []string) bool {
	for _, v := range list {
		if v == elm {
			return true
		}
	}
	return false
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func MergeMaps(maps ...map[string]string) map[string]string {
	out := make(map[string]string, 0)
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

func MapsEqual(a, b map[string]string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}

type BySortIndex []map[string]interface{}

func (a BySortIndex) Len() int { return len(a) }
func (a BySortIndex) Less(i, j int) bool {
	return fmt.Sprintf("%v", a[i]["sortIndex"]) < fmt.Sprintf("%v", a[j]["sortIndex"])
}
func (a BySortIndex) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

var allLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var RandomString = func(length uint) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = allLetters[seededRand.Intn(len(allLetters))]
	}
	return string(b)
}

func StripLeftMargin(data string) string {
	lines := strings.Split(strings.TrimSpace(strings.Replace(data, "\t", "  ", -1)), "\n")
	minPrefix := len(data)

	for _, l := range lines[1:] {
		if subs := len(strings.TrimSpace(l)); subs > 0 {
			if curPrefix := len(l) - subs; curPrefix < minPrefix {
				minPrefix = curPrefix
			}
		}
	}

	for i, l := range lines {
		if len(l) > minPrefix {
			lines[i] = l[minPrefix:]
		}
	}

	return strings.Join(lines, "\n")
}

func WriteTemp(data []byte, callback func(string) error) error {
	tmpfile, err := ioutil.TempFile("", "serve-")
	if err != nil {
		return fmt.Errorf("Error create tmpfile: %v", err)
	}

	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	if _, err := tmpfile.Write(data); err != nil {
		return fmt.Errorf("Error write to tmpfile: %v", err)
	}

	return callback(tmpfile.Name())
}
