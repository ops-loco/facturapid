package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds common security-enhancing HTTP headers to responses.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevents browsers from trying to MIME-sniff the content-type of a response away from the declared content-type.
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevents the page from being displayed in an iframe, unless it is on the same origin.
		// DENY: The page cannot be displayed in a frame, regardless of the site attempting to do so.
		// SAMEORIGIN: The page can only be displayed in a frame on the same origin as the page itself.
		// ALLOW-FROM uri: The page can only be displayed in a frame on the specified origin. (obsolete in some browsers)
		c.Header("X-Frame-Options", "DENY")

		// HTTP Strict Transport Security (HSTS) informs browsers that the site should only be accessed using HTTPS.
		// max-age is in seconds (e.g., 1 year = 31536000 seconds).
		// includeSubDomains (optional): If this optional parameter is specified, this rule applies to all of the site's subdomains as well.
		// preload (optional): Not part of the spec, but used by browser vendors to include domains in their HSTS preload lists.
		// This header is only effective if the site is actually served over HTTPS.
		if c.Request.TLS != nil { // Only send HSTS if served over HTTPS
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		// Provides XSS protection (though modern browsers have built-in XSS protection, this can enforce it).
		// 0: Disables XSS filtering.
		// 1: Enables XSS filtering (usually default in browsers). If a cross-site scripting attack is detected, the browser will sanitize the page (remove the unsafe parts).
		// 1; mode=block: Enables XSS filtering. Rather than sanitizing the page, the browser will prevent rendering of the page if an attack is detected.
		// 1; report=<reporting-URI>: Enables XSS filtering. If a cross-site scripting attack is detected, the browser will sanitize the page and report the violation.
		c.Header("X-XSS-Protection", "1; mode=block")

		// Content Security Policy (CSP) is a more powerful and flexible alternative/addition to X-XSS-Protection.
		// Example (restrictive): c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none';")
		// Implementing CSP requires careful consideration of your application's resources.
		// For this example, we'll stick to the requested headers.

		c.Next()
	}
}
