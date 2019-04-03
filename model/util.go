package model

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"unsafe"
)

// PathExists check path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// HTTPGetWithCache call http.GET with cache
func HTTPGetWithCache(uri, cacheDir string) (res []byte, err error) {
	exist, err := PathExists(cacheDir)
	if err != nil {
		return
	}

	// load from cache
	if exist {
		if res, err = loadCache(uri, cacheDir); len(res) > 0 {
			return
		}
	}

	resp, err := http.Get(uri)
	if err != nil {
		return
	}

	defer func() {
		err = resp.Body.Close()
	}()

	res, err = ioutil.ReadAll(resp.Body)

	// dump data into local cache
	go dumpCache(res, uri, cacheDir)
	return
}

// MD5Sum md5sum string
func MD5Sum(str string) string {
	h := md5.New()
	h.Write(stringToByte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// dumpCache dump data into cache
func dumpCache(data []byte, label, cacheDir string) {
	label = MD5Sum(label)
	cacheFile := path.Join(cacheDir, label)
	f, err := os.Create(cacheFile)
	if err != nil {
		Logger.Println("dumpCache os.Create(cacheFile): ", err)
		return
	}

	// gzip data
	wr := gzip.NewWriter(f)
	defer func() {
		err = wr.Flush()
		if err != nil {
			Logger.Println("dumpCache wr.Flush(): ", err)
			return
		}

		err = wr.Close()
		if err != nil {
			Logger.Println("dumpCache wr.Close(data): ", err)
			return
		}
	}()

	_, err = wr.Write(data)
	if err != nil {
		Logger.Println("dumpCache wr.Write(data): ", err)
		return
	}

	return
}

// load data from a local file cache
func loadCache(label, cacheDir string) (data []byte, err error) {
	label = MD5Sum(label)
	cacheFile := path.Join(cacheDir, label)
	f, err := os.Open(cacheFile)
	if err != nil {
		return nil, err
	}

	// gzip reader
	gzr, err := gzip.NewReader(f)
	if err != nil {
		Logger.Println("loadCache zlib.NewReader(f):", err)
		return nil, err
	}

	return ioutil.ReadAll(gzr)
}

// no safe string to bytes
func stringToByte(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// no safe byte to string
func byteToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: bh.Data,
		Len:  bh.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

// compress space
func compressStr(str string) string {
	if str == "" {
		return ""
	}
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, " ")
}
