package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	pbHelper "github.com/rzaf/youtube-clone/database/pbs/helper"

	// "os"
	"github.com/rzaf/youtube-clone/file/helpers"
	"github.com/rzaf/youtube-clone/file/models"
)

func generateSecureToken(length int) string {
	buff := make([]byte, int(math.Ceil(float64(length)/2)))
	if _, err := rand.Read(buff); err != nil {
		panic(err.Error())
	}
	str := hex.EncodeToString(buff)
	return str[:length]
}

func getUniqueFileUrl() string {
	for i := 0; i < 10; i++ {
		newUrl := generateSecureToken(16)
		if models.ExistsUrl(newUrl) {
			return newUrl
		}
		// path := "storage/" + newUrl
		// if _, err := os.Stat(path); os.IsNotExist(err) {
		// 	return newUrl
		// }
	}
	panic("unable to create unique url !!!")
}

func sanitizeFileName(name string) string {
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)

	re := regexp.MustCompile(`[^a-zA-Z0-9_.\s-]`)
	name = re.ReplaceAllString(name, "")

	if len(name) > 128 {
		name = name[:128]
	}

	return name
}

func CreateAndWriteUrl(src io.Reader, url string, t pbHelper.MediaType, contentType string, originalName string, userId int64) {
	originalName = sanitizeFileName(originalName)
	fmt.Printf("creating file with url:%s and name:%d", url, originalName)
	newFile, err := os.Create("storage/temp/" + url)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	hasher := sha256.New()
	multiWriter := io.MultiWriter(newFile, hasher)
	var size int64
	if size, err = io.Copy(multiWriter, src); err != nil {
		panic(err)
	}
	checksum := hex.EncodeToString(hasher.Sum(nil))
	err2 := helpers.CheckUserUploadBandwidth(size, userId)
	if err2 != nil {
		os.Remove("storage/temp/" + url)
		panic(err2)
	}

	models.CreateUrl(url, size, t, contentType, originalName, checksum, userId)
}

// func getStoragePath() string {
// 	return ""
// }
