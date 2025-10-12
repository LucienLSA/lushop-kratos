package device

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"

	"github.com/go-kratos/kratos/v2/transport"
)

// GetDeviceFingerprint 生成设备唯一标识
func GetDeviceFingerprint(ctx context.Context) string {
	// 从HTTP传输层获取请求信息
	if tr, ok := transport.FromServerContext(ctx); ok {
		// 获取User-Agent
		userAgent := tr.RequestHeader().Get("User-Agent")
		if userAgent == "" {
			userAgent = "unknown"
		}

		// 获取客户端IP
		ip := getClientIP(tr)

		// 获取其他请求头信息
		acceptLanguage := tr.RequestHeader().Get("Accept-Language")
		acceptEncoding := tr.RequestHeader().Get("Accept-Encoding")
		acceptCharset := tr.RequestHeader().Get("Accept-Charset")

		// 组合特征字符串
		featureStr := strings.Join([]string{
			userAgent,
			ip,
			acceptLanguage,
			acceptEncoding,
			acceptCharset,
		}, "|")

		// SHA256哈希生成唯一标识
		h := sha256.New()
		h.Write([]byte(featureStr))
		return hex.EncodeToString(h.Sum(nil))
	}

	// 如果无法获取传输信息，返回默认值
	return "unknown_device"
}

// getClientIP 获取客户端真实IP
func getClientIP(tr transport.Transporter) string {
	// 优先从X-Forwarded-For获取
	if xff := tr.RequestHeader().Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 从X-Real-IP获取
	if xri := tr.RequestHeader().Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 从RemoteAddr获取
	if addr := tr.RequestHeader().Get("Remote-Addr"); addr != "" {
		host, _, err := net.SplitHostPort(addr)
		if err == nil {
			return host
		}
		return addr
	}

	return "unknown"
}
