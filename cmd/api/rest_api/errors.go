package restapi

import "net/http"

func (app *AppServer) logError(r *http.Request, err error) {
	app.Slogger.Error(err.Error(), "request_method", r.Method, "request_url", r.URL.String())
}

func (app *AppServer) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *AppServer) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
