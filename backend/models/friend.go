package models

import (
	e "backend/entities"
)

func (m model) IsFriend(userAID, userBID uint64) (bool, error) {
	var count int64

	err := m.Db.Debug().
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

func (m *model) CreateRequest(userAID, userBID uint64) error {
	request := &e.Friends{FromUserID: userAID, ToUserID: userBID, FriendStatus: "pending"}

	err := m.Db.Debug().Create(request).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *model) GetMyRequests(userID uint64) ([]e.Requests, error) {
	requests := []e.Requests{}

	err := m.Db.Debug().
		Model(&e.Friends{}).
		Select("friends.from_user_id", "users.username").
		Joins("INNER JOIN users ON friends.from_user_id = users.id").
		Where("friends.to_user_id = ? AND friend_status = 'pending'", userID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (m *model) AcceptRequest(userAID, userBID uint64) error {
	err := m.Db.Debug().Model(&e.Friends{}).Where("from_user_id = ? AND to_user_id = ?", userAID, userBID).Update("friend_status", "accepted").Error
	if err != nil {
		return err
	}
	return nil
}

func (m *model) RejectRequest(userAID, userBID uint64) error {
	err := m.Db.Debug().Model(&e.Friends{}).Where("from_user_id = ? AND to_user_id = ?", userAID, userBID).Update("friend_status", "rejected").Error
	if err != nil {
		return err
	}
	return nil
}

func (m *model) IsRequestSent(userAID, userBID uint64) (bool, error) {
	var count int64

	err := m.Db.Debug().Model(e.Friends{}).
		Select("friends.id", "users.username").
		Joins("INNER JOIN users ON friends.from_user_id = users.id").
		Where("friends.from_user_id = ? AND friends.to_user_id = ? AND friend_status = 'pending'", userAID, userBID).Count(&count).Error
	if err != nil {
		return false, err
	} else if count == 0 {
		return false, nil
	}

	return true, nil
}

func (m *model) IsRequestReceived(userAID, userBID uint64) (bool, error) {
	var count int64

	err := m.Db.Debug().Model(e.Friends{}).
		Select("friends.id", "users.username").
		Joins("INNER JOIN users ON friends.from_user_id = users.id").
		Where("friends.from_user_id = ? AND friends.to_user_id = ? AND friend_status = 'pending'", userAID, userBID).Count(&count).Error
	if err != nil {
		return false, err
	} else if count == 0 {
		return false, nil
	}

	return true, nil
}
