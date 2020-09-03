package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type HTTPHandlerFuncWithUser func(w http.ResponseWriter, r *http.Request, userID string) (int, error)

type APIErrorResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func writeAPIError(w http.ResponseWriter, err *APIErrorResponse) {
	b, _ := json.Marshal(err)
	w.WriteHeader(err.StatusCode)
	_, _ = w.Write(b)
}

const (
	routeAPIPrefix             = "/api/v1"
	routeAutocompleteLabels    = "/autocomplete/labels"
	routeAutocompleteBookmarks = "/autocomplete/bookmarks"
)

func (p *Plugin) initialiseAPI() {
	p.router = mux.NewRouter()
	apiRouter := p.router.PathPrefix(routeAPIPrefix).Subrouter()

	apiRouter.HandleFunc(routeAutocompleteLabels, p.extractUserMiddleWare(p.handleAutoCompleteLabels, true)).Methods("GET")
	apiRouter.HandleFunc(routeAutocompleteBookmarks, p.extractUserMiddleWare(p.handleAutoCompleteBookmarks, true)).Methods("GET")
	apiRouter.HandleFunc("/view", p.extractUserMiddleWare(p.handleViewBookmarks, true)).Methods("POST")
	apiRouter.HandleFunc("/add", p.extractUserMiddleWare(p.handleAddBookmark, true)).Methods("POST")
	apiRouter.HandleFunc("/get", p.extractUserMiddleWare(p.handleGetBookmark, true)).Methods("GET")
	apiRouter.HandleFunc("/labels/get", p.extractUserMiddleWare(p.handleLabelsGet, true)).Methods("GET")
	apiRouter.HandleFunc("/labels/add", p.extractUserMiddleWare(p.handleLabelsAdd, true)).Methods("POST")
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) extractUserMiddleWare(handler HTTPHandlerFuncWithUser, jsonResponse bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID == "" {
			if jsonResponse {
				writeAPIError(w, &APIErrorResponse{ID: "", Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
			} else {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
			}
			return
		}

		_, _ = handler(w, r, userID)
	}
}

// handleAddBookmark saves a bookmark to the bookmarks store
func (p *Plugin) handleAddBookmark(w http.ResponseWriter, r *http.Request, userID string) (int, error) {
	pluginapi := pluginapi.New(p.API)

	type bmarkWithChannel struct {
		Bookmark  *bookmarks.Bookmark `json:"bookmark"`
		ChannelID string              `json:"channelId"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return respondErr(w, http.StatusBadRequest, err)
	}

	var req *bmarkWithChannel
	if err = json.Unmarshal(body, &req); err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}
	bmark := req.Bookmark
	channelID := req.ChannelID

	bmarks, err := bookmarks.NewBookmarksWithUser(pluginapi, userID)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	l, err := bookmarks.NewLabelsWithUser(pluginapi, userID)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}
	ids := bmark.GetLabelIDs()

	var newIDs []string
	for _, id := range ids {
		label := l.ByID[id]
		// if doesn't exist, this is a name and needs to be added to the labels
		// store.  also save the id to the bookmark, not the name
		if label == nil {
			var labelNew *bookmarks.Label
			labelNew, err = l.AddLabel(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			newIDs = append(newIDs, labelNew.ID)
			continue
		}
		newIDs = append(newIDs, id)
	}

	// update bmark with UUID values, not the names
	bmark.LabelIDs = newIDs
	err = bmarks.AddBookmark(bmark)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	var names []string
	var name string
	for _, id := range newIDs {
		name, err = l.GetNameFromID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		names = append(names, name)
	}

	text, err := bmarks.GetBmarkTextOneLine(bmark, names)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}
	message := "Saved Bookmark:\n" + text

	post := &model.Post{
		UserId:    p.GetBotID(),
		ChannelId: channelID,
		Message:   message,
	}
	_ = p.API.SendEphemeralPost(userID, post)

	return http.StatusOK, nil
}

// func (p *Plugin) handleDelete(w http.ResponseWriter, r *http.Request) {
// 	return
// }

// handleViewBookmarks makes an ephemeral post listing a users bookmarks
func (p *Plugin) handleViewBookmarks(w http.ResponseWriter, r *http.Request, userID string) (int, error) {
	type requestStruct struct {
		ChannelID string `json:"channelId"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return respondErr(w, http.StatusBadRequest, err)
	}

	var req *requestStruct
	if err = json.Unmarshal(body, &req); err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}
	channelID := req.ChannelID

	pluginapi := pluginapi.New(p.API)
	bmarks, err := bookmarks.NewBookmarksWithUser(pluginapi, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	text, err := bmarks.GetBmarksEphemeralText(userID, nil)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	post := &model.Post{
		UserId:    p.GetBotID(),
		ChannelId: channelID,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(userID, post)

	return http.StatusOK, nil
}

// handleGetBookmark returns a bookmark
func (p *Plugin) handleGetBookmark(w http.ResponseWriter, r *http.Request, userID string) (int, error) {
	query := r.URL.Query()
	postID := query["postID"][0]

	pluginapi := pluginapi.New(p.API)
	bmarks, err := bookmarks.NewBookmarksWithUser(pluginapi, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// return nil if bookmark does not exist
	bmark, err := bmarks.GetBookmark(postID)
	if bmark == nil {
		var bb []byte
		_, err = w.Write(bb)
		if err != nil {
			return respondErr(w, http.StatusInternalServerError, err)
		}
		return http.StatusOK, nil
	}
	if err != nil {
		p.handleErrorWithCode(w, http.StatusBadRequest, "Unable to get bookmark", err)
	}

	resp, err := json.Marshal(bmark)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	_, err = w.Write(resp)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	return http.StatusOK, nil
}

// handleAutoCompleteLabels returns all autocomplete labels
func (p *Plugin) handleAutoCompleteLabels(w http.ResponseWriter, r *http.Request, userID string) (int, error) {
	pluginapi := pluginapi.New(p.API)
	labels, err := bookmarks.NewLabelsWithUser(pluginapi, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	out := []model.AutocompleteListItem{}
	for _, label := range labels.ByID {
		out = append(out, model.AutocompleteListItem{
			Item: label.Name,
		})
	}
	return respondJSON(w, out)
}

// handleAutoCompleteBookmarks returns all autocomplete bookmark postIDs
func (p *Plugin) handleAutoCompleteBookmarks(w http.ResponseWriter, r *http.Request, userID string) (int, error) {
	pluginapi := pluginapi.New(p.API)
	bmarks, err := bookmarks.NewBookmarksWithUser(pluginapi, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	out := []model.AutocompleteListItem{}
	for _, bmark := range bmarks.ByID {
		out = append(out, model.AutocompleteListItem{
			Item: bmark.PostID,
		})
	}
	return respondJSON(w, out)
}

func respondErr(w http.ResponseWriter, code int, err error) (int, error) {
	http.Error(w, err.Error(), code)
	return code, err
}

func respondJSON(w http.ResponseWriter, obj interface{}) (int, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, errors.WithMessage(err, "failed to marshal response"))
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return http.StatusInternalServerError, errors.WithMessage(err, "failed to write response")
	}
	return http.StatusOK, nil
}

// handleLabelsGet returns all labels
func (p *Plugin) handleLabelsGet(w http.ResponseWriter, r *http.Request, userID string) (int, error) {
	pluginapi := pluginapi.New(p.API)
	labels, err := bookmarks.NewLabelsWithUser(pluginapi, userID)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	resp, err := json.Marshal(labels)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	_, err = w.Write(resp)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}
	return http.StatusOK, nil
}

// handleLabelsAdd adds a label to the labels store
func (p *Plugin) handleLabelsAdd(w http.ResponseWriter, r *http.Request, userID string) (int, error) {
	pluginapi := pluginapi.New(p.API)
	query := r.URL.Query()
	labelName := query["labelName"][0]
	labels, err := bookmarks.NewLabelsWithUser(pluginapi, userID)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	label, err := labels.AddLabel(labelName)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	resp, err := json.Marshal(label)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}

	_, err = w.Write(resp)
	if err != nil {
		return respondErr(w, http.StatusInternalServerError, err)
	}
	return http.StatusOK, nil
}

func (p *Plugin) handleErrorWithCode(w http.ResponseWriter, code int, errTitle string, err error) {
	w.WriteHeader(code)
	b, _ := json.Marshal(struct {
		Error   string `json:"error"`
		Details string `json:"details"`
	}{
		Error:   errTitle,
		Details: err.Error(),
	})
	_, _ = w.Write(b)
}
