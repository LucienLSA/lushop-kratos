package domain

type Address struct {
	ID           int32
	UserID       int32
	Province     string
	City         string
	District     string
	Address      string
	SignerName   string
	SignerMobile string
}

type AddressListResponse struct {
	Total int64
	List  []*Address
}
