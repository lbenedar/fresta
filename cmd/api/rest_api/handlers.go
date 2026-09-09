package restapi

import (
	"fmt"
	"net/http"

	"github.com/lbenedar/fresta/internal/models/db"
)

func (app *AppServer) ShowText(w http.ResponseWriter, r *http.Request) {
	dbConn := app.FoundryApp.Transport.DB

	world, err := db.GetWorld(dbConn, "kingmaker")
	if err != nil {
		w.Write([]byte("Got error"))
		return
	}

	fmt.Fprintf(w, "Got world with name %s, core_version - %s, next_session - %v", world.ID, world.CoreVersion, world.NextSession)
}

func (app *AppServer) WhenNextSession(w http.ResponseWriter, r *http.Request) {
	// data, err := app.foundryApp.HandleWSRequest("/setup")
	// if err != nil {
	// 	app.slogger.Error("", "error", err)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	// setupData, err := json_model.ParseSetupModel(data)
	// if err != nil {
	// 	app.slogger.Error("", "error", err)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	// // nextSession, err := setupData.GetWorlds().GetSessionTime("")
	// // if err != nil {
	// // 	app.slogger.Error("", "error", err)
	// // 	w.Write([]byte(err.Error()))
	// // 	return
	// // }

	// app.slogger.Info("Next session data is ready to send", "sessionTime", nextSession)
	// w.Write([]byte(nextSession.Local().String()))
	return
}
