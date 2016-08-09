package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrKeyNotFound         = errors.New("key not found in session values")
	ErrTypeAssertionFailed = errors.New("type assertion failed")
)

func GetString(r *http.Request, key string) (string, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return "", err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	v, exists := s.values[key]
	if exists == false {
		return "", ErrKeyNotFound
	}

	str, ok := v.(string)
	if ok == false {
		return "", ErrTypeAssertionFailed
	}

	return str, nil
}

func PutString(r *http.Request, key string, val string) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	s.values[key] = val
	s.modified = true
	return nil
}

func PopString(r *http.Request, key string) (string, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return "", ErrAlreadyWritten
	}
	v, exists := s.values[key]
	if exists == false {
		return "", ErrKeyNotFound
	}

	str, ok := v.(string)
	if ok == false {
		return "", ErrTypeAssertionFailed
	}

	delete(s.values, key)
	s.modified = true
	return str, nil
}

func GetBool(r *http.Request, key string) (bool, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return false, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	v, exists := s.values[key]
	if exists == false {
		return false, ErrKeyNotFound
	}

	b, ok := v.(bool)
	if ok == false {
		return false, ErrTypeAssertionFailed
	}

	return b, nil
}

func PutBool(r *http.Request, key string, val bool) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	s.values[key] = val
	s.modified = true
	return nil
}

func PopBool(r *http.Request, key string) (bool, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return false, ErrAlreadyWritten
	}
	v, exists := s.values[key]
	if exists == false {
		return false, ErrKeyNotFound
	}

	b, ok := v.(bool)
	if ok == false {
		return false, ErrTypeAssertionFailed
	}

	delete(s.values, key)
	s.modified = true
	return b, nil
}

func GetInt(r *http.Request, key string) (int, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return 0, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	v, exists := s.values[key]
	if exists == false {
		return 0, ErrKeyNotFound
	}

	switch v := v.(type) {
	case int:
		return v, nil
	case json.Number:
		return strconv.Atoi(v.String())
	}

	return 0, ErrTypeAssertionFailed
}

func PutInt(r *http.Request, key string, val int) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	s.values[key] = val
	s.modified = true
	return nil
}

func PopInt(r *http.Request, key string) (int, error) {
	s, err := sessionFromContext(r)
	if err != nil {
		return 0, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.written == true {
		return 0, ErrAlreadyWritten
	}
	v, exists := s.values[key]
	if exists == false {
		return 0, ErrKeyNotFound
	}

	var i int
	switch v := v.(type) {
	case int:
		i = v
	case json.Number:
		i, err = strconv.Atoi(v.String())
		if err != nil {
			return 0, err
		}
	default:
		return 0, ErrTypeAssertionFailed
	}

	delete(s.values, key)
	s.modified = true
	return i, nil
}

func Remove(r *http.Request, key string) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	delete(s.values, key)
	s.modified = true
	return nil
}

func Clear(r *http.Request) error {
	s, err := sessionFromContext(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.written == true {
		return ErrAlreadyWritten
	}
	for key := range s.values {
		delete(s.values, key)
	}
	s.modified = true
	return nil
}