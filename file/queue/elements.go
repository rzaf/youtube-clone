package queue

import (
	"fmt"
	pbHelper "github.com/rzaf/youtube-clone/database/pbs/helper"
	"github.com/rzaf/youtube-clone/file/models"
	"io"
	"os"
	"os/exec"
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

func (u *FileFormat) process() error {
	switch u.MediaType {
	case pbHelper.MediaType_VIDEO:
		return u.processVideo()
	case pbHelper.MediaType_MUSIC:
		return u.processMusic()
	case pbHelper.MediaType_PHOTO:
		return u.processPhoto()
	}
	return nil
}

func (u *FileFormat) name() string {
	switch u.MediaType {
	case pbHelper.MediaType_VIDEO:
		return "video_format"
	case pbHelper.MediaType_MUSIC:
		return "music_format"
	case pbHelper.MediaType_PHOTO:
		return "photo_format"
	}
	return ""
}

func (u *FileFormat) url() string {
	return u.Url
}

func (u *FileFormat) processVideo() error {
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
		// err2, _ := err.(*exec.ExitError)
		// if err2 != nil {
		// 	fmt.Println(err2.Stderr)
		// 	return err2
		// }
		return err
	}
	models.SetUrlState(u.Url, models.Ready)
	removeTemp(u.Url)
	return nil
}

func (u *FileFormat) processMusic() error {
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
		// err2, _ := err.(*exec.ExitError)
		// if err2 != nil {
		// 	fmt.Println(err2.Stderr)
		// 	return err2
		// }
		return err
	}
	models.SetUrlState(u.Url, models.Ready)
	removeTemp(u.Url)
	return nil
}

func (u *FileFormat) processPhoto() error {
	file, err := os.Open("storage/temp/" + u.Url)
	if err != nil {
		return err
	}
	c, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	// newContent, err := bimg.NewImage(c).Process(bimg.Options{Quality: 5})
	// if err != nil {
	// return err
	// }
	newFile, err := os.Create("storage/photos/" + u.Url)
	if err != nil {
		return err
	}
	_, err = newFile.Write(c)
	if err != nil {
		return err
	}

	models.SetUrlState(u.Url, models.Ready)
	removeTemp(u.Url)
	return nil
}

func removeTemp(url string) {
	err := os.Remove("storage/temp/" + url)
	if err != nil {
		fmt.Printf("error occured while removing queue element . err:%v \n", err)
	}
}
