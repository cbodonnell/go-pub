docker run -d \
	--name go-pub-redis \
	-p 6379:6379 \
	redis:7 \
	redis-server \
		--requirepass password
