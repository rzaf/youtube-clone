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

		"-filter_complex",
		"[0:v]split=4[v1][v2][v3][v4];"+
			"[v1]scale=-2:144[v1out];"+
			"[v2]scale=-2:360[v2out];"+
			"[v3]scale=-2:480[v3out];"+
			"[v4]scale=-2:720[v4out]",

		// 144p
		"-map", "[v1out]", "-map", "a:0", "-c:v:0", "libx264", "-b:v:0", "95k", "-maxrate:v:0", "110k", "-bufsize:v:0", "200k", "-c:a:0", "aac", "-b:a:0", "48k",

		// 360p
		"-map", "[v2out]", "-map", "a:0", "-c:v:1", "libx264", "-b:v:1", "800k", "-maxrate:v:1", "856k", "-bufsize:v:1", "1200k", "-c:a:1", "aac", "-b:a:1", "96k",

		// 480p
		"-map", "[v3out]", "-map", "a:0", "-c:v:2", "libx264", "-b:v:2", "1400k", "-maxrate:v:2", "1500k", "-bufsize:v:2", "2100k", "-c:a:2", "aac", "-b:a:2", "128k",

		// 720p
		"-map", "[v4out]", "-map", "a:0", "-c:v:3", "libx264", "-b:v:3", "2800k", "-maxrate:v:3", "2996k", "-bufsize:v:3", "4200k", "-c:a:3", "aac", "-b:a:3", "128k",

		"-f", "hls",
		"-hls_time", "6",
		"-hls_flags", "temp_file",
		"-hls_playlist_type", "vod",

		"-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2 v:3,a:3",
		"-master_pl_name", url,
		"-hls_segment_filename", "storage/videos/"+url+".%v.%03d.ts",
		"-hls_base_url", os.Getenv("URL")+"/api/videos/",
		"storage/videos/"+url+".%v",
	)

	fmt.Println("command:", cmd.String())
	stdErr := os.Stderr
	stdOut := os.Stdout
	cmd.Stderr = stdErr
	cmd.Stdout = stdOut
	err := cmd.Run()
	println(err)
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
		"-b:a", "128k",
		// "-f", "segment",
		// "-segment_time", "10",
		// // "-segment_warp", "5",
		// "-segment_list", "storage/"+url+".m3u8",
		// "-segment_format", "mpegts",

		"-f", "hls",
		"-hls_playlist_type", "vod",
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
