package part

import (
	repoModel "Jopa/inventory/internal/repository/model"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func initParts() []interface{} {
	fmt.Printf("🔍 initParts() called - preparing %d test parts", 3)
	stringPtr := "ZaebicyVersia1.0"
	part1 := repoModel.Part{
		Uuid:          uuid.New().String(),
		Name:          "Крылышко",
		Description:   "Эта огненный длакон",
		Price:         1000.0,
		StockQuantity: 2,
		Category:      string(repoModel.CategoryWing),
		Dimensions: repoModel.Dimensions{
			Length: 10.0,
			Width:  10.0,
			Height: 10.0,
			Weight: 15.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "ChinaHolocoColpolation",
			Country: "China",
			Website: "www.chinacloudapi.cn",
		},
		Tags: []string{"Крыло", "дракон", "нефритовый стержень"},
		Metadata: map[string]repoModel.MetaValue{
			"Versia:": {
				StringValue: &stringPtr},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	stringPtr2 := "ГОСТ 1 1010-57"
	part2 := repoModel.Part{
		Uuid:          uuid.New().String(),
		Name:          "Бак",
		Description:   "Бак-колпак",
		Price:         100.0,
		StockQuantity: 10,
		Category:      string(repoModel.CategoryFuel),
		Dimensions: repoModel.Dimensions{
			Length: 100.0,
			Width:  100.0,
			Height: 100.0,
			Weight: 150.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "UnitedTank-buildingCorp",
			Country: "Russia",
			Website: "www.russiaTankCorp.ru",
		},
		Tags: []string{"Мечта герметизатора"},
		Metadata: map[string]repoModel.MetaValue{
			"Версия бака": {
				StringValue: &stringPtr2},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	stringPtr3 := "V0.00015"
	part3 := repoModel.Part{
		Uuid:          uuid.New().String(),
		Name:          "Двигатель",
		Description:   "F3500. С тягой СтоТыщьМельенов",
		Price:         100_000_000.0,
		StockQuantity: 1,
		Category:      string(repoModel.CategoryEngine),
		Dimensions: repoModel.Dimensions{
			Length: 15.0,
			Width:  5.0,
			Height: 5.0,
			Weight: 30.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "Capitalist1 & Capitalist2 & Co",
			Country: "USA",
			Website: "www.Amazon-Fabrica№394.com",
		},
		Tags: []string{"дорого/богато"},
		Metadata: map[string]repoModel.MetaValue{
			"Версия двигателя": {
				StringValue: &stringPtr3},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return []interface{}{part1, part2, part3}
}
