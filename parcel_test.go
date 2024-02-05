package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Error opening DB")

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	require.NoError(t, err, "Adding parcel must not return error")
	require.NotNil(t, number, "Number must not be nil")
	require.NotEqual(t, 0, number, "Number must not be 0")

	parcel.Number = number

	// get
	parcelFromDb, err := store.Get(number)
	require.NoError(t, err, "Getting parcel from DB must not return error")
	require.Equal(t, parcel, parcelFromDb, "ALl fields of parcel must match")

	// delete
	// удаляем добавленную посылку
	err2 := store.Delete(number)
	require.NoError(t, err2)
	// проверяем, что посылку больше нельзя получить из БД
	parcelFromDb2, err2 := store.Get(number)
	require.ErrorIs(t, sql.ErrNoRows, err2)
	require.Equal(t, parcelFromDb2, *new(Parcel))
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Error opening DB")

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	require.NoError(t, err, "Adding parcel must not return error")
	require.NotNil(t, number, "Number must not be nil")
	require.NotEqual(t, 0, number, "Number must not be 0")

	parcel.Number = number

	// set address
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)


	// check
	parcelFromDb, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, newAddress, parcelFromDb.Address)
}

// TestSetAddress_WhenNotRegistered проверяет обновление адреса для посылки с неверным статусом
func TestSetAddress_WhenNotRegistered(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Error opening DB")

	store := NewParcelStore(db)
	parcel := getTestParcel()
	parcel.Status = ParcelStatusSent

	// add
	number, err := store.Add(parcel)
	require.Error(t, err, "Adding parcel must not return error")
	require.NotNil(t, number, "Number must not be nil")
	require.NotEqual(t, 0, number, "Number must not be 0")

	parcel.Number = number

	// set address
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.ErrorContains(t, err, "not registered")


	// check
	parcelFromDb, err := store.Get(number)
	require.NoError(t, err)
	require.NotEqual(t, newAddress, parcelFromDb.Address)
	require.Equal(t, parcel.Address, parcelFromDb.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Error opening DB")

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	require.NoError(t, err, "Adding parcel must not return error")
	require.NotNil(t, number, "Number must not be nil")
	require.NotEqual(t, 0, number, "Number must not be 0")

	parcel.Number = number
	newStatus := ParcelStatusDelivered

	// set status
	err = store.SetStatus(number, newStatus)
	require.NoError(t, err)

	// check
	parcelFromDb, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, newStatus, parcelFromDb.Status)

}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Error opening DB")

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	}
}
