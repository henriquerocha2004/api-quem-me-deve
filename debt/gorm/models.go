package gorm

import (
	"time"

	"github.com/lib/pq"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Debt struct {
	ID                   string         `gorm:"column:id;primaryKey;type:char(26)"`
	Description          string         `gorm:"column:description;type:text;not null"`
	TotalValue           float64        `gorm:"column:total_value;type:decimal(10,2);not null"`
	DueDate              *time.Time     `gorm:"column:due_date;type:timestamp;not null"`
	InstallmentsQuantity int            `gorm:"column:installments_quantity;type:int;not null"`
	UserClientId         string         `gorm:"column:user_client_id;type:char(26);not null"`
	ProductIds           pq.StringArray `gorm:"column:product_ids;type:text[];not null"`
	ServiceIds           pq.StringArray `gorm:"column:service_ids;type:text[];not null"`
	Status               string         `gorm:"column:status;type:text;not null"`
	DebtDate             *time.Time     `gorm:"column:debt_date;type:timestamp;not null"`
	Installments         []Installment  `gorm:"foreignKey:DebtId"`
	CancelInfo           CancelInfo     `gorm:"foreignKey:DebtId"`
	ReversalInfo         ReversalInfo   `gorm:"foreignKey:DebtId"`
}

func (d *Debt) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = ulid.Make().String()
	}
	return nil
}

func (d *Debt) TableName() string {
	return "debts"
}

type Installment struct {
	Id            string     `gorm:"column:id"`
	Description   string     `gorm:"column:description"`
	Value         float64    `gorm:"column:value"`
	DueDate       *time.Time `gorm:"column:due_date"`
	DebDate       *time.Time `gorm:"column:deb_date"`
	Status        string     `gorm:"column:status"`
	PaymentDate   *time.Time `gorm:"column:payment_date"`
	PaymentMethod string     `gorm:"column:payment_method"`
	Number        int        `gorm:"column:number"`
	DebtId        string     `gorm:"column:debt_id"`
}

func (d *Installment) BeforeCreate(tx *gorm.DB) (err error) {
	if d.Id == "" {
		d.Id = ulid.Make().String()
	}
	return nil
}

func (d *Installment) TableName() string {
	return "installments"
}

type CancelInfo struct {
	Id          string     `gorm:"column:id"`
	Reason      string     `gorm:"column:reason"`
	CancelDate  *time.Time `gorm:"column:cancel_date"`
	CancelledBy string     `gorm:"column:cancelled_by"`
	DebtId      string     `gorm:"column:debt_id"`
}

func (d *CancelInfo) BeforeCreate(tx *gorm.DB) (err error) {
	if d.Id == "" {
		d.Id = ulid.Make().String()
	}
	return nil
}

func (d *CancelInfo) TableName() string {
	return "cancel_info"
}

type ReversalInfo struct {
	Id                      string     `gorm:"column:id"`
	Reason                  string     `gorm:"column:reason"`
	ReversalDate            *time.Time `gorm:"column:reversal_date"`
	ReversedBy              string     `gorm:"column:reversed_by"`
	ReversedInstallmentQtd  int        `gorm:"column:reversed_installment_qtd"`
	CancelledInstallmentQtd int        `gorm:"column:cancelled_installment_qtd"`
	DebtId                  string     `gorm:"column:debt_id"`
}

func (d *ReversalInfo) BeforeCreate(tx *gorm.DB) (err error) {
	if d.Id == "" {
		d.Id = ulid.Make().String()
	}
	return nil
}

func (d *ReversalInfo) TableName() string {
	return "reversal_info"
}
