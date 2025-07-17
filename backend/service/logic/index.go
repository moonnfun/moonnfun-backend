package logic

// import (
// 	"meme3/service/model"
// 	"sync"
// )

// const (
// 	c_page_tokens_max = 100
// )

// var (
// 	v_lock_pages = new(sync.RWMutex)

// 	v_cache_pages []*model.Token
// )

// func GetCacheTokens(offset, limit int) []*model.Token {
// 	pageIndex := offset - 1

// 	bIndex0 := true
// 	workingLock := v_lock_pages
// 	cacheTokens := v_cache_list.Tokens
// 	if pageIndex/5 > 0 {
// 		bIndex0 = false
// 		cacheTokens = v_cache_pages
// 	}

// 	// prev
// 	if !bIndex0 && offset%5 == 0 {

// 	}

// 	// next
// 	if !bIndex0 && offset%5 == 1 {

// 	}
// }
