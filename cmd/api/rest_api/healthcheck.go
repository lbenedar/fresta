package restapi

import "net/http"

func (app *AppServer) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.Cfg.Env,
			"version":     app.Version,
			"build_time":  app.BuildTime,
		},
		"foundry": map[string]any{
			"status":       app.FoundryApp.Status,
			"is_available": app.FoundryApp.IsAvailable,
		},
		"gRPC": map[string]any{
			"status": app.GrpcStatus.String(),
		},
	}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
