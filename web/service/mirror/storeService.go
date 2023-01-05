package mirror

import (
	"errors"
	"x-ui/database"
	"x-ui/database/model"
)

type StoreService struct {
}

func (*StoreService) checkIpAndPortExist(mirror *model.Mirror) (bool, error) {
	var count int64

	db := database.GetDB()
	err := db.Model(model.Mirror{}).Where("port = ?", mirror.Port).Where("ip = ?", mirror.Ip).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *StoreService) SaveMirror(mirror *model.Mirror) error {
	exists, err := s.checkIpAndPortExist(mirror)
	if err != nil {
		return nil
	}
	if exists {
		return errors.New("record already exists")
	}

	db := database.GetDB()
	return db.Save(mirror).Error
}
