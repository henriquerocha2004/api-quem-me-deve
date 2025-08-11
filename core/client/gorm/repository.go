package gorm

import (
	"context"
	"errors"

	"github.com/henriquerocha2004/quem-me-deve-api/core/client"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/document"
	"github.com/henriquerocha2004/quem-me-deve-api/pkg/paginate"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GormClientRepository struct {
	db *gorm.DB
}

func NewGormClientRepository(db *gorm.DB) *GormClientRepository {
	return &GormClientRepository{db: db}
}

func (c *GormClientRepository) Create(ctx context.Context, client *client.Client) error {

	clientModel := Client{
		ID:         client.Id.String(),
		Name:       client.Name,
		LastName:   client.LastName,
		EntityType: string(client.EntityType),
		Document:   string(client.Document),
		BirthDay:   client.BirthDay,
		Addresses:  c.convertAddressToModel(client.Addresses, client.Id),
		Phones:     c.convertPhoneToModel(client.Phones, client.Id),
	}

	tx := c.db.Begin()

	if err := tx.Create(&clientModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (c *GormClientRepository) Update(ctx context.Context, client *client.Client) error {
	clientModel := &Client{
		ID:         client.Id.String(),
		Name:       client.Name,
		LastName:   client.LastName,
		EntityType: string(client.EntityType),
		Document:   string(client.Document),
		BirthDay:   client.BirthDay,
		Addresses:  c.convertAddressToModel(client.Addresses, client.Id),
		Phones:     c.convertPhoneToModel(client.Phones, client.Id),
	}

	tx := c.db.Begin()

	err := tx.Model(&Client{}).Where("id = ?", clientModel.ID).Updates(clientModel).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(clientModel.Addresses) != 0 {
		err = tx.WithContext(ctx).Model(clientModel).
			Association("Addresses").
			Replace(clientModel.Addresses)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(clientModel.Phones) != 0 {
		err = tx.WithContext(ctx).Model(clientModel).
			Association("Phones").
			Replace(clientModel.Phones)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (c *GormClientRepository) Delete(ctx context.Context, id ulid.ULID) error {
	tx := c.db.Begin()

	if err := tx.Where("id = ?", id.String()).Delete(&Client{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (c *GormClientRepository) FindById(ctx context.Context, id ulid.ULID) (*client.Client, error) {
	var clientModel Client

	result := c.db.WithContext(ctx).Where("id = ?", id.String()).
		Preload("Addresses").
		Preload("Phones").
		First(&clientModel)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("client not found")
		}
		return nil, result.Error
	}

	return &client.Client{
		Id:         id,
		Name:       clientModel.Name,
		LastName:   clientModel.LastName,
		EntityType: client.EntityType(clientModel.EntityType),
		Document:   document.Document(clientModel.Document),
		BirthDay:   clientModel.BirthDay,
		Addresses:  c.convertModelAddressToDomainAddress(clientModel.Addresses),
		Phones:     c.convertModelPhoneToDomainPhone(clientModel.Phones),
	}, nil
}

func (c *GormClientRepository) FindByDocument(ctx context.Context, doc string) (*client.Client, error) {
	var clientModel Client

	result := c.db.WithContext(ctx).Where("document = ?", doc).
		Preload("Addresses").
		Preload("Phones").
		First(&clientModel)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("client not found")
		}
		return nil, result.Error
	}

	id, err := ulid.Parse(clientModel.ID)
	if err != nil {
		return nil, err
	}

	return &client.Client{
		Id:         id,
		Name:       clientModel.Name,
		LastName:   clientModel.LastName,
		EntityType: client.EntityType(clientModel.EntityType),
		Document:   document.Document(clientModel.Document),
		BirthDay:   clientModel.BirthDay,
		Addresses:  c.convertModelAddressToDomainAddress(clientModel.Addresses),
		Phones:     c.convertModelPhoneToDomainPhone(clientModel.Phones),
	}, nil
}

func (c *GormClientRepository) FindAll(ctx context.Context, criteria paginate.SearchDto) (*client.PaginationResult, error) {
	var models []Client
	var total int64

	query := c.db.Model(&Client{}).
		Count(&total).
		Offset(criteria.Offset()).
		Limit(criteria.Limit).
		Order("created_at DESC").
		Preload("Addresses").
		Preload("Phones")

	if criteria.TermSearch != "" {
		query = query.
			Where("name LIKE ?", "%"+criteria.TermSearch+"%").
			Or("last_name LIKE ?", "%"+criteria.TermSearch+"%").
			Or("document LIKE ?", "%"+criteria.TermSearch+"%")
	}

	if len(criteria.ColumnSearch) >= 1 {
		for _, value := range criteria.ColumnSearch {
			query = query.Where(value.ColumnName+" = ?", value.ColumnValue)
		}
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	var clients []*client.Client

	for _, model := range models {
		clients = append(clients, c.convertClientModelToDomain(model))
	}

	return &client.PaginationResult{
		TotalRecords: int(total),
		Data:         clients,
	}, nil
}

func (c *GormClientRepository) convertAddressToModel(address []client.Address, clientId ulid.ULID) []Address {
	var addressModel []Address
	for _, addr := range address {
		addressModel = append(addressModel, Address{
			ID:           addr.Id.String(),
			Street:       addr.Street,
			Neighborhood: addr.Neighborhood,
			City:         addr.City,
			State:        addr.State,
			ZipCode:      addr.ZipCode,
			OwnerID:      clientId.String(),
		})
	}
	return addressModel
}

func (c *GormClientRepository) convertPhoneToModel(phones []client.Phone, clientId ulid.ULID) []Phone {
	var phoneModel []Phone
	for _, phone := range phones {
		phoneModel = append(phoneModel, Phone{
			ID:          phone.Id.String(),
			Description: phone.Description,
			Number:      phone.Number,
			OwnerID:     clientId.String(),
		})
	}
	return phoneModel
}

func (c *GormClientRepository) convertModelAddressToDomainAddress(addresses []Address) []client.Address {
	var domainAddresses []client.Address
	for _, addr := range addresses {
		domainAddresses = append(domainAddresses, client.Address{
			Street:       addr.Street,
			Neighborhood: addr.Neighborhood,
			City:         addr.City,
			State:        addr.State,
			ZipCode:      addr.ZipCode,
		})
	}
	return domainAddresses
}

func (c *GormClientRepository) convertModelPhoneToDomainPhone(phones []Phone) []client.Phone {
	var domainPhones []client.Phone
	for _, phone := range phones {
		domainPhones = append(domainPhones, client.Phone{
			Description: phone.Description,
			Number:      phone.Number,
		})
	}
	return domainPhones
}

func (c *GormClientRepository) convertClientModelToDomain(clientModel Client) *client.Client {
	id, err := ulid.Parse(clientModel.ID)
	if err != nil {
		return nil
	}

	return &client.Client{
		Id:         id,
		Name:       clientModel.Name,
		LastName:   clientModel.LastName,
		EntityType: client.EntityType(clientModel.EntityType),
		Document:   document.Document(clientModel.Document),
		BirthDay:   clientModel.BirthDay,
		Addresses:  c.convertModelAddressToDomainAddress(clientModel.Addresses),
		Phones:     c.convertModelPhoneToDomainPhone(clientModel.Phones),
	}
}

type ClientReaderGormRepository struct {
	db *gorm.DB
}

func NewClientReaderGormRepository(db *gorm.DB) *ClientReaderGormRepository {
	return &ClientReaderGormRepository{db: db}
}

func (c *ClientReaderGormRepository) ClientExists(ctx context.Context, id ulid.ULID) (bool, error) {
	var count int64
	err := c.db.WithContext(ctx).Model(&Client{}).Where("id = ?", id.String()).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
