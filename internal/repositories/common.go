package repositories

import (
	"LinuxOnM/internal/constant"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type DBOption func(*gorm.DB) *gorm.DB

type ICommonRepository interface {
	WithByID(id uint) DBOption
	WithByName(name string) DBOption
	WithByType(t string) DBOption
	WithOrderBy(orderStr string) DBOption
	WithOrderRuleBy(orderBy, order string) DBOption
	WithByGroupID(groupID uint) DBOption
	WithLikeName(name string) DBOption
	WithIDsIn(ids []uint) DBOption
	WithByDate(startTime, endTime time.Time) DBOption
	WithByStartDate(StartTime time.Time) DBOption
	WithByFrom(from string) DBOption
	WithByStatus(status string) DBOption
}

type CommonRepository struct{}

func NewCommonRepository() ICommonRepository {
	return &CommonRepository{}
}

func (c *CommonRepository) WithByID(id uint) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (c *CommonRepository) WithByName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (c *CommonRepository) WithByType(t string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type = ?", t)
	}
}

func (c *CommonRepository) WithOrderBy(orderStr string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(orderStr)
	}
}

func (c *CommonRepository) WithOrderRuleBy(orderBy, order string) DBOption {
	switch order {
	case constant.OrderDesc:
		order = "desc"
	case constant.OrderAsc:
		order = "asc"
	default:
		orderBy = "created_at"
		order = "desc"
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s %s", orderBy, order))
	}
}

func (c *CommonRepository) WithByGroupID(groupID uint) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if groupID == 0 {
			return db
		}
		return db.Where("group_id = ?", groupID)
	}
}

func (c *CommonRepository) WithLikeName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(name) == 0 {
			return db
		}
		return db.Where("name LIKE ?", "%"+name+"%")
	}
}

func (c *CommonRepository) WithIDsIn(ids []uint) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN (?)", ids)
	}
}

func (c *CommonRepository) WithByDate(startTime, endTime time.Time) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at > ?", startTime).Where("created_at < ?", endTime)
	}
}

func (c *CommonRepository) WithByStartDate(startTime time.Time) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at > ?", startTime)
	}
}

func (c *CommonRepository) WithByFrom(from string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("from = ?", from)
	}
}

func (c *CommonRepository) WithByStatus(status string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(status) == 0 {
			return g
		}
		return g.Where("status = ?", status)
	}
}
