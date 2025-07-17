package api

import (
	"fmt"
	"io"
	"log/slog"
	"meme3/global"
	"meme3/server/web"
	"meme3/service/logic"
	"meme3/service/model"
	"meme3/service/store"
	"mime"
	"net/http"
	"path/filepath"
	"slices"
	"time"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BannerImage struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

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
			endpoint.Query("tokenAddress", "string", "token address", false),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption(model.ListingWait{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/system/listing/prepare",
			endpoint.Tags("System"),
			endpoint.Handler(WebSystemListingPrepare),
			endpoint.Summary("System listing prepare"),
			endpoint.Description("System listing prepare"),
			endpoint.Query("cancel", "string", "cancel or not", false),
			endpoint.Query("tokenAddress", "string", "token address", false),
			endpoint.Response(http.StatusOK, "Successfully add user", endpoint.SchemaResponseOption(model.ListingWait{})),
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
			endpoint.Response(http.StatusOK, "Successfully added help", endpoint.SchemaResponseOption(BannerImage{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/system/orders",
			endpoint.Tags("System"),
			endpoint.Handler(webSystemOrders),
			endpoint.Summary("system order list"),
			endpoint.Description("get system order list"),
			endpoint.Query("kind", "string", "order kind", false),
			endpoint.Query("address", "string", "wallet address", false),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]*model.Listing{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
}

func webSystemBanner(w http.ResponseWriter, r *http.Request) {
	WebResponseJson(w, r, ApiResponse(logic.GetCacheListing(), true), http.StatusOK)
}

func webSystemOrders(w http.ResponseWriter, r *http.Request) {
	kind := WebParams(r).Get("kind")
	if kind == "" {
		kind = "listing"
	}
	dbName := kind

	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}

	timeNow := time.Now().UnixMilli()
	orderList, _, err := logic.GetModelListPageEx[model.Listing](dbName, bson.M{"wallet": userID}, "", "createdAt", 0, 0)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	for i, ol := range orderList {
		if orderList[i].Status == model.C_Status_listed {
			orderList[i].Remain = time.UnixMilli(ol.Start).Add(logic.V_listing_timeout).UnixMilli() - timeNow
			if orderList[i].Remain < 0 {
				orderList[i].Status = model.C_Status_removed
				orderList[i].Remain = 0
			}
		}
	}
	WebResponseJson(w, r, ApiResponse(orderList, true), http.StatusOK)
}

func webSystemListingWait(w http.ResponseWriter, r *http.Request) {
	WebResponseJson(w, r, ApiResponse(logic.GetListingWait(), true), http.StatusOK)
}

func WebSystemListingPrepare(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}
	cancel := WebParams(r).Get("cancel")

	l := logic.GetListingWait()
	if l != nil && l.Total >= logic.C_listing_max && cancel != "true" {
		WebResponseJson(w, r, ApiError("push listing faild with reaching the limit"), http.StatusInternalServerError)
		return
	}

	// check CA
	tokenAddress := WebParams(r).Get("tokenAddress")
	if tokenAddress != "" {
		if _, err := logic.GetToken(tokenAddress); err != nil {
			WebResponseJson(w, r, ApiError("invalid CA"), http.StatusNotFound)
			return
		}
	}

	if WebParams(r).Get("cancel") == "true" {
		logic.RemoveListingWait(fmt.Sprintf("%v", userID), tokenAddress, "user cancel")
		l = logic.GetListingWait()
	} else {
		logic.PushListingWait(fmt.Sprintf("%v", userID), tokenAddress)
	}
	WebResponseJson(w, r, ApiResponse(l, true), http.StatusOK)
}

func WebSystemListing(w http.ResponseWriter, r *http.Request) {
	_, listing, err := WebBody[model.Listing](r)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	listing.Entry = false
	listing.Listed = false
	listing.System = false

	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = listing.Wallet
		}
	}
	listing.Wallet = fmt.Sprintf("%v", userID)

	// check txhash
	if userListing, err := store.DBGet[model.Listing](model.C_Trade, bson.M{"txhash": listing.TxHash}); err == nil && userListing != nil {
		WebResponseJson(w, r, ApiError("permission denied: invalid txhash already exists"), http.StatusForbidden)
		return
	}
	// if !monitor.IsValidTx(listing.Wallet, listing.TxHash) {
	// 	WebResponseJson(w, r, ApiError("permission denied, invalid txhash"), http.StatusForbidden)
	// 	return
	// }

	// check
	listing.Status = model.C_Status_verifying

	// save to db
	listing.DBID = primitive.NewObjectID()
	listing.ID, err = logic.GetListingIDFromCache(userID, false)
	if err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}
	if err := logic.SaveListing(listing); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
		return
	}

	go logic.CheckListing(listing)

	WebResponseJson(w, r, ApiResponse(listing, true), http.StatusOK)
}

func webSystemImageUpload(w http.ResponseWriter, r *http.Request) {
	userID := web.PopFromSession(r, web.C_Session_User)
	if userID == nil {
		if !global.Config.Debug {
			WebResponseJson(w, r, ApiError("permission denied"), http.StatusForbidden)
			return
		} else {
			userID = WebParams(r).Get("address")
		}
	}

	// directory, err := os.Getwd()
	// if err != nil {
	// 	WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
	// 	return
	// }
	// directory = filepath.Join(directory, "dist", "images", "banner")
	// if _, err := os.Stat(directory); err != nil {
	// 	if derr := os.Mkdir(directory, os.ModePerm); derr != nil {
	// 		WebResponseJson(w, r, ApiError(derr.Error()), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

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

	logic.GetListingIDFromCache(userID, true)
	listingID := fmt.Sprintf("%v", time.Now().UnixNano())
	fileName := fmt.Sprintf("%v%s", listingID, filepath.Ext(fileHeader.Filename))
	slog.Info("get image successed", "mimeType", mimeType, "fileName", fileName)

	l := &model.Listing{
		ID:            listingID,
		ImageFileBuf:  fileBuf,
		ImageFileName: fileName,
	}
	store.CacheSetByTime(fmt.Sprintf("wait_listing_%v", userID), l, true, time.Duration(5)*time.Minute, nil)

	ret := struct {
		ID  string `json:"id"`
		Url string `json:"url"`
	}{
		ID:  listingID,
		Url: fmt.Sprintf("%s/images/banner/%s", global.Config.HostURL, fileName),
	}
	WebResponseJson(w, r, ApiResponse(ret, true), http.StatusOK)

	// fullPath := filepath.Join(directory, fileName)
	// f, err := os.Create(fullPath)
	// if err != nil {
	// 	WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
	// 	return
	// }
	// defer f.Close()

	// if _, err = f.Write(fileBuf); err != nil {
	// 	WebResponseJson(w, r, ApiError(err.Error()), http.StatusInternalServerError)
	// 	return
	// }
	// WebResponseJson(w, r, ApiResponse(fmt.Sprintf("%s/images/banner/%s", global.Config.HostURL, fileName), true), http.StatusOK)
}
