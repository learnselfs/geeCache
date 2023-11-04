// Package core @Author Bing
// @Date 2023/10/25 10:38:00
// @Desc
package core

import (
	"encoding/json"
	"fmt"
	"github.com/learnselfs/geeCache/utils"
	"net/http"
	"strings"
)

const (
	defaultBasePath = "/_geeCache/"
	defaultApiPath  = "/api/"
)

type Http struct {
	*Remote
	*Group
	*Map
	addr     string
	basePath string
	isApi    bool
	// request
	r *http.Request
	// response
	w http.ResponseWriter
}

func (h *Http) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.w = w
	h.r = r
	if strings.HasPrefix(r.URL.Path, h.basePath) {
		h.cacheHandle()
	} else if strings.HasPrefix(r.URL.Path, defaultApiPath) && h.isApi {
		h.apiHandle()
	} else {
		h.JSON(utils.FailedWithMsg(http.StatusNotFound, utils.StringJoin("not found", r.URL.Path)))
	}

}

func (h *Http) JSON(obj utils.H) {
	h.SetHeader("Content-Type", "application/json")
	w := json.NewEncoder(h.w)
	err := w.Encode(obj)
	if err != nil {
		h.Log("JSON: %s", obj)
	}
}

func (h *Http) Log(f string, v ...any) {
	_, err := h.w.Write([]byte(fmt.Sprintf(f, v...)))
	utils.Log(f, v...)
	if err != nil {
		return
	}
}

func (h *Http) SetHeader(k, v string) {
	h.w.Header().Set(k, v)
}

func (h *Http) cacheHandle() {
	path := strings.SplitN(h.r.URL.Path[len(h.basePath):], "/", 2)
	path = utils.ExcludeEmpty(path)
	groupName := path[0]
	key := path[1]
	g, ok := h.Group.Groups[groupName]
	if !ok {
		h.JSON(utils.Failed())
		return
	}
	v, ok := g.Get(key)
	if ok {
		h.JSON(utils.OkWithDetail(http.StatusOK, "selected successfully", utils.H{"groupName": groupName, "key": key, "value": v}))
	} else {
		h.JSON(utils.FailedWithMsg(http.StatusNotFound, utils.StringJoin("not fund", groupName, "-", key)))
	}
}

func (h *Http) pick(group, key string) ([]byte, error) {
	url := h.Map.Get(key)
	result, err := h.remoteGet(url, group, key)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h *Http) apiHandle() {
	paths := strings.SplitN(h.r.URL.Path[len(defaultApiPath):], "/", 2)
	paths = utils.ExcludeEmpty(paths)
	group := paths[0]
	key := paths[1]
	res, err := h.pick(group, key)
	if err != nil {
		h.JSON(utils.FailedWithMsg(http.StatusInternalServerError, err.Error()))
		return
	}
	var obj utils.Obj
	err = json.Unmarshal(res, &obj)
	if err != nil {
		h.JSON(utils.FailedWithMsg(http.StatusServiceUnavailable, err.Error()))
		return
	}
	if obj.Code != http.StatusOK {
		h.JSON(utils.FailedWithMsg(http.StatusServiceUnavailable, obj.Msg))
		return
	}

	data, err := json.Marshal(obj.Data)
	if err != nil {
		h.JSON(utils.FailedWithMsg(http.StatusInternalServerError, err.Error()))
		return
	}

	h.JSON(utils.OkWithMsg(http.StatusOK, string(data)))
	if err != nil {
		return
	}
}

func NewHttp(group *Group, remote *Remote, hashMap *Map, url, basePath string, api bool) *Http {
	if len(basePath) == 0 {
		basePath = defaultBasePath
	}
	return &Http{
		Remote:   remote,
		Map:      hashMap,
		Group:    group,
		addr:     url,
		basePath: basePath,
		isApi:    api,
	}
}

func Run(group *Group, remote *Remote, hashMap *Map, url, basePath string, api bool) {
	pool := NewHttp(group, remote, hashMap, url, basePath, api)
	err := http.ListenAndServe(pool.addr, pool)
	if err != nil {
		return
	}
}
