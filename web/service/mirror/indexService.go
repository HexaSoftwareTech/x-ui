package mirror

import (
	"gorm.io/gorm"
	"x-ui/database"
	"x-ui/database/model"
)

type IndexService struct {
}

func (s *IndexService) GetMirrors() ([]*model.Mirror, error) {
	db := database.GetDB()
	var mirrors []*model.Mirror
	err := db.Model(model.Mirror{}).Find(&mirrors).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return mirrors, nil
}
