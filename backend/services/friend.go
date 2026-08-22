package services

import (
	e "backend/entities"
	"errors"
)

func (s *service) CheckIsFriend(userAID, userBID uint64) (bool, error) {
	resp, err := s.m.IsFriend(userAID, userBID)

	if err != nil {
		return false, err
	}

	return resp, nil
}

func (s *service) CreateFriendRequest(userAID, userBID uint64) error {
	err := s.m.CreateRequest(userAID, userBID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetFriendRequests(userID uint64) ([]e.Requests, error) {
	requests, err := s.m.GetMyRequests(userID)
	if err != nil {
		return nil, errors.New("could not get friend requests")
	}

	return requests, nil
}

func (s *service) AcceptFriendRequest(userAID, userBID uint64) error {
	err := s.m.AcceptRequest(userAID, userBID)
	if err != nil {
		return errors.New("could not accept friend request")
	}
	return nil
}

func (s *service) RejecttFriendRequest(userAID, userBID uint64) error {
	err := s.m.RejectRequest(userAID, userBID)
	if err != nil {
		return errors.New("could not reject friend request")
	}
	return nil
}

func (s *service) CheckIsFriendRequestSent(userAID, userBID uint64) (bool, error) {
	resp, err := s.m.IsRequestSent(userAID, userBID)

	if err != nil {
		return false, err
	}

	return resp, nil
}

func (s *service) CheckIsFriendRequestReceived(userAID, UserBID uint64) (bool, error) {
	resp, err := s.m.IsRequestReceived(userAID, UserBID)
	if err != nil {
		return false, err
	}
	return resp, nil
}
