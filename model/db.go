package model

// Custom model
type MyModel struct {
	ID uint `gorm:"primary_key"`
}

type MyUserId struct {
	UID string `gorm:"not null default 0 ;type:varchar(20)" json:"uid"`
}

// Models represents all models..
var Models = []interface{}{
	&User{}, &UserOauth{}, &UserAccounts{}, &UserMessage{}, &WorkOrder{}, &UserLoginSets{}, &MinerAccounts{},
	&PayTx{}, &UserAssetflow{}, &VirtualRecharge{},
	&MgdTx{}, &WalletRecord{}, &WalletAddress{}, &WalletSync{}, &EthTx{}, &BtcTx{},
	&WalletBasic{}, &WalletSimple{},
	&SystemBasic{}, &SystemSimple{},
	&CoinMarket{}, &CoinPrice{},
	&Articles{}, &ArticleInfo{}, &ArticleCategory{},
	&Products{}, &ProductDetails{},
	&Orders{}, &FarmServer{}, &MinerOrder{}, &MinerPriceList{},
	&Assets{}, &SystemAccount{},
	&Settings{},
	&WebHook{},
}
