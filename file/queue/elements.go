package queue

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	pbHelper "youtube-clone/database/pbs/helper"
	"youtube-clone/file/models"
	// "github.com/h2non/bimg"
)

type FileFormat struct {
	Url string
	pbHelper.MediaType
}

func NewFileFormat(url string, t pbHelper.MediaType) *FileFormat {
	return &FileFormat{
		Url:       url,
		MediaType: t,
	}
}

func (u *FileFormat) process() {
	switch u.MediaType {
	case pbHelper.MediaType_VIDEO:
		u.processVideo()
	case pbHelper.MediaType_MUSIC:
		u.processMusic()
	case pbHelper.MediaType_PHOTO:
		u.processPhoto()
	}
}

func (u *FileFormat) processVideo() {
	url := u.Url
	cmd := exec.Command("ffmpeg",
		"-i", "storage/temp/"+url,
		// "-vcodec", "libx265",
		// "-c:v", "h264",
		"-c:v", "libx265",
		"-preset:v", "ultrafast",
		// "-b:v", "500k",
		"-b:a", "80k",
		"-crf", "35",
		"-f", "hls",
		"-hls_playlist_type", "vod",
		// "-hls_list_size", "0",
		"-hls_allow_cache", "1",
		"-hls_time", "15",
		"-hls_flags", "independent_segments",
		"-hls_segment_type", "mpegts",
		"-hls_segment_filename", "storage/videos/"+url+"_%03d.ts",
		"-hls_base_url", os.Getenv("URL")+"/api/videos/",
		"storage/videos/"+url,
	)
	fmt.Println("command:", cmd.String())
	stdErr := os.Stderr
	stdOut := os.Stdout
	cmd.Stderr = stdErr
	cmd.Stdout = stdOut
	err := cmd.Run()
	if err != nil {
		err2, _ := err.(*exec.ExitError)
		if err2 != nil {
			fmt.Println(err2.Stderr)
			panic(err2)
		}
		panic(err)
	}
	models.SetUrlState(u.Url, models.Ready)
}

func (u *FileFormat) processMusic() {
	url := u.Url
	cmd := exec.Command("ffmpeg",
		"-i", "storage/temp/"+url,
		"-map", "0:a:0",
		"-c:a", "aac",
		"-b:a", "96k",
		// "-f", "segment",
		// "-segment_time", "10",
		// // "-segment_warp", "5",
		// "-segment_list", "storage/"+url+".m3u8",
		// "-segment_format", "mpegts",

		"-f", "hls",
		"-hls_playlist_type", "vod",
		// "-hls_list_size", "0",
		"-hls_allow_cache", "1",
		"-hls_time", "15",
		"-hls_flags", "independent_segments",
		"-hls_segment_type", "mpegts",
		"-hls_segment_filename", "storage/musics/"+url+"_%03d.ts",
		"-hls_base_url", os.Getenv("URL")+"/api/musics/",
		"storage/musics/"+url,
	)
	fmt.Println(cmd.String())
	// var errBuffer []byte
	stdErr := os.Stderr
	stdOut := os.Stdout
	cmd.Stderr = stdErr
	cmd.Stdout = stdOut
	err := cmd.Run()
	if err != nil {
		err2, _ := err.(*exec.ExitError)
		if err2 != nil {
			fmt.Println(err2.Stderr)
			panic(err2)
		}
		panic(err)
	}
	models.SetUrlState(u.Url, models.Ready)
}

func (u *FileFormat) processPhoto() {
	file, err := os.Open("storage/temp/" + u.Url)
	if err != nil {
		panic(err)
	}
	c, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	// newContent, err := bimg.NewImage(c).Process(bimg.Options{Quality: 5})
	// if err != nil {
	// 	panic(err)
	// }
	newFile, err := os.Create("storage/photos/" + u.Url)
	if err != nil {
		panic(err)
	}
	_, err = newFile.Write(c)
	if err != nil {
		panic(err)
	}

	models.SetUrlState(u.Url, models.Ready)
}

func (u *FileFormat) remove() {
	err := os.Remove("storage/temp/" + u.Url)
	if err != nil {
		panic(err)
	}
}
