package session

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ SessionStore = (*MemcachedStore)(nil)

const (
	DefaultSessionExpirationTime = 86400 * 7
	storageSessionKeyLength      = 20
	sessionStorageKey            = "session_id"

	SessionCookieName = "SESSION_ID"
)

// SessionStore defines methods for manipulating sessions in memcache.
type SessionStore interface {
	Get(c echo.Context) ([]byte, error)
	CreateNew(value []byte) (string, error)
	Update(key string, value []byte, expiration int32) error
	Delete(key string) error
}

// Session represents stored values in memcache.
type Session struct {
	Email string `json:"email"`
}

// Middlewares adds the stored session from the memcache to the request context.
func Middleware(sessionStore SessionStore, skipper middleware.Skipper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			sessionValue, err := sessionStore.Get(c)
			if err != nil {
				return err
			}

			session := Session{}

			if err := json.Unmarshal(sessionValue, &session); err != nil {
				panic(errors.Join(errors.New("failed to decode the value from session"), err))
			}

			c.Set("session", &session)

			return next(c)
		}
	}
}

// MemcachedStore implements methods for managing the memcached session.
type MemcachedStore struct {
	memcachedClient *memcache.Client
	itemPool        *sync.Pool // used sync.Pool to reduce memcached.Item allocation
}

// NewSessionStorage returns a new instance of MemcachedStore.
func NewSessionStorage(cacheClient *memcache.Client) *MemcachedStore {
	return &MemcachedStore{
		memcachedClient: cacheClient,
		itemPool: &sync.Pool{
			New: func() interface{} {
				return new(memcache.Item)
			},
		},
	}
}

// generateSessionID creates a random string used as a key to memcached.
func generateSessionID() string {
	buf := make([]byte, storageSessionKeyLength)

	if _, err := rand.Read(buf); err != nil {
		panic("failed to generate random bytes for session ID")
	}

	return base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(buf)
}

// getItemFromThePool returns an instance of memcache.Item from the pool.
func (ms *MemcachedStore) getItemFromThePool() *memcache.Item {
	return ms.itemPool.Get().(*memcache.Item)
}

// returnItemToThePool clears the state of the item and it returns to the pool.
func (ms *MemcachedStore) returnItemToThePool(item *memcache.Item) {
	item.Key = ""
	item.Value = nil
	item.Expiration = 0

	ms.itemPool.Put(item)
}

// Get fetches the value of SESSION_ID from a cookie in the echo context.
func (ms *MemcachedStore) Get(c echo.Context) ([]byte, error) {
	cookieValue, err := c.Cookie(SessionCookieName)
	if err != nil {
		return nil, err
	}

	item, err := ms.memcachedClient.Get(cookieValue.Value)
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

// CreateNew creates and saves an entirely new session to memcached.
// The returned string type is a ID of the newly created session.
func (ms *MemcachedStore) CreateNew(value []byte) (string, error) {
	generatedSessionID := generateSessionID()

	item := ms.getItemFromThePool()
	defer ms.returnItemToThePool(item)

	item.Key = generatedSessionID
	item.Value = value
	item.Expiration = DefaultSessionExpirationTime

	err := ms.memcachedClient.Set(item)
	if err != nil {
		return "", err
	}

	return generatedSessionID, nil
}

// Update updates already existing session in memcached.
func (ms *MemcachedStore) Update(key string, value []byte, expiration int32) error {
	item := ms.getItemFromThePool()
	defer ms.returnItemToThePool(item)

	item.Key = key
	item.Value = value
	item.Expiration = expiration

	err := ms.memcachedClient.Replace(item)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes session by the key in memcached.
func (ms *MemcachedStore) Delete(key string) error {
	return ms.memcachedClient.Delete(key)
}
