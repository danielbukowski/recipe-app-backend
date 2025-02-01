package session

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	DefaultSessionExpirationTime = 86400 * 14
	storageSessionKeyLength      = 20
	sessionStorageKey            = "session_id"

	SessionCookieName = "SESSION_ID"
)

// Session represents stored values in memcache.
type Session struct {
	Email string `json:"email"`
}

// Middlewares adds the stored session from the memcache to the request context.
func Middleware(memcachedStore *MemcachedStore, skipper middleware.Skipper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			session := Session{}

			sessionValue, err := memcachedStore.Get(c)
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					// Pass the request with empty session
					c.Set("session", &session)
					return next(c)
				case errors.Is(err, memcache.ErrCacheMiss):
					// The session cookie does not exist in the cache, so just delete the cookie from client.
					// Also pass the request with empty session.
					memcachedStore.deleteCookieFromClient(c)
					c.Set("session", &session)
					return next(c)
				default:
					return err
				}
			}

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
	isDev           bool
}

// NewSessionStorage returns a new instance of MemcachedStore.
func NewSessionStorage(cacheClient *memcache.Client, isDev bool) *MemcachedStore {
	return &MemcachedStore{
		memcachedClient: cacheClient,
		itemPool: &sync.Pool{
			New: func() interface{} {
				return new(memcache.Item)
			},
		},
		isDev: isDev,
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

	err := ms.memcachedClient.Add(item)
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

// Delete deletes session from all storages.
func (ms *MemcachedStore) Delete(c echo.Context) {
	cookie, err := c.Cookie(SessionCookieName)
	if err != nil {
		return
	}

	_ = ms.memcachedClient.Delete(cookie.Value)

	ms.deleteCookieFromClient(c)
}

// DeleteCookieFromClient deletes a session cookie from a client's browser.
func (ms *MemcachedStore) deleteCookieFromClient(c echo.Context) {
	cookie := http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   !ms.isDev,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(&cookie)
}

// AttachSessionCookieToClient saves session cookie to a client's browser.
func (ms *MemcachedStore) AttachSessionCookieToClient(sessionID string, c echo.Context) {
	cookie := http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   DefaultSessionExpirationTime,
		Secure:   !ms.isDev,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(&cookie)
}
