package dbop

import (
	"HappyShopTogether/model"
	"strconv"
)

type shopinfo struct {
	MerchantId uint   `json:"id"`
	ShopName   string `json:"shop_name"`
	ShopIntro  string `json:"shop_intro"`
	ShopAvatar string `json:"shop_avatar"`
}

func SearchCommoditiesLimitPage(limit, page string, keys []string) ([]*model.CommodityInfo, []*shopinfo) {

    var commodities []*model.CommodityInfo
    var shops []*shopinfo

	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)

	if limitInt == 0 {
		limitInt = -1
		pageInt = 1
	}

    db := model.Db.Self

    if len(keys) == 0 {

        commoditySql := "SELECT * FROM `commodity_infos`"
        merchantSql := "SELECT merchant_id, shop_name,shop_intro ,shop_avatar FROM `merchant_infos`"

        db.Raw(commoditySql).Offset(limitInt * pageInt).Limit(limitInt).Scan(&commodities)
        db.Raw(merchantSql).Offset(limitInt * pageInt).Limit(limitInt).Scan(&shops)

    } else {
        link := ""

        for _, key := range keys {
            link += key + "|"
        }

        link = "'" + link[:len(link)-1] + "'"

        commoditySql := "SELECT * FROM `commodity_infos` WHERE name REGEXP" + link
        merchantSql := "SELECT merchant_id, shop_name,shop_intro ,shop_avatar FROM `merchant_infos` WHERE shop_name REGEXP" + link

        db.Raw(commoditySql).Offset(limitInt * pageInt).Limit(limitInt).Scan(&commodities)
        db.Raw(merchantSql).Offset(limitInt * pageInt).Limit(limitInt).Scan(&shops)
    }

	return commodities, shops
}
