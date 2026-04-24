package src

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

func parseInt(v interface{}) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	default:
		return 0, false
	}
}

func GetMessageTime(message Message) int64 {
	if message.Time != nil {
		return *message.Time
	}
	return 0
}

func SanitizeFilename(name string) string {
	safe := ""
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '-' || r == ' ' {
			safe += string(r)
		}
	}
	for len(safe) > 0 && safe[len(safe)-1] == '.' {
		safe = safe[:len(safe)-1]
	}
	return strings.TrimSpace(safe)
}

func DownloadPhoto(baseURL, photoToken string, photoID int, downloadPath string, userAgent string) string {
	urlStr := fmt.Sprintf("%s&sig=%s", baseURL, photoToken)
	filePath := filepath.Join(downloadPath, "images", fmt.Sprintf("%d.webp", photoID))

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		Logf("Failed to create request for photo %d: %v", photoID, err)
		return ""
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		Logf("Failed to download photo %d: %v", photoID, err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Logf("Failed to download image: HTTP %d", resp.StatusCode)
		return ""
	}

	file, err := os.Create(filePath)
	if err != nil {
		Logf("Failed to create file for photo %d: %v", photoID, err)
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		Logf("Failed to save photo %d: %v", photoID, err)
		return ""
	}

	Logf("Image downloaded: %s", filePath)
	return filePath
}

func DownloadVideo(urlStr string, videoID int, downloadPath string, videoHeaders string, userAgent string) string {
	filePath := filepath.Join(downloadPath, "videos", fmt.Sprintf("%d.mp4", videoID))

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		Logf("Failed to parse video URL %d: %v", videoID, err)
		return ""
	}

	client := &http.Client{Timeout: 120 * time.Second}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		Logf("Failed to create request for video %d: %v", videoID, err)
		return ""
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Host", parsedURL.Host)
	headers := videoHeaders
	for _, line := range strings.Split(headers, "\n") {
		if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		Logf("Failed to download video %d: %v", videoID, err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Logf("Failed to download video: HTTP %d", resp.StatusCode)
		return ""
	}

	file, err := os.Create(filePath)
	if err != nil {
		Logf("Failed to create file for video %d: %v", videoID, err)
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		Logf("Failed to save video %d: %v", videoID, err)
		return ""
	}

	Logf("Video downloaded: %s", filePath)
	return filePath
}

func DownloadFile(urlStr string, fileID int, fileName string, downloadPath string, userAgent string) string {
	safeName := SanitizeFilename(fileName)
	if safeName == "" {
		safeName = fmt.Sprintf("file-%d", fileID)
	}
	filePath := filepath.Join(downloadPath, "files", fmt.Sprintf("%d-%s", fileID, safeName))

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		Logf("Failed to create request for file %d: %v", fileID, err)
		return ""
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		Logf("Failed to download file %d: %v", fileID, err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Logf("Failed to download file: HTTP %d", resp.StatusCode)
		return ""
	}

	file, err := os.Create(filePath)
	if err != nil {
		Logf("Failed to create file for file %d: %v", fileID, err)
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		Logf("Failed to save file %d: %v", fileID, err)
		return ""
	}

	Logf("File downloaded: %s", filePath)
	return filePath
}

func DownloadAudio(urlStr string, audioID int, downloadPath string, audioHeaders string, userAgent string) string {
	filePath := filepath.Join(downloadPath, "audio", fmt.Sprintf("%d.mp3", audioID))

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		Logf("Failed to parse audio URL %d: %v", audioID, err)
		return ""
	}

	client := &http.Client{Timeout: 120 * time.Second}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		Logf("Failed to create request for audio %d: %v", audioID, err)
		return ""
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Host", parsedURL.Host)
	headers := audioHeaders
	for _, line := range strings.Split(headers, "\n") {
		if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		Logf("Failed to download audio %d: %v", audioID, err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Logf("Failed to download audio: HTTP %d", resp.StatusCode)
		return ""
	}

	file, err := os.Create(filePath)
	if err != nil {
		Logf("Failed to create file for audio %d: %v", audioID, err)
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		Logf("Failed to save audio %d: %v", audioID, err)
		return ""
	}

	Logf("Audio downloaded: %s", filePath)
	return filePath
}

func CountVisibleCharacters(text string) int {
	count := 0
	inTag := false
	for _, r := range text {
		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
		} else if !inTag {
			count++
		}
	}
	return count
}

func TruncateMessage(text string, isCaption bool) (string, bool) {
	maxLen := 4096
	if isCaption {
		maxLen = 1024
	}

	visibleLen := CountVisibleCharacters(text)
	if visibleLen <= maxLen {
		return text, false
	}

	targetLen := maxLen - 3
	result := ""
	count := 0
	inTag := false
	var tagBuffer strings.Builder

	for _, r := range text {
		if r == '<' {
			inTag = true
			tagBuffer.Reset()
			tagBuffer.WriteRune(r)
		} else if r == '>' {
			inTag = false
			tagBuffer.WriteRune(r)
			result += tagBuffer.String()
		} else if inTag {
			tagBuffer.WriteRune(r)
		} else {
			if count >= targetLen {
				result += closeOpenTags(result) + "..."
				return result, true
			}
			result += string(r)
			count++
		}
	}

	return result, false
}

func closeOpenTags(text string) string {
	openTags := []string{}
	i := 0
	for i < len(text) {
		if i < len(text)-1 && text[i:i+2] == "</" {
			endIdx := strings.Index(text[i:], ">")
			if endIdx != -1 {
				tagName := strings.TrimSpace(text[i+2 : i+endIdx])
				if len(openTags) > 0 && openTags[len(openTags)-1] == tagName {
					openTags = openTags[:len(openTags)-1]
				}
				i += endIdx + 1
			} else {
				i++
			}
		} else if text[i] == '<' {
			endIdx := strings.Index(text[i:], ">")
			if endIdx != -1 {
				tagContent := text[i+1 : i+endIdx]
				tagName := strings.Fields(tagContent)[0]
				if !strings.HasSuffix(tagContent, "/") && !strings.Contains(tagContent, "/") {
					openTags = append(openTags, tagName)
				}
				i += endIdx + 1
			} else {
				i++
			}
		} else {
			i++
		}
	}

	result := ""
	for i := len(openTags) - 1; i >= 0; i-- {
		result += "</" + openTags[i] + ">"
	}
	return result
}

func CheckAndHandleMessageLength(text string, isCaption bool, truncate bool) (string, bool) {
	if truncate {
		newText, wasTruncated := TruncateMessage(text, isCaption)
		if wasTruncated {
			Logf("Message was truncated to fit Telegram limits")
		}
		return newText, true
	}

	maxLen := 4096
	if isCaption {
		maxLen = 1024
	}

	visibleLen := CountVisibleCharacters(text)
	if visibleLen > maxLen {
		Logf("Message is too long (%d chars, max %d). Skipping send.", visibleLen, maxLen)
		return text, false
	}

	return text, true
}
