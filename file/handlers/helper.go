package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"

	pbHelper "github.com/rzaf/youtube-clone/database/pbs/helper"

	// "os"
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

func CreateAndWriteUrl(src io.Reader, url string, t pbHelper.MediaType, size int64, userId int64) {
	fmt.Printf("createing file with url:%s and size:%d", url, size)
	newFile, err := os.Create("storage/temp/" + url)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	if _, err := io.Copy(newFile, src); err != nil {
		panic(err)
	}
	models.CreateUrl(url, size, t, userId)
}

// func getStoragePath() string {
// 	return ""
// }
