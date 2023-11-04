// Package core @Author Bing
// @Date 2023/11/2 9:48:00
// @Desc
package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/learnselfs/geeCache/utils"
	"io"
	"net/http"
)

type Remote struct {
	baseUrl   string
	remoteUrl []string
}

func newRemote(baseurl string, remoteUrl []string) *Remote {
	return &Remote{baseUrl: baseurl, remoteUrl: remoteUrl}
}

func (r *Remote) remoteGet(nodeUrl, group, key string) ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))
	urls := []string{nodeUrl, "/_geeCache/", group, "/", key}
	remoteUrl := utils.StringJoin(urls...)
	if len(nodeUrl)|len(group)|len(key) == 0 {
		return []byte(""), fmt.Errorf("nodeUrl cannot be empty")
	}

	result, err := http.Get(remoteUrl)
	if err != nil {
		fmt.Println(err, remoteUrl, nodeUrl)
		return []byte(""), err
	}
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		return []byte(""), errors.New("response status is not OK")
	}
	_, err = io.Copy(res, result.Body)
	if err != nil {
		return []byte(""), errors.New("response read copy is not OK")
	}
	return res.Bytes(), nil
}
