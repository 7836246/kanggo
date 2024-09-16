package encryptcookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/7836246/kanggo"
	"github.com/7836246/kanggo/core"
)

// Config 是 encryptcookie 中间件的配置结构体
type Config struct {
	Next      func(c *kanggo.Context) bool                      // 可选：跳过此中间件的函数
	Except    []string                                          // 可选：不需要加密的 Cookie 名称列表
	Key       string                                            // 必填：Base64 编码的唯一密钥，用于加密和解密 Cookie
	Encryptor func(decryptedString, key string) (string, error) // 可选：自定义加密函数
	Decryptor func(encryptedString, key string) (string, error) // 可选：自定义解密函数
}

// ConfigDefault 默认配置
var ConfigDefault = Config{
	Next:      nil,
	Except:    []string{},
	Key:       "",
	Encryptor: EncryptCookie,
	Decryptor: DecryptCookie,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.Next == nil {
			cfg.Next = ConfigDefault.Next
		}

		if cfg.Except == nil {
			cfg.Except = ConfigDefault.Except
		}

		if cfg.Encryptor == nil {
			cfg.Encryptor = ConfigDefault.Encryptor
		}

		if cfg.Decryptor == nil {
			cfg.Decryptor = ConfigDefault.Decryptor
		}
	}

	if cfg.Key == "" {
		panic("kanggo: encryptcookie middleware requires key")
	}

	return cfg
}

// New 创建一个新的 encryptcookie 中间件
func New(config ...Config) core.MiddlewareFunc {
	// 使用默认配置
	cfg := configDefault(config...)

	// 返回中间件函数
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 创建一个新的 KangGo 上下文
			ctx := &kanggo.Context{
				Writer:  w,
				Request: r,
			}

			// 如果 Next 返回 true，则跳过此中间件
			if cfg.Next != nil && cfg.Next(ctx) {
				next(w, r)
				return
			}

			// 解密请求 Cookie
			for _, cookie := range r.Cookies() {
				if !isDisabled(cookie.Name, cfg.Except) {
					decryptedValue, err := cfg.Decryptor(cookie.Value, cfg.Key)
					if err != nil {
						// 如果解密失败，删除该 Cookie
						http.SetCookie(w, &http.Cookie{Name: cookie.Name, MaxAge: -1})
					} else {
						// 如果解密成功，设置解密后的值
						r.AddCookie(&http.Cookie{Name: cookie.Name, Value: decryptedValue})
					}
				}
			}

			// 执行下一个处理函数
			next(w, r)

			// 加密响应 Cookie
			for _, cookie := range w.Header()["Set-Cookie"] {
				name := cookie[:len(cookie)-len(cookie)]
				if !isDisabled(name, cfg.Except) {
					encryptedValue, err := cfg.Encryptor(name, cfg.Key)
					if err != nil {
						panic(err)
					}
					// 设置加密后的 Cookie 值
					http.SetCookie(w, &http.Cookie{Name: name, Value: encryptedValue})
				}
			}
		}
	}
}

// EncryptCookie 加密 Cookie 值
func EncryptCookie(value, key string) (string, error) {
	keyDecoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("无法解码密钥: %w", err)
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return "", fmt.Errorf("无法创建 AES 密码: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("无法创建 GCM 模式: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("无法读取随机数: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptCookie 解密 Cookie 值
func DecryptCookie(value, key string) (string, error) {
	keyDecoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("无法解码密钥: %w", err)
	}
	enc, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("无法解码值: %w", err)
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return "", fmt.Errorf("无法创建 AES 密码: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("无法创建 GCM 模式: %w", err)
	}

	nonceSize := gcm.NonceSize()

	if len(enc) < nonceSize {
		return "", errors.New("加密值无效")
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("无法解密密文: %w", err)
	}

	return string(plaintext), nil
}

// GenerateKey 生成加密密钥
func GenerateKey() string {
	const keyLen = 32
	ret := make([]byte, keyLen)

	if _, err := rand.Read(ret); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(ret)
}

// 检查给定的 Cookie 密钥是否已禁用加密
func isDisabled(key string, except []string) bool {
	for _, k := range except {
		if key == k {
			return true
		}
	}

	return false
}
