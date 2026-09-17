package chat

import (
	"backend/apperrors"
	p "backend/entities/packet"
)

func (c *chat) CheckIsFriend(userAID, userBID string) (bool, error) {
	resp, err := c.m.IsFriend(userAID, userBID)
	if err != nil {
		return false, err
	}

	return resp, nil
}

func (c *chat) CreateFriendRequest(userAID, userBID string) error {
	if userAID == userBID {
		return apperrors.ErrInvalidFriendRequest
	}
	accepted, err := c.CheckIsFriend(userAID, userBID)
	if err != nil {
		return err
	}
	if accepted {
		return apperrors.ErrFriendRequestExists
	}
	sent, err := c.CheckIsFriendRequestSent(userAID, userBID)
	if err != nil {
		return err
	}
	received, err := c.CheckIsFriendRequestReceived(userAID, userBID)
	if err != nil {
		return err
	}
	if sent || received {
		return apperrors.ErrFriendRequestExists
	}
	err = c.m.CreateRequest(userAID, userBID)
	if err != nil {
		return err
	}
	return nil
}

func (c *chat) GetFriendRequests(userID string) ([]*p.Requests, error) {
	requests, err := c.m.GetMyRequests(userID)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (c *chat) AcceptFriendRequest(userAID, userBID string) error {
	err := c.m.AcceptRequest(userAID, userBID)
	if err != nil {
		return err
	}
	return nil
}

func (c *chat) RejecttFriendRequest(userAID, userBID string) error {
	err := c.m.RejectRequest(userAID, userBID)
	if err != nil {
		return err
	}
	return nil
}

func (c *chat) CheckIsFriendRequestSent(userAID, userBID string) (bool, error) {
	resp, err := c.m.IsRequestSent(userAID, userBID)

	if err != nil {
		return false, err
	}

	return resp, nil
}

func (c *chat) CheckIsFriendRequestReceived(userAID, UserBID string) (bool, error) {
	resp, err := c.m.IsRequestReceived(userAID, UserBID)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (c *chat) GetMyFriends(userID string, limit int, cursor string) ([]*p.Users, error) {
	users, err := c.m.GetMyFriends(userID, limit, cursor)
	if err != nil {
		return nil, err
	}

	return users, nil
}
