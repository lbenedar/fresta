package transport

import (
	"slices"
	"time"

	"github.com/lbenedar/fresta/internal/foundry/requests"
	"github.com/lbenedar/fresta/internal/foundry/types"
	"github.com/lbenedar/fresta/internal/models/db"
	"github.com/lbenedar/fresta/internal/models/json"
)

func (tr *FoundryTransport) FillDBWithFoundryData() error {
	now := time.Now()

	setupExist, err := db.HasSetup(tr.DB)
	if err != nil {
		return err
	}

	if !setupExist {
		err := tr.InsertSetupToDB()
		if err != nil {
			return err
		}
	}

	err = tr.InitWorldsData()
	if err != nil {
		return err
	}

	tr.Logger.Info("Initialization has been successful", "elapsed_time", time.Since(now).String())
	return nil
}

func (tr *FoundryTransport) InsertSetupToDB() error {
	msgJson, err := tr.GetJsonDataByType(requests.SetupPath)
	if err != nil {
		return err
	}

	foundryStateJson, err := json.ParseSetup(msgJson)
	if err != nil {
		return err
	}

	tr.Logger.Info("Setup data received")

	var foundryStateDb db.Setup
	foundryStateJson.ToDB(&foundryStateDb)
	err = foundryStateDb.Insert(tr.DB)
	if err != nil {
		return err
	}
	return nil
}

func (tr *FoundryTransport) InitWorldsData() error {
	worldNames, err := db.GetNotStartedWorlds(tr.DB)
	if err != nil {
		return err
	}

	for i := range tr.InitWorlds {
		initWorld := tr.InitWorlds[i]
		if !slices.Contains(worldNames, initWorld.Name) {
			tr.Logger.Warn("World name does not exist in database or have been already inserted. Skip!", "world_name", initWorld.Name)
			continue
		}

		err = tr.InitWorldData(initWorld)
		if err != nil {
			return err
		}

		_, err := tr.Http.PostReturnToSetup()
		if err != nil {
			return err
		}
	}
	return nil
}

func (tr *FoundryTransport) InitWorldData(worldData types.WorldData) error {
	tr.Logger.Debug("Launching world...", "world_name", worldData.Name)
	err := tr.LaunchWorld(worldData.Name)
	if err != nil {
		return err
	}

	if worldData.UserId == "" {
		worldData.UserId, err = tr.GetUserId(worldData.Username)
		if err != nil {
			return err
		}
	}

	err = tr.ConnectToWorld(worldData.UserId, worldData.UserPass)
	if err != nil {
		return err
	}
	return tr.InsertGameToDB()
}

func (tr *FoundryTransport) InsertGameToDB() error {
	msgJson, err := tr.GetJsonData("world")
	if err != nil {
		return err
	}

	game, err := json.ParseGame(msgJson)
	if err != nil {
		return err
	}

	var gameDb db.Game
	game.ToDB(&gameDb)

	isInserted, err := db.IsWorldInserted(tr.DB, game.World.ID)
	if err != nil {
		return err
	}
	if isInserted {
		return nil
	}

	tr.Logger.Info("Game data received", "world_name", game.World.ID)

	err = gameDb.Insert(tr.DB)
	if err != nil {
		return err
	}

	tr.Logger.Info("Game data inserted to DB", "world_name", game.World.ID)
	return nil
}

func (tr *FoundryTransport) GetSetup() (*db.Setup, error) {
	var setup db.Setup

	err := setup.Get(tr.DB)
	if err != nil {
		return nil, err
	}

	return &setup, err
}
