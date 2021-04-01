package cache

type Cash interface{
	Session() SessionRepository
}