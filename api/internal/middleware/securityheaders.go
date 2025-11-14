package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// HSTS (HTTP Strict Transport Security)
		// Forces HTTPS for 1 year, including subdomains
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Content-Type-Options
		// Prevents MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options
		// Prevents clickjacking attacks
		c.Header("X-Frame-Options", "DENY")

		// X-XSS-Protection
		// Legacy XSS protection (for older browsers)
		c.Header("X-XSS-Protection", "1; mode=block")

		// Content-Security-Policy
		// Restricts resource loading to prevent XSS and data injection
		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: https:; "+
			"font-src 'self' data:; "+
			"connect-src 'self'; "+
			"frame-ancestors 'none'; "+
			"base-uri 'self'; "+
			"form-action 'self'")

		// Referrer-Policy
		// Controls referrer information sent to other sites
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy (formerly Feature-Policy)
		// Disables potentially dangerous browser features
		c.Header("Permissions-Policy",
			"geolocation=(), "+
			"microphone=(), "+
			"camera=(), "+
			"payment=(), "+
			"usb=(), "+
			"magnetometer=(), "+
			"gyroscope=(), "+
			"accelerometer=()")

		// X-Permitted-Cross-Domain-Policies
		// Prevents Adobe Flash and PDF from loading content
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// X-Download-Options
		// Prevents Internet Explorer from executing downloads in site context
		c.Header("X-Download-Options", "noopen")

		// Cache-Control for API responses
		// Prevent caching of sensitive data
		if c.Request.URL.Path != "/health" && c.Request.URL.Path != "/version" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
		}

		c.Next()
	}
}

// SecurityHeadersRelaxed provides relaxed CSP for development
// Use only in development environments
func SecurityHeadersRelaxed() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Same headers as SecurityHeaders() but with relaxed CSP
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN") // Allow same-origin framing for dev
		c.Header("X-XSS-Protection", "1; mode=block")

		// Relaxed CSP for development
		c.Header("Content-Security-Policy",
			"default-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
			"img-src 'self' data: https:; "+
			"connect-src 'self' ws: wss: http: https:")

		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		c.Header("X-Download-Options", "noopen")

		c.Next()
	}
}
