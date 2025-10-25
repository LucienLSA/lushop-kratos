package biz

//go:generate mockgen -destination=../mocks/address_repo.go -package=mocks userop/internal/biz AddressRepo
//go:generate mockgen -destination=../mocks/favorite_repo.go -package=mocks userop/internal/biz FavoriteRepo
