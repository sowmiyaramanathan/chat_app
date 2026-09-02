package models

import e "backend/entities"

func (m *model) SaveUser(user *e.User) (*e.User, error) {
	err := m.Db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// func (m *model) GetUserById(id uint64) (*e.User, error) {
// 	user := &e.User{}
// 	err := m.Db.First(user, id).Take(user).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

func (m *model) GetUserByUsername(username string) (*e.User, error) {
	user := &e.User{}
	err := m.Db.Model(user).Where("username = ?", username).Take(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *model) GetUserByMobilenumber(number string) (*e.User, error) {
	user := &e.User{}
	err := m.Db.Model(user).Where("mobilenumber = ?", number).Take(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *model) GetUsers(username string) (*[]e.Users, error) {
	users := &[]e.Users{}
	err := m.Db.Model(&e.Users{}).Select("id", "username").Where("username != ?", username).Limit(100).Find(users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (m *model) GetPublicKey(userID uint64) (string, error) {
	user := &e.User{}
	err := m.Db.Model(&e.User{}).Where("id = ?", userID).Take(user).Error
	if err != nil {
		return "", err
	}
	return user.PubKey, nil
}
