package structs

import "errors"

type AllRooms map[string]Room

func (all *AllRooms) AddRoom(roomId string) error {
	if all == nil {
		return errors.New("nil pointer, all rooms is not accessible here")
	}

	(*all)[roomId] = Room{}
	return nil
}

func (all *AllRooms) DeleteRoom(roomId string) error {
	if all == nil {
		return errors.New("nil pointer, all rooms is not accessible here")
	}

	delete(*all, roomId)

	return nil
}
