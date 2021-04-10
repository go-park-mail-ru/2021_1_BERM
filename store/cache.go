package store

type Cash interface{
	Session() SessionRepository
}