{


  func main() {
	r := gin.New()

	rl := NewRateLimiter()

	// Apply the rate limiter middleware to all routes.
	r.Use(RateLimiterMiddleware(rl))

	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	r.Run() // Listen and serve on 0.0.0.0:8080
 }




}
