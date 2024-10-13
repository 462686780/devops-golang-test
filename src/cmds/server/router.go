package server

import (
	"net/http"
	"statefulset/base"
	validateadmissionwebhook "statefulset/cmds/server/ValidateAdmissionWebhook"
	"statefulset/cmds/server/context"
)

func (server *ApiServer) InitRoute() {
	server.router.Use()
	server.AddRoute(http.MethodGet, "/ping", ping)
	server.AddRoute(http.MethodGet, "/vaildate", validateadmissionwebhook.ValidateAdmission)
}

func ping(c *context.Context) error {
	base.Context.Logger.Debugf("handler ping")

	c.WriteJSONResponse(http.StatusOK, "pong")
	return nil
}
