package app

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (a *Application) QbittorrentsAdd(c echo.Context, URL string) error {
	b, err := base64.StdEncoding.DecodeString(URL)
	if err != nil {
		return err
	}
	u := string(b)

	infohash, err := app.Qbt.Add(u, nil)
	if err != nil {
		return err
	}

	a.Log.Infof("added: %s", infohash)

	return c.JSON(http.StatusOK, H{"error": false, "infohash": infohash})
}

func (a *Application) QbittorrentsRemove(c echo.Context, infohash string, del bool) error {
	if infohash == "" {
		return errors.New("infohash is required")
	}

	err := app.Qbt.Delete(infohash, del)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

func (a *Application) QbittorrentsIndex(c echo.Context) error {
	// read the json string from cache
	// read the json string from cache
	res := ""
	ok, err := a.Cache.Get("flame-qbittorrents", &res)
	if err != nil {
		return err
	}
	if !ok {
		return err
	}

	return c.String(http.StatusOK, res)
}

func (a *Application) QbittorrentsLabel(c echo.Context, infohash, label string) error {
	err := app.Qbt.SetLabel([]string{infohash}, label)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

func (a *Application) QbittorrentsPause(c echo.Context, infohash string) error {
	if infohash == "" {
		a.QbittorrentsPauseAll(c)
		return nil
	}
	err := app.Qbt.Pause(infohash)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, H{"error": false})
}

func (a *Application) QbittorrentsPauseAll(c echo.Context) error {
	err := app.Qbt.PauseAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, H{"error": false})
}

func (a *Application) QbittorrentsResume(c echo.Context, infohash string) error {
	if infohash == "" {
		a.QbittorrentsResumeAll(c)
		return nil
	}

	err := app.Qbt.Resume(infohash)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

func (a *Application) QbittorrentsResumeAll(c echo.Context) error {
	err := app.Qbt.ResumeAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, H{"error": false})
}

//func (a *Application) QbittorrentsStart(c echo.Context, infohash string) error {
//	err := app.Qbt.Start(infohash)
//	if err != nil {
//		return err
//	}
//	return c.JSON(http.StatusOK, H{"error": false})
//}
//
//func (a *Application) QbittorrentsStop(c echo.Context, infohash string) error {
//	err := App().Utorrent.Stop(infohash)
//	if err != nil {
//		return err
//	}
//	return c.JSON(http.StatusOK, H{"error": false})
//}

func (a *Application) QbittorrentsWant(c echo.Context, infohash, files string) error {
	if files == "none" {
		a.QbittorrentsWantNone(c, infohash)
		return nil
	}

	ids := strings.Split(files, ",")
	err := app.Qbt.Want(infohash, ids)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

func (a *Application) QbittorrentsWantNone(c echo.Context, infohash string) error {
	err := app.Qbt.WantNone(infohash)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

func (a *Application) QbittorrentsWanted(c echo.Context, infohash string) error {
	want, err := app.Qbt.WantedCount(infohash)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false, "wanted": want > 0, "count": want})
}
