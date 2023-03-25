package repository

import (
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func GetOfferID(Db *sqlx.DB, vkID int) (offerID int, err error) {
	err = Db.QueryRow("SELECT offer_id from offers WHERE executor_vk_id =$1 ORDER BY offer_id DESC LIMIT 1", vkID).Scan(&offerID)
	if err != nil {
		log.Error(err, "couldn't get the offer ID from the database")
	}
	return offerID, err
}

func GetOffer(Db *sqlx.DB, OfferID int) (offer ds.Offer, err error) {
	err = Db.QueryRow("SELECT * from offers WHERE offer_id =$1 ORDER BY offer_id DESC LIMIT 1", OfferID).Scan(&offer.OfferID, &offer.ExecutorVKID, &offer.OrderID, &offer.Price)
	if err != nil {
		log.Error(err, "couldn't get an offer from the database")
	}
	return offer, err
}
func WriteOffer(Db *sqlx.DB, orderNumber int, vkID int) error {
	var err error
	_, err = Db.Exec("INSERT INTO offers (executor_vk_id, order_id) VALUES ($1, $2)", vkID, orderNumber)
	if err != nil {
		log.WithError(err).Error("offer doesn't write")
		return err
	}
	return nil
}

func WritePriceOffer(Db *sqlx.DB, vkID int, price int) error {
	offerID, err := GetOfferID(Db, vkID)
	if err != nil {
		log.WithError(err).Error("couldn't find an offer")
		return err
	}
	_, err = Db.Exec("UPDATE offers SET price=$1 WHERE offer_id=$2", price, offerID)
	if err != nil {
		log.WithError(err).Error("offer doesn't write")
		return err
	}
	return nil
}

func DeleteOffer(Db *sqlx.DB, vkID int) error {
	offerID, err := GetOfferID(Db, vkID)
	if err != nil {
		log.WithError(err).Error("couldn't find an offer")
		return err
	}
	_, err = Db.Exec("DELETE FROM offers WHERE offer_id='$1'", offerID)
	if err != nil {
		log.WithError(err).Error("offer doesn't write")
		return err
	}
	return nil
}
