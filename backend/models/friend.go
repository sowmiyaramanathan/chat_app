package models

import (
	e "backend/entities"
	p "backend/entities/packet"
)

func (m model) IsFriend(userAID, userBID string) (bool, error) {
	var count int64

	err := m.Db.
		Model(&e.Friends{}).Select("friends.friend_status").
		Where("((from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?)) and friend_status = 'accepted'", userAID, userBID, userBID, userAID).
		Count(&count).Error
	if err != nil {
		return false, err
	} else if count == 0 {
		return false, nil
	}
	return true, nil
}

func (m *model) CreateRequest(userAID, userBID string) error {
	request := &e.Friends{FromUserID: userAID, ToUserID: userBID, FriendStatus: "pending"}

	err := m.Db.Create(request).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *model) GetMyRequests(userID string) ([]*p.Requests, error) {
	requests := []*p.Requests{}

	err := m.Db.
		Model(&e.Friends{}).
		Select("friends.from_user_id", "users.user_name").
		Joins("INNER JOIN users ON friends.from_user_id = users.id").
		Where("friends.to_user_id = ? AND friend_status = 'pending'", userID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (m *model) AcceptRequest(userAID, userBID string) error {
	err := m.Db.Model(&e.Friends{}).Where("from_user_id = ? AND to_user_id = ?", userAID, userBID).Update("friend_status", "accepted").Error
	if err != nil {
		return err
	}
	return nil
}

func (m *model) RejectRequest(userAID, userBID string) error {
	err := m.Db.Model(&e.Friends{}).Where("from_user_id = ? AND to_user_id = ?", userAID, userBID).Update("friend_status", "rejected").Error
	if err != nil {
		return err
	}
	return nil
}

func (m *model) IsRequestSent(userAID, userBID string) (bool, error) {
	var count int64

	err := m.Db.Model(e.Friends{}).
		Select("friends.id", "users.user_name").
		Joins("INNER JOIN users ON friends.from_user_id = users.id").
		Where("friends.from_user_id = ? AND friends.to_user_id = ? AND friend_status = 'pending'", userAID, userBID).Count(&count).Error
	if err != nil {
		return false, err
	} else if count == 0 {
		return false, nil
	}

	return true, nil
}

func (m *model) IsRequestReceived(userAID, userBID string) (bool, error) {
	var count int64

	err := m.Db.Model(e.Friends{}).
		Select("friends.id", "users.user_name").
		Joins("INNER JOIN users ON friends.from_user_id = users.id").
		Where("friends.from_user_id = ? AND friends.to_user_id = ? AND friend_status = 'pending'", userAID, userBID).Count(&count).Error
	if err != nil {
		return false, err
	} else if count == 0 {
		return false, nil
	}

	return true, nil
}
