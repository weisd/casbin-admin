package jwt

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Errors
var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")
)

// SessionManager SessionManager
type SessionManager struct {
	store Store
	opts  Options
}

// Options Options
type Options struct {
	MaxActive int // 同一uid 有效jwt数量 默认：0 无限
	MaxAge    int // jwt 有效天数   默认：3天 0点到期

	// StoreEnable  bool   // session更新失效机制需要开户store
	StoreAdapter string // store 存储方式 mysql,redis,memory,file...
	StoreConfig  string // store Adapter 对应的配置

	// JwtLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	JwtLookup        string            // 默认 header:Authorization
	JwtScheme        string            // token前缀 默认 Bearer
	JwtSigningMethod jwt.SigningMethod // jwt 算法
	JwtPrivateKey    []byte            // 只需要一个的时候存这
	JwtPublicKey     []byte            //
}

func preOptions(opts *Options) {
	if opts.MaxAge == 0 {
		opts.MaxAge = 3
	}

	if len(opts.StoreAdapter) == 0 {
		opts.StoreAdapter = "memery"
	}

	if len(opts.JwtLookup) == 0 {
		opts.JwtLookup = "header:Authorization"
	}

	if len(opts.JwtScheme) == 0 {
		opts.JwtScheme = "Bearer"
	}

	if opts.JwtSigningMethod == nil {
		panic("JWT SigningMethod is required")
	}

	if len(opts.JwtPrivateKey) == 0 {
		panic("JWT key is required")
	}

}

func (p Options) UseCounter() bool {
	return p.MaxActive > 0
}

func (p Options) GetSignKey(alg jwt.SigningMethod) []byte {
	switch alg {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		return p.JwtPrivateKey
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512, jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		// 私钥签名
		return p.JwtPrivateKey
	default:
		return []byte{}
	}
}

func (p Options) GetVerifyKey(alg jwt.SigningMethod) []byte {
	switch alg {
	case jwt.SigningMethodHS256, jwt.SigningMethodHS384, jwt.SigningMethodHS512:
		return p.JwtPrivateKey
	case jwt.SigningMethodRS256, jwt.SigningMethodRS384, jwt.SigningMethodRS512, jwt.SigningMethodES256, jwt.SigningMethodES384, jwt.SigningMethodES512:
		// 公钥验证
		return p.JwtPublicKey
	default:
		return []byte{}
	}
}

// NewSessionManger NewSessionManger
func NewSessionManger(opts ...Options) *SessionManager {

	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	preOptions(&opt)

	m := &SessionManager{}
	m.opts = opt
	m.store = NewStore(opt.StoreAdapter, opt.StoreConfig)

	return m
}

// GetSession GetSession
func (p *SessionManager) GetSession(r *http.Request) (*Session, error) {
	// 从header中取到jwt

	parts := strings.Split(p.opts.JwtLookup, ":")
	extractor := jwtFromHeader(parts[1], p.opts.JwtScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	}

	tokenString := extractor(r)

	var token *jwt.Token
	if len(tokenString) == 0 {
		token = jwt.NewWithClaims(p.opts.JwtSigningMethod, NewSessionClaims())
	} else {
		var err error
		token, err = p.ParseJWT(tokenString)
		if err != nil {
			return nil, err
		}
	}

	sess := p.NewSession(token)
	sess.GetCliams().SetVerify("ip", RealIP(r))
	return sess, nil

}

// ParseJWT ParseJWT
func (p *SessionManager) ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.opts.GetVerifyKey(token.Method), nil
	})

	return token, err
}

// NewSession NewSession
func (p *SessionManager) NewSession(token ...*jwt.Token) *Session {
	var t *jwt.Token
	if len(token) > 0 {
		t = token[0]
	} else {
		t = jwt.NewWithClaims(p.opts.JwtSigningMethod, NewSessionClaims())
	}
	return &Session{
		store: p.store,
		opts:  p.opts,
		Token: t,
	}
}
