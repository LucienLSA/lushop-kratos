package biz

//go:generate mockgen -destination=../mocks/goods_repo.go -package=mocks . GoodsRepo
//go:generate mockgen -destination=../mocks/brand_repo.go -package=mocks . BrandRepo
//go:generate mockgen -destination=../mocks/category_repo.go -package=mocks . CategoryRepo
//go:generate mockgen -destination=../mocks/transaction.go -package=mocks . Transaction
//go:generate mockgen -destination=../mocks/es_goods_repo.go -package=mocks . EsGoodsRepo
