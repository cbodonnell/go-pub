docker run -d \
	--name go-pub-postgres \
	-p 5432:5432 \
	-e POSTGRES_DB=activitypub \
	-e POSTGRES_USER=activitypub \
	-e POSTGRES_PASSWORD=my-secret-password \
	-v $PWD/pgdata/:/var/lib/postgresql/data/ \
	-v $PWD/deploy/db/:/docker-entrypoint-initdb.d/ \
	postgres:14
