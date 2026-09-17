package models

import (
	e "backend/entities"
	p "backend/entities/packet"
)

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
	err := m.Db.Model(user).Where("user_name = ?", username).Take(user).Error
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

func (m *model) GetNonFriends(userID string, limit int, cursor string) ([]*p.Users, error) {
	users := []*p.Users{}
	query := m.Db.Model(&e.User{}).
		Select("users.id, users.user_name").
		Where("users.id != ?", userID).
		Where("NOT EXISTS (SELECT 1 FROM friends WHERE (friends.from_user_id = users.id AND friends.to_user_id = ?) OR (friends.to_user_id = users.id AND friends.from_user_id = ?))", userID, userID).
		Order("users.id ASC")
	if cursor != "" {
		query = query.Where("users.id > ?", cursor)
	}
	err := query.Limit(limit + 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
