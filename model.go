package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT * FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, quantity, price FROM Products WHERE id=%v", p.ID)
	row := db.QueryRow(query)
	if err := row.Scan(&p.Name, &p.Quantity, &p.Price); err != nil {
		return err
	}
	return nil
}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO Products(name, quantity, price) VALUES('%v', %v, %v)", p.Name, p.Quantity, p.Price)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// id, err := res.LastInsertId() // not supported by this postgres driver
	// if err != nil {
	// 	return err
	// }
	// p.ID = int(id)
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE Products SET name='%v', quantity=%v, price=%v WHERE id=%v", p.Name, p.Quantity, p.Price, p.ID)
	res, err := db.Exec(query)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("No such row exists")
	}

	return nil
}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM Products WHERE id=%v", p.ID)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
