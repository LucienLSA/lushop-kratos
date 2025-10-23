package domain

type Favorite struct {
	UserID  int32
	GoodsID int32
}

type UserFavResponse struct {
	UserID  int32
	GoodsID int32
}

type UserFavListResponse struct {
	Total int64
	List  []*UserFavResponse
}
