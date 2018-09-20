package jwt

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/endiangroup/compandauth"
)

// Author Author
type Author struct {
	ID   int64
	Name string
}

// SessionClaims SessionClaims
type SessionClaims struct {
	Author  Author
	CAAType CAAType
	CAA     compandauth.SessionCAA // 记录有效用户序号
	Verify  map[string]interface{} // 验证数据，默认存ip,host
	Data    map[string]interface{} // 用户自定义数据
	jwt.StandardClaims
}

// NewSessionClaims NewSessionClaims
func NewSessionClaims() *SessionClaims {
	return &SessionClaims{
		Verify: make(map[string]interface{}),
		Data:   make(map[string]interface{}),
	}
}

// SetAuthor  自定义数据
func (p *SessionClaims) SetAuthor(author Author) *SessionClaims {
	p.Author = author
	return p
}

// GetData GetData
func (p *SessionClaims) GetData(key string) interface{} {
	return p.Data[key]
}

// SetData SetData
func (p *SessionClaims) SetData(key string, value interface{}) *SessionClaims {
	p.Data[key] = value
	return p
}

// VerifyKey VerifyKey
func (p *SessionClaims) VerifyKey(key string, value interface{}) bool {
	return p.Verify[key] == value
}

// SetVerify SetVerify
func (p *SessionClaims) SetVerify(key string, value interface{}) *SessionClaims {
	p.Verify[key] = value
	return p
}

// SetCAA SetCAA
func (p *SessionClaims) SetCAA(caa compandauth.SessionCAA) *SessionClaims {
	p.CAA = caa
	return p
}

// Session Session
type Session struct {
	opts  Options
	store Store
	*jwt.Token
}

// Valid Valid
func (p *Session) Valid() (bool, error) {
	uid := p.GetCliams().Author.ID
	// uid can not be zero
	if uid == 0 {
		fmt.Println("id 00 ")
		return false, nil
	}

	if !p.Token.Valid {
		fmt.Println("token expired")
		return false, nil
	}

	// caa类型
	if p.GetCliams().CAAType != p.opts.CAAType {
		fmt.Println("caa type err")
		return false, nil
	}

	if p.opts.UseCounter() {
		c, err := p.store.GetCounter(uid)
		if err != nil {
			return false, err
		}

		cp := compandauth.Counter(c)
		fmt.Println("caa c isvalid ?")
		return cp.IsValid(p.GetCliams().CAA, int64(p.opts.MaxActive)), nil
	}

	c, err := p.store.GetTimeout(uid)
	if err != nil {
		return false, err
	}

	cp := compandauth.Timeout(c)
	fmt.Println("caa isvalid ?", cp, p.GetCliams().CAA, int64(p.opts.MaxAge)*86400)
	return cp.IsValid(p.GetCliams().CAA, int64(p.opts.MaxAge)*86400), nil
}

// GetCliams GetCliams
func (p *Session) GetCliams() *SessionClaims {
	return p.Claims.(*SessionClaims)
}

// SignedString SignedString 生成SignedString
func (p *Session) SignedString() (string, error) {
	var (
		token string
		err   error
	)

	if p.opts.UseCounter() {
		c, err := p.store.GetCounter(p.GetCliams().Author.ID)
		if err != nil {
			return "", err
		}

		cp := compandauth.Counter(c)
		// @note 并发登陆的时候，可能会出现多个sessionCAA一样的情况，机率很小
		p.GetCliams().SetCAA(cp.Issue())

		fmt.Println("use counter", cp)
		defer func() {
			// counter 每次登陆+1， 要更新store
			err = p.store.SetCounter(p.GetCliams().Author.ID, int64(cp))
		}()
	} else {
		c, err := p.store.GetTimeout(p.GetCliams().Author.ID)
		if err != nil {
			return "", err
		}

		cp := compandauth.Timeout(c)

		// 是否首次分配
		hasFirstIssued := !cp.HasIssued()

		issueUnix := cp.Issue()

		if hasFirstIssued && p.opts.MinNight.True() {
			t := time.Unix(int64(issueUnix), 0)
			issueUnix = compandauth.SessionCAA(MinNight(t).Unix())
		}

		// timout 首次分配时更新就行
		if hasFirstIssued {
			defer func() {
				err = p.store.SetTimeout(p.GetCliams().Author.ID, int64(issueUnix))
			}()
		}

		fmt.Println("use timeout", cp)

		p.GetCliams().SetCAA(issueUnix)

	}

	var startTime time.Time

	if p.opts.MinNight.True() {
		startTime = MinNight()
	} else {
		startTime = time.Now()
	}

	p.GetCliams().ExpiresAt = startTime.AddDate(0, 0, p.opts.MaxAge).Unix()
	p.GetCliams().IssuedAt = time.Now().Unix()

	token, err = p.Token.SignedString(p.opts.GetSignKey(p.opts.JwtSigningMethod))
	if err != nil {
		return "", err
	}

	return token, err
}

// Flush Flush
func (p *Session) Flush() error {
	var err error
	uid := p.GetCliams().Author.ID
	// uid can not be zero
	if uid == 0 {
		return nil
	}

	if p.opts.UseCounter() {
		c, err := p.store.GetCounter(p.GetCliams().Author.ID)
		if err != nil {
			return err
		}

		cp := compandauth.Counter(c)

		// timout 没有分配说明没有登陆过
		if !cp.HasIssued() {
			return nil
		}

		cp.Revoke(int64(p.opts.MaxActive))

		defer func() {
			// counter， 要更新store
			err = p.store.SetCounter(p.GetCliams().Author.ID, int64(cp))
		}()
	} else {
		c, err := p.store.GetTimeout(p.GetCliams().Author.ID)
		if err != nil {
			return err
		}

		cp := compandauth.Timeout(c)

		// timout 没有分配说明没有登陆过
		if !cp.HasIssued() {
			return err
		}

	}

	return err
}
