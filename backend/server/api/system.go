package api

import (
	"fmt"
	"io"
	"log/slog"
	"meme3/global"
	"meme3/service/logic"
	"meme3/service/model"
	"meme3/service/monitor"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/types"
)

func initSystem(api *swag.API) {
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/system/banner",
			endpoint.Tags("System"),
			endpoint.Handler(webSystemBanner),
			endpoint.Summary("system banner"),
			endpoint.Description("get system banner"),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]string{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/system/listing/wait",
			endpoint.Tags("System"),
			endpoint.Handler(webSystemListingWait),
			endpoint.Summary("listing wait"),
			endpoint.Description("get listing wait"),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption(model.ListingWait{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/system/listing",
			endpoint.Tags("System"),
			endpoint.Handler(WebSystemListing),
			endpoint.Summary("System listing"),
			endpoint.Description("System listing"),
			endpoint.Body(model.Listing{}, "Help object that needs to be added to the store", true),
			endpoint.Response(http.StatusOK, "Successfully add user", endpoint.SchemaResponseOption(model.Listing{})),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/system/image/upload",
			endpoint.Tags("System"),
			endpoint.Handler(webSystemImageUpload),
			endpoint.Summary("Banner image upload"),
			endpoint.Description("Banner image upload"),
			endpoint.FormData("file", types.File, "upload banner image", true),
			endpoint.Response(http.StatusOK, "Successfully added help", endpoint.SchemaResponseOption("imageUrl")),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
}

func webSystemBanner(w http.ResponseWriter, r *http.Request) {
	WebResponseJson(w, r, ApiResponse(logic.GetCacheListing(), true), http.StatusOK)
}

func webSystemListingWait(w http.ResponseWriter, r *http.Request) {
	// userID := web.PopFromSession(r, web.C_Session_User)
	// if userID == nil {
	// 	WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
	// 	return
	// }
	WebResponseJson(w, r, ApiResponse(logic.GetListingWait(), true), http.StatusOK)
}

func WebSystemListing(w http.ResponseWriter, r *http.Request) {
	// userID := web.PopFromSession(r, web.C_Session_User)
	// if userID == nil {
	// 	WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
	// 	return
	// }

	// user, err := logic.GetUser(fmt.Sprintf("%v", userID), false)
	// if err != nil {
	// 	WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
	// 	return
	// }

	_, listing, err := WebBody[model.Listing](r)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	// listing.Wallet = user.Address

	// push, 先霸占位置, 如果超时未获取到已支付进行清除
	if err := logic.PushListing(listing, true); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	// check
	listing.Status = model.C_Status_pending
	if monitor.IsListToken(listing.Token, listing.Wallet) == "true" {
		listing.Status = model.C_Status_paid
	}

	// save to db
	if err := logic.SaveListing(listing); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	// check with timeout
	if listing.Status != model.C_Status_paid {
		go logic.CheckListing(listing)
	} else {
		// 直接上新
		listing.Listed = true
		slog.Info("enable listing display successed", "listing", listing)
	}
	WebResponseJson(w, r, ApiResponse(listing, true), http.StatusOK)
}

func webSystemImageUpload(w http.ResponseWriter, r *http.Request) {
	// userID := web.PopFromSession(r, web.C_Session_User)
	// if userID == nil {
	// 	WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
	// 	return
	// }

	directory, err := os.Getwd()
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	directory = filepath.Join(directory, "dist", "images", "banner")
	if _, err := os.Stat(directory); err != nil {
		if derr := os.Mkdir(directory, os.ModePerm); derr != nil {
			WebResponseJson(w, r, ApiError(derr.Error()), http.StatusInternalServerError)
			return
		}
	}

	r.ParseMultipartForm(32 << 20)
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileBuf, err := io.ReadAll(file)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if mimeType == "" {
		mimeType = http.DetectContentType(fileBuf)
	}
	if !slices.Contains(imageTypes, mimeType) {
		WebResponseJson(w, r, ApiError("invalid image"), http.StatusInternalServerError)
		return
	}
	fileName := fmt.Sprintf("%v%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))

	fullPath := filepath.Join(directory, fileName)
	f, err := os.Create(fullPath)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err = f.Write(fileBuf); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	WebResponseJson(w, r, ApiResponse(fmt.Sprintf("%s/images/banner/%s", global.Config.HostURL, fileName), true), http.StatusOK)
}
