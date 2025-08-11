package gorm

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Client struct {
	ID         string     `gorm:"column:id;primaryKey;type:char(26)"`
	Name       string     `gorm:"column:name;type:text;not null"`
	LastName   string     `gorm:"column:last_name;type:text;not null"`
	EntityType string     `gorm:"column:entity_type;type:text;not null"`
	Document   string     `gorm:"column:document;type:text;not null"`
	BirthDay   *time.Time `gorm:"column:birth_day;type:timestamp"`
	Addresses  []Address  `gorm:"foreignKey:OwnerID"`
	Phones     []Phone    `gorm:"foreignKey:OwnerID"`
	DeletedAt  gorm.DeletedAt
}

func (d *Client) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	return nil
}

func (d *Client) TableName() string {
	return "clients"
}

type Address struct {
	ID           string `gorm:"column:id;primaryKey;type:char(26)"`
	Street       string `gorm:"column:street"`
	Neighborhood string `gorm:"column:neighborhood"`
	City         string `gorm:"column:city"`
	State        string `gorm:"column:state"`
	ZipCode      string `gorm:"column:zip_code"`
	OwnerID      string `gorm:"column:owner_id"`
}

func (d *Address) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	return nil
}

func (d *Address) TableName() string {
	return "addresses"
}

type Phone struct {
	ID          string `gorm:"column:id;primaryKey;type:char(26)"`
	Description string `gorm:"column:description"`
	Number      string `gorm:"column:number"`
	OwnerID     string `gorm:"column:owner_id;type:char(26);not null"`
}

func (d *Phone) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	return nil
}

func (d *Phone) TableName() string {
	return "phones"
}
