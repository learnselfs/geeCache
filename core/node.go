// Package core @Author Bing
// @Date 2023/11/1 15:30:00
// @Desc
package core

import (
	"github.com/learnselfs/geeCache/utils"
	"net/http"
)

const (
	defaultReplicas = 2 << 6
	defaultSize     = 2 << 10
)

type Node struct {
	*Remote
	*Group
	*Http
	*Map
	baseUrl string
}

func NewNode(host, port string, api bool, remoteNodeUrl []string, callback callback, basePath string) *Node {

	urls := []string{host, ":", port}
	url := utils.StringJoin(urls...)
	remote := newRemote(url, remoteNodeUrl)
	hashMap := NewMap(defaultReplicas, nil)
	hashMap.Add(remoteNodeUrl...)
	group := NewGroup(defaultSize, "default", callback)
	h := NewHttp(group, remote, hashMap, url, basePath, api)

	return &Node{
		Map:     hashMap,
		Group:   group,
		Http:    h,
		Remote:  remote,
		baseUrl: url,
	}

}

func (n *Node) Run() {
	err := http.ListenAndServe(n.baseUrl, n)
	if err != nil {
		return
	}
	utils.OkWithMsg(http.StatusOK, utils.StringJoin("【Server】listen", n.baseUrl))
}
