package burgers

import (
	"context"
	"errors"
	errors1 "github.com/JAbduvohidov/burger-shop.tj/pkg/crud/errors"
	"github.com/JAbduvohidov/burger-shop.tj/pkg/crud/models"
	"github.com/JAbduvohidov/burger-shop.tj/pkg/crud/services"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BurgersSvc struct {
	pool *pgxpool.Pool // dependency
}

func NewBurgersSvc(pool *pgxpool.Pool) *BurgersSvc {
	if pool == nil {
		panic(errors.New("pool can't be nil")) // <- be accurate
	}
	return &BurgersSvc{pool: pool}
}

func (service *BurgersSvc) InitDB() error {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return errors1.ApiError("can't init db: ", err)
	}
	_, err = conn.Query(context.Background(), services.BurgersDDL)
	if err != nil {
		return errors1.ApiError("can't init db: ", err)
	}
	return nil
}

func (service *BurgersSvc) BurgersList() (list []models.Burger, err error) {
	list = make([]models.Burger, 0)
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors1.ApiError("can't execute pool: ", err)
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), services.GetBurgers)
	if err != nil {
		return nil, errors1.ApiError("can't query: execute pool", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Burger{}
		err := rows.Scan(&item.Id, &item.Name, &item.Price, &item.Description)
		if err != nil {
			return nil, errors1.ApiError("can't scan row: ", err)
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (service *BurgersSvc) Save(model models.Burger) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return errors1.ApiError("can't execute pool: ", err)
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), services.SaveBurger, model.Name, model.Price, model.Description)
	if err != nil {
		return errors1.ApiError("can't save burger: ", err)
	}
	return nil
}

func (service *BurgersSvc) RemoveById(id int) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return errors1.ApiError("can't execute pool: ", err)
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), services.RemoveBurger, id)
	if err != nil {
		return errors1.ApiError("can't remove burger: ", err)
	}
	return nil
}
