package api

import (
	"log/slog"
	"meme3/global"
	"meme3/server/web"
	"net/http"
	"net/url"
	"strings"

	"github.com/99nil/gopkg/ctr"
	"github.com/julienschmidt/httprouter"
	"github.com/zc2638/swag"
	"github.com/zc2638/swag/option"
)

func RouterInit(router *httprouter.Router) {
	hurl, _ := url.Parse(global.Config.HostURL)
	api := swag.New(
		option.Version("v1"),
		option.BasePath(global.Config.APIPrefix+"/v1"),
		option.Host(hurl.Host),
		option.Title("Moonn API document"),
		// option.Security("petstore_auth", "read:pets"),
		// option.SecurityScheme("petstore_auth",
		//  option.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token"),
		//  option.OAuth2Scope("write:pets", "modify pets in your account"),
		//  option.OAuth2Scope("read:pets", "read your pets"),
		// ),
	)

	restAPI := swag.New(
		option.Version("v1"),
		option.BasePath(global.Config.APIPrefix+"/rest"),
		option.Host(hurl.Host),
		option.Title("Moonn API document"),
		// option.Security("petstore_auth", "read:pets"),
		// option.SecurityScheme("petstore_auth",
		//  option.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token"),
		//  option.OAuth2Scope("write:pets", "modify pets in your account"),
		//  option.OAuth2Scope("read:pets", "read your pets"),
		// ),
	)

	// init
	initUser(api)
	initToken(api)
	initTrade(api)
	initSystem(api)

	// swagger
	api.Walk(func(path string, e *swag.Endpoint) {
		h := e.Handler.(http.Handler)
		path = swag.ColonPath(path)

		// verify
		isUser := strings.HasPrefix(path, global.Config.APIPrefix+"/v1/user")
		isLogin := strings.HasSuffix(path, "login")
		isVerify := strings.HasSuffix(path, "verify")
		if isUser && !global.Config.Debug {
			if isLogin || isVerify {
				routerHook(router, e.Method, path, h.ServeHTTP)
			} else {
				routerHook(router, e.Method, path, h.ServeHTTP, UserVerify)
			}
		} else {
			routerHook(router, e.Method, path, h.ServeHTTP)
		}
		// routerHook(router, e.Method, path, h.ServeHTTP)
	})

	// rest swagger
	restAPI.Walk(func(path string, e *swag.Endpoint) {
		h := e.Handler.(http.Handler)
		path = swag.ColonPath(path)
		routerHook(router, e.Method, path, h.ServeHTTP)
	})
	if global.Config.Debug || global.Config.Testnet {
		slog.Info("enable api swagger......")
		router.Handler(http.MethodGet, "/swagger/json", Handler(api))
	} else {
		slog.Info("enable rest swagger......")
		router.Handler(http.MethodGet, "/swagger/json", Handler(restAPI))
	}
	// router.Handler(http.MethodGet, "/swagger/ui/*any", UIHandler("/swagger/ui", "/swagger/json", true))

	// websocket
	routerHook(router, http.MethodGet, global.Config.APIPrefix+"/v1/ws", websocketHandler)
}

func routerHook(router *httprouter.Router, method string, path string, handleFunc http.HandlerFunc, hookFuncs ...func(w http.ResponseWriter, r *http.Request) bool) {
	// fmt.Printf("%s %s\n", method, path)
	router.HandlerFunc(method, path, func(w http.ResponseWriter, r *http.Request) {
		for _, hookFunc := range hookFuncs {
			if hookFunc != nil && !hookFunc(w, r) {
				WebResponseJson(w, r, ApiResponse([]byte("permission denied"), false), http.StatusForbidden)
				return
			}
		}
		handleFunc(w, r)
	})
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if err := global.WebsocketHandler(w, r); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
}

func UserVerify(w http.ResponseWriter, r *http.Request) bool {
	if userID := web.PopFromSession(r, web.C_Session_User); userID == nil {
		return false
	}
	return true
}

func Handler(api *swag.API) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// customize the swagger header based on host
		scheme := "https"
		doc := api.Clone()
		doc.Schemes = []string{scheme}
		ctr.OK(w, doc)
	}
}
