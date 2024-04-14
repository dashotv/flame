package app

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/dashotv/flame/nzbget"
)

func (a *Application) NzbsAdd(c echo.Context, URL, cat, name string) error {
	pri, err := QueryDefaultInteger(c, "priority", nzbget.PriorityNormal)
	if err != nil {
		return err
	}

	b, err := base64.StdEncoding.DecodeString(URL)
	if err != nil {
		return err
	}

	u := strings.Replace(string(b), "&amp;", "&", -1)

	options := nzbget.NewOptions()
	options.Category = cat
	options.Priority = pri
	if name != "" {
		options.NiceName = name
	}

	id, err := a.Nzb.Add(u, options)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false, Result: id})
}

func (a *Application) NzbsRemove(c echo.Context, id int) error {
	var err error

	if c.QueryParam("delete") == "true" {
		err = a.Nzb.Delete(id)
	} else {
		err = a.Nzb.Remove(id)
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}

func (a *Application) NzbsDestroy(c echo.Context, id int) error {
	err := a.Nzb.Destroy(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false})
}

func (a *Application) NzbsIndex(c echo.Context) error {
	// read the json string from cache
	res := &nzbget.GroupResponse{}
	ok, err := a.Cache.Get("flame-nzbs", &res)
	if err != nil {
		return err
	}
	if !ok {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false, Result: res})
}

func (a *Application) NzbsHistory(c echo.Context, hidden bool) error {
	r, err := a.Nzb.History(hidden)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nzbget.HistoryResponse{Response: &nzbget.Response{Timestamp: time.Now()}, Result: r})
}

func (a *Application) NzbsPause(c echo.Context, id int) error {
	var err error

	if id == -1 {
		a.NzbsPauseAll(c)
		return nil
	}

	err = a.Nzb.Pause(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Response{Error: false})
}

func (a *Application) NzbsPauseAll(c echo.Context) error {
	err := a.Nzb.PauseAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}

func (a *Application) NzbsResume(c echo.Context, id int) error {
	var err error

	if id == -1 {
		a.NzbsResumeAll(c)
		return nil
	}
	err = a.Nzb.Resume(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}

func (a *Application) NzbsResumeAll(c echo.Context) error {
	err := a.Nzb.ResumeAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &Response{Error: false})
}
