package app

import (
	"github.com/pkg/errors"
)

func checkDisk(freedisk int, paused bool) {
	// if less than 25GB
	if freedisk < 25000 {
		// and already paused, return
		if diskPaused() {
			return
		}

		err := pauseAll()
		if err != nil {
			app.Workers.Log.Errorf("checkdisk: failed to pause all qbts: %s", err)
		}
		return
	}

	// disk not low and paused, resume
	if paused {
		err := resumeAll()
		if err != nil {
			app.Workers.Log.Errorf("checkdisk: failed to resume all qbts: %s", err)
		}
	}
}

// pauseAll pauses all torrents and sets a flag in the cache
func pauseAll() error {
	err := app.Qbt.PauseAll()
	if err != nil {
		return errors.Wrap(err, "pausing all")
	}

	err = app.Cache.Set("flame-disk-paused", true)
	if err != nil {
		app.Workers.Log.Errorf("checkdisk: set cache failed: %s", err)
	}

	return nil
}

func resumeAll() error {
	err := app.Qbt.ResumeAll()
	if err != nil {
		app.Workers.Log.Errorf("checkdisk: failed to resume all qbts: %s", err)
	}

	err = app.Cache.Delete("flame-disk-paused")
	if err != nil {
		app.Workers.Log.Errorf("checkdisk: set cache failed: %s", err)
	}

	return nil
}

// diskPaused checks the cache for the flag
func diskPaused() bool {
	paused := "false"
	_, err := app.Cache.Get("flame-disk-paused", &paused)
	if err != nil {
		// app.Workers.Log.Errorf("paused: %s", err)
		return false
	}
	return paused == "true"
}
