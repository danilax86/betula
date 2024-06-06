package web

import (
	"git.sr.ht/~bouncepaw/betula/types"
	"net/http"
)

/* Error templates for edit link currentPage */

func (d dataEditLink) emptyURL(post types.Bookmark, data *dataCommon, w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	templateExec(w, rq, templateEditLink, dataEditLink{
		Bookmark:      post,
		dataCommon:    data,
		ErrorEmptyURL: true,
	})
}

func (d dataEditLink) invalidURL(post types.Bookmark, data *dataCommon, w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	templateExec(w, rq, templateEditLink, dataEditLink{
		Bookmark:        post,
		dataCommon:      data,
		ErrorInvalidURL: true,
	})
}

func (d dataEditLink) titleNotFound(post types.Bookmark, data *dataCommon, w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	templateExec(w, rq, templateEditLink, dataEditLink{
		Bookmark:           post,
		dataCommon:         data,
		ErrorTitleNotFound: true,
	})
}

/* Error templates for save link currentPage */

func (d dataSaveLink) emptyURL(post types.Bookmark, data *dataCommon, w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	templateExec(w, rq, templateSaveLink, dataSaveLink{
		Bookmark:      post,
		dataCommon:    data,
		ErrorEmptyURL: true,
	})
}

func (d dataSaveLink) invalidURL(post types.Bookmark, data *dataCommon, w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	templateExec(w, rq, templateSaveLink, dataSaveLink{
		Bookmark:        post,
		dataCommon:      data,
		ErrorInvalidURL: true,
	})
}

func (d dataSaveLink) titleNotFound(post types.Bookmark, data *dataCommon, w http.ResponseWriter, rq *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	templateExec(w, rq, templateSaveLink, dataSaveLink{
		Bookmark:           post,
		dataCommon:         data,
		ErrorTitleNotFound: true,
	})
}
