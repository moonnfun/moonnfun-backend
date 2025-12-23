package logic

import (
	"errors"
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/service/model"
	"meme3/service/monitor"
	"meme3/service/store"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	C_listing_max = 6

	c_entry_image_path = "/images/banner/entry.jpg"

	c_listing_timeout = time.Duration(72) * time.Hour

	c_listing_timeout_debug = time.Duration(30) * time.Minute
)

var (
	v_cache_listing sync.Map

	v_cache_waitting sync.Map

	v_entry_listing *model.Listing

	v_system_listing []*model.Listing

	V_listing_timeout time.Duration

	V_listing_waitting = store.NewCache()
)

func InitSystemDB() error {
	if store.DBCount(model.C_System, bson.M{}) == 0 {
		sys := &model.System{
			Hot:    make([]string, 0),
			Banner: make([]string, 0),
		}
		sys.DBID = primitive.NewObjectID()
		if store.DBCount(model.C_System, bson.M{}) == 0 {
			store.DBSet(model.C_System, sys, bson.M{})
		}
	}

	v_system_listing = make([]*model.Listing, 0)
	sys, err := store.DBGet[model.System](model.C_System, bson.M{})
	if err != nil {
		return err
	}
	for _, bannerUrl := range sys.Banner {
		v_system_listing = append(v_system_listing, &model.Listing{
			BaseModel:      model.BaseModel{DBID: primitive.NewObjectID()},
			System:         true,
			Entry:          false,
			Listed:         true,
			BannerImageUrl: bannerUrl,
		})
	}
	return nil
}

func InitListing() error {
	V_listing_timeout = c_listing_timeout
	if global.Config.Debug {
		V_listing_timeout = c_listing_timeout_debug
	}
	v_entry_listing = &model.Listing{
		BaseModel:      model.BaseModel{DBID: primitive.NewObjectID()},
		System:         true,
		Entry:          true,
		Listed:         false,
		Start:          -1,
		BannerImageUrl: fmt.Sprintf("%s%s", global.Config.HostURL, c_entry_image_path),
	}

	listingTotal := 0
	listings, _, _ := store.DBListPage[model.Listing](model.C_Listing, bson.M{"status": model.C_Status_listed}, "desc", "createdAt", 0, 0, false)
	if listings != nil {
		for _, l := range listings {
			if time.Since(time.UnixMilli(l.Start)) < V_listing_timeout {
				l.Entry = false
				l.Listed = true
				l.System = false
				PushListing(l, true)
				listingTotal += 1
			}
		}
	}
	if listingTotal == 0 {
		return PushSysListing(true)
	}
	return nil
}

func listingKey(l *model.Listing) string {
	return l.DBID.String()
}

func PushSysListing(bInit bool) error {
	if bInit {
		AddListing(v_entry_listing, "add entry listing")
	} else {
		userCount := 0
		totalCount := 0
		v_cache_listing.Range(func(key, value any) bool {
			totalCount += 1
			if !value.(*model.Listing).System && !value.(*model.Listing).Entry {
				userCount += 1
				return true
			}
			return true
		})
		if userCount < C_listing_max {
			AddListing(v_entry_listing, "add entry listing")
		}
		if userCount > 0 {
			return nil
		}
	}

	for i := 0; i < len(v_system_listing); i++ {
		AddListing(v_system_listing[i], "add system listing")
	}
	return nil
}

func ClearSysListing() error {
	for i := 0; i < len(v_system_listing); i++ {
		RemoveListing(v_system_listing[i], "remove system listing")
	}
	RemoveListing(v_entry_listing, "remove entry listing")
	return nil
}

func PushListing(l *model.Listing, bWatch bool) error {
	count := 0
	var existListing *model.Listing
	v_cache_listing.Range(func(key, value any) bool {
		if fmt.Sprintf("%v", key) == listingKey(l) {
			existListing = value.(*model.Listing)
			return false
		}
		if !value.(*model.Listing).System {
			count += 1
		}
		return true
	})
	if existListing != nil {
		return errors.New("push listing faild with already exists")
	}
	if count >= C_listing_max {
		return errors.New("push listing faild with reaching the limit")
	}
	ClearSysListing()

	count += 1
	AddListing(l, "add user listing")
	if bWatch {
		go listingWithTimeout(l)
	}
	if count < C_listing_max {
		AddListing(v_entry_listing, "add entry listing")
	}
	return nil
}

func listingWithTimeout(l *model.Listing) {
	slog.Info("start watch listing", "listing", l, "time", time.Now().Format(time.DateTime))
	for {
		select {
		case <-time.After(V_listing_timeout - time.Since(time.UnixMilli(l.Start))):
			slog.Info("end watch listing with timeout successed", "listing", l, "time", time.Now().Format(time.DateTime))
			RemoveListing(l, "watch timeout")
			PushSysListing(false)
			return
		}
	}
}

func PushListingWait(walletAddress, tokenAddress string) {
	slog.Info("before add listing wait", "walletAddress", walletAddress, "tokenAddress", tokenAddress)
	v_cache_waitting.Store(fmt.Sprintf("%s_%s", walletAddress, tokenAddress), byte('0'))
}

func RemoveListingWait(walletAddress, tokenAddress, reason string) {
	slog.Info("before remove listing wait", "walletAddress", walletAddress, "tokenAddress", tokenAddress, "reason", reason)
	v_cache_waitting.Delete(fmt.Sprintf("%s_%s", walletAddress, tokenAddress))
}

func AddListing(l *model.Listing, reason string) {
	slog.Info("before add listing", "listing", l, "reason", reason)
	v_cache_listing.Store(listingKey(l), l)
}

func RemoveListing(l *model.Listing, reason string) {
	slog.Info("before remove listing", "listing", l, "reason", reason)
	v_cache_listing.Delete(listingKey(l))
}

func EnableListing(l *model.Listing) {
	// 直接上新
	l.Listed = true
	l.Start = time.Now().UnixMilli()
	PushListing(l, true)
	slog.Info("enable listing display successed", "listing", l)
}

func GetListingWait() *model.ListingWait {
	var total, waitTotal int64
	v_cache_listing.Range(func(key, value any) bool {
		if !value.(*model.Listing).System {
			total += 1
		}
		return true
	})
	v_cache_waitting.Range(func(key, value any) bool {
		waitTotal += 1
		return true
	})
	return &model.ListingWait{
		Total: total + waitTotal,
	}
}

func GetCacheListing() []*model.Listing {
	ret := make([]*model.Listing, 0)
	v_cache_listing.Range(func(key, value any) bool {
		if l, ok := value.(*model.Listing); ok && l != nil && l.Listed {
			ret = append(ret, value.(*model.Listing))
		}
		return true
	})

	clist := model.ListingList(ret)
	sort.Sort(clist)
	return clist[:]
}

func SaveListing(l *model.Listing) error {
	l.Token = common.HexToAddress(l.Token).String()
	// l.Wallet = common.HexToAddress(l.Wallet).String()

	var createdAt int64
	var dbID primitive.ObjectID
	if dbListing, err := store.DBGet[model.Listing](model.C_Listing, bson.M{"_id": l.DBID.String()}); err == nil && dbListing != nil {
		dbID = dbListing.DBID
		createdAt = dbListing.CreatedAt
	} else {
		dbID = l.DBID
		createdAt = time.Now().UnixMilli()
	}

	if createdAt != 0 {
		l.CreatedAt = createdAt
	}
	if err := store.DBSet(model.C_Listing, l, bson.M{"_id": dbID}); err != nil {
		slog.Error("update listing failed", "listing", l, "error", err.Error())
		return err
	}
	slog.Info("update listing successed", "listing", l)
	return nil
}

func CheckListing(l *model.Listing) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if monitor.IsListToken(l.ID, l.Token, l.Wallet) == "true" {
				RemoveListingWait(l.Wallet, l.Token, "payment successed")
				GetListingIDFromCache(l.Wallet, true)
				l.Status = model.C_Status_listed
				EnableListing(l)
				SaveListing(l)
				return
			}
		case <-time.After(time.Duration(5) * time.Minute):
			RemoveListingWait(l.Wallet, l.Token, "payment failed")
			l.Status = model.C_Status_verify_failed
			l.VerifyResult = "payment failed"
			SaveListing(l)
			ticker.Stop()
			return
		}
	}
}

func GetListingIDFromCache(userID any, bRemove bool) (string, error) {
	if listing := store.CacheGet(fmt.Sprintf("wait_listing_%v", userID), false, nil); listing != nil {
		if bRemove {
			if err := SaveImageForSub("banner", listing.(*model.Listing).ImageFileName, listing.(*model.Listing).ImageFileBuf); err != nil {
				slog.Error("save listing image failed", "listing", listing, "error", err.Error())
				return "", err
			}
			store.CacheGet(fmt.Sprintf("wait_listing_%v", userID), true, nil)
		}
		return listing.(*model.Listing).ID, nil
	} else {
		return "", errors.New("can not find listing ID")
	}
}
