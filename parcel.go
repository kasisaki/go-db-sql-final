package main

import (
	"database/sql"
	"errors"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуем добавление строки в таблицу parcel, используя данные из переменной p (Parcel)
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :date)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		log.Print("Error adding a new parcel into DB")
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Print("Error getting ID of a last insert while adding a new parcel into DB")
		return 0, err
	}

	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуем чтение строки по заданному number
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	// здесь из таблицы должна вернуться только одна строка, т.к. number является первичным ключом в таблице parcel

	// заполняем объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		log.Println("Error scanning data from a row while getting a parcel by number")
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// чтение строк из таблицы parcel по заданному client
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		log.Println("Error getting parcels by client")
		return []Parcel{}, err
	}

	// заполняем срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		// Если возникает ошибка при получении данных из запроса, то возвращаем список, при этом не новый
		// чтобы часть данных, которая успешно добавлена ранее могла быть выдана по результату запроса
		if err != nil {
			log.Println("Error scanning data from a row while getting a parcel by client")
			return res, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		log.Println("Error updating status of a parcel")
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	isRegistered, err := checkParcelIsRegistered(number, s.db)
	if err != nil {
		return err
	}

	if isRegistered {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			log.Println("Error updating address of a parcel")
			return err
		}
	} else {
		log.Println("Address can be changed only for registered parcels")
		return errors.New("parcel is not registered")
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// удалять строку можно только если значение статуса registered
	isRegistered, err := checkParcelIsRegistered(number, s.db)
	if err != nil {
		return err
	}

	if isRegistered {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number",
			sql.Named("number", number))
		if err != nil {
			log.Println("Error deleting a parcel")
			return err
		}
	} else {
		log.Println("Only registered parcels can be deleted")
		return errors.New("parcel is not registered")
	}

	return nil
}

func checkParcelIsRegistered(number int, db *sql.DB) (bool, error) {
	row := db.QueryRow("SELECT status FROM parcel WHERE number = :number",
		sql.Named("number", number))

	status := ""

	err := row.Scan(&status)
	if err != nil {
		log.Println("Error scanning data from a row while getting a parcels status")
		return false, err
	}

	return status == ParcelStatusRegistered, nil
}
