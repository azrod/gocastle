// screens/savescreen.go

package screens

import (
	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"github.com/zwindler/gocastle/maps"
	"github.com/zwindler/gocastle/model"
	"github.com/zwindler/gocastle/utils"
)

// ShowSaveGameScreen is the main function of the save game screen.
func ShowSaveGameScreen(window fyne.Window) {
	// Remove Images from character & inventory before saving
	playerSaveData := *player
	playerSaveData.Avatar.CanvasImage.Image = nil
	for index := range playerSaveData.Inventory {
		playerSaveData.Inventory[index].CanvasImage = nil
	}

	// Remove Images from NPCs & Objects before saving
	mapSaveData := *currentMap
	for index := range mapSaveData.NPCList {
		mapSaveData.NPCList[index].Avatar.CanvasImage.Image = nil
	}
	for index := range mapSaveData.ObjectList {
		mapSaveData.ObjectList[index].CanvasImage.Image = nil
	}

	// Get the data to save
	gameData := struct {
		Player         model.CharacterStats
		CurrentMap     maps.Map
		TimeSinceBegin int
	}{
		Player:         playerSaveData,
		CurrentMap:     mapSaveData,
		TimeSinceBegin: model.TimeSinceBegin,
	}

	// Show file save dialog
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err == nil && writer != nil {
			defer writer.Close()

			// Create JSON encoder
			encoder := json.NewEncoder(writer)

			// Write game data to JSON .sav file
			if err := encoder.Encode(gameData); err != nil {
				dialog.ShowError(err, window)
			} else {
				dialog.ShowInformation("Game Saved", "Game data has been successfully saved.", window)
			}
		}
	}, window)
	// only allow .sav files
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".sav"}))
	fd.SetFileName("backup.sav")
	// TODO: don't change path for iOS, Android, Flatpak
	location, err := utils.GetBaseDirectory()
	if err != nil {
		dialog.ShowError(err, window)
	}
	fd.SetLocation(location)
	fd.Show()
}
