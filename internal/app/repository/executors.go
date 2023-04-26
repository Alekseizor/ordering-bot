package repository

import (
	"database/sql"
	"fmt"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func IsExecutor(Db *sqlx.DB, vkID int) (bool, error) {
	var exec ds.Executor
	err := Db.QueryRow("SELECT 1 from executors WHERE vk_id = $1", vkID).Scan(&exec.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No executor with vk_id founded")
			return false, err
		} else {
			log.Println(err)
			log.Println("Query error")
			return false, err
		}
	}
	return true, nil
}

func IsExecutorInOrder(Db *sqlx.DB, orderID, vkID int) (bool, error) {
	var exec ds.Executor
	err := Db.QueryRow("SELECT 1 from executors WHERE vk_id = $1", vkID).Scan(&exec.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No executor with vk_id founded")
			return false, err
		} else {
			log.Println(err)
			log.Println("Query error")
			return false, err
		}
	}
	var ord ds.Order
	err = Db.QueryRow("SELECT 1 from orders WHERE id = $1 AND executor_vk_id = $2", orderID, vkID).Scan(&ord.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No executor with vk_id in this order founded")
			return false, err
		} else {
			log.Println(err)
			log.Println("Query error")
			return false, err
		}
	}
	return true, nil
}

func DeleteExecutor(Db *sqlx.DB, vkID int) error {
	log.Println(vkID)
	_, err := Db.Exec("DELETE from executors WHERE vk_id = $1", vkID)
	if err != nil {
		log.Println("Executor with vk_id unknown. Delete failed")
		return err
	}
	return nil
}

func ExecutorsDiscipline(Db *sqlx.DB, disciplineID int) ([]int, error) {
	var executorsVKID []int
	query := fmt.Sprintf("SELECT vk_id FROM executors WHERE disciplines_id@> ARRAY[%d]", disciplineID)
	err := Db.Select(&executorsVKID, query)
	log.Println(executorsVKID)
	if err != nil {
		log.WithError(err).Error("failed to request a executor by discipline")
		return nil, err
	}
	return executorsVKID, nil
}

func GetExecutor(Db *sqlx.DB, vkID int) (ds.Executor, error) {
	exec := ds.Executor{}
	var DisciplinesID []sql.NullInt64
	err := Db.QueryRow("SELECT * from executors WHERE vk_id = $1", vkID).Scan(&exec.Id, &exec.VkID, pq.Array(&DisciplinesID), &exec.PercentExecutor, &exec.Rating, &exec.AmountRating, &exec.Profit, &exec.AmountOrders, &exec.Requisite)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No executor with vk_id founded")
		} else {
			log.Println(err)
			log.Println("Query error")
		}
	}
	var IDs []int
	for _, val := range DisciplinesID {
		IDs = append(IDs, int(val.Int64))
	}
	exec.DisciplinesID = IDs
	return exec, err
}

func GetExecutorByID(Db *sqlx.DB, ID int) (ds.Executor, error) {
	exec := ds.Executor{}
	var DisciplinesID []sql.NullInt64
	err := Db.QueryRow("SELECT * from executors WHERE id = $1", ID).Scan(&exec.Id, &exec.VkID, pq.Array(&DisciplinesID), &exec.PercentExecutor, &exec.Rating, &exec.AmountRating, &exec.Profit, &exec.AmountOrders, &exec.Requisite)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No executor with vk_id founded")
		} else {
			log.Println(err)
			log.Println("Query error")
		}
	}
	var IDs []int
	for _, val := range DisciplinesID {
		IDs = append(IDs, int(val.Int64))
	}
	exec.DisciplinesID = IDs
	return exec, err
}

//func AddUniqueDiscipline(num string, disciplines *map[string]bool, unique *[]string) {
//	if disciplines[num] {
//		return // Already in the map
//	}
//	unique = append(unique, num)
//	disciplines[num] = true
//}
