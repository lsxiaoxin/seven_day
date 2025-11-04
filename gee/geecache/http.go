package geecache

import (
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	parts :=  strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2) //按照第一个"/"切分成俩部分
	if len(parts) != 2 {
		http.Error(w, "bad reuuest", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := groups[groupName]
	if group == nil {
		http.Error(w, "no such group: "+ groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())

}