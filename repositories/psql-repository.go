package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cheebz/arb"
	"github.com/cheebz/go-pub/cache"
	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/models"
	"github.com/cheebz/go-pub/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PSQLRepository struct {
	conf  config.Configuration
	cache cache.Cache
	db    *pgxpool.Pool
}

func NewPSQLRepository(_conf config.Configuration, _cache cache.Cache) Repository {
	return &PSQLRepository{
		conf:  _conf,
		cache: _cache,
		db:    connectDb(_conf.Db),
	}
}

func connectDb(s config.DataSource) *pgxpool.Pool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.User, s.Password, s.Dbname)
	db, err := pgxpool.Connect(context.Background(), psqlInfo)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	log.Printf("Connected to %s as %s\n", s.Dbname, s.User)
	return db
}

func (r *PSQLRepository) Close() {
	r.db.Close()
}

func (r *PSQLRepository) QueryUserByName(name string) (models.User, error) {
	sql := `SELECT * FROM users
	WHERE name = $1
	LIMIT 1`

	var user models.User
	err := r.db.QueryRow(context.Background(), sql, name).Scan(
		&user.ID,
		&user.Name,
		&user.Discoverable,
		&user.IRI,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *PSQLRepository) CheckUser(name string) error {
	sql := `SELECT 1 from users
	WHERE name = $1`

	var result int
	_ = r.db.QueryRow(context.Background(), sql, name).Scan(&result)
	if result != 1 {
		return errors.New("user does not exist")
	}
	return nil
}

func (r *PSQLRepository) CreateUser(name string) (string, error) {
	sql := `INSERT INTO users (name, discoverable, iri)
	VALUES ($1, true, $2)`

	iri := fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name)
	_, err := r.db.Exec(context.Background(), sql, name, iri)
	if err != nil {
		return iri, err
	}
	return iri, nil
}

func (r *PSQLRepository) QueryInboxTotalItemsByUserName(name string) (int, error) {
	var count int
	_, err := r.cache.Get(fmt.Sprintf("inbox-totalItems-%s", name), &count)
	if err == nil {
		return count, nil
	}
	log.Println(fmt.Sprintf("no cached %s", fmt.Sprintf("inbox-totalItems-%s", name)))

	sql := `SELECT COUNT(act.*)
	FROM activities as act
	JOIN activities_to AS act_to ON act_to.activity_id = act.id
	WHERE act_to.iri = $1`

	err = r.db.QueryRow(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
	).Scan(&count)
	if err != nil {
		return count, err
	}

	err = r.cache.Set(fmt.Sprintf("inbox-totalItems-%s", name), count)
	if err != nil {
		log.Println(fmt.Sprintf("error setting cache %s", fmt.Sprintf("inbox-totalItems-%s", name)))
	}

	return count, nil
}

func (r *PSQLRepository) QueryInboxByUserName(name string, pageNum int) ([]models.Activity, error) {
	var activities []models.Activity
	r.cache.Get(fmt.Sprintf("inbox-%s-%d", name, pageNum), &activities)
	if activities != nil {
		return activities, nil
	}
	log.Println(fmt.Sprintf("no cached %s", fmt.Sprintf("inbox-%s-%d", name, pageNum)))

	sql := `SELECT act.*
	FROM activities as act
	JOIN activities_to AS act_to ON act_to.activity_id = act.id
	WHERE act_to.iri = $1
	ORDER BY id DESC
	OFFSET $2
	LIMIT $3`

	rows, err := r.db.Query(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
		pageNum*r.conf.PageLength,
		r.conf.PageLength+1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var activity_id int
		var object_id int
		activity := models.NewActivity()
		err = rows.Scan(
			&activity_id,
			&activity.Type,
			&activity.Actor,
			&object_id,
			&activity.Id,
		)
		if err != nil {
			return activities, err
		}
		object_iri, err := r.queryObjectIRIById(object_id)
		if err != nil {
			return activities, err
		}
		object, err := r.queryObjectByIRI(object_iri)
		if err != nil {
			activity.ChildObject = object_iri

		} else {
			activity.ChildObject = object
		}
		activity.To, err = r.queryToByActivityId(activity_id)
		if err != nil {
			return activities, err
		}
		activities = append(activities, activity)
	}
	err = rows.Err()
	if err != nil {
		return activities, err
	}

	err = r.cache.Set(fmt.Sprintf("inbox-%s-%d", name, pageNum), activities)
	if err != nil {
		log.Println(fmt.Sprintf("error setting cache %s", fmt.Sprintf("inbox-%s-%d", name, pageNum)))
	}

	return activities, nil
}

func (r *PSQLRepository) QueryOutboxTotalItemsByUserName(name string) (int, error) {
	var count int
	_, err := r.cache.Get(fmt.Sprintf("outbox-totalItems-%s", name), &count)
	if err == nil {
		return count, nil
	}
	log.Println(fmt.Sprintf("no cached %s", fmt.Sprintf("outbox-totalItems-%s", name)))

	sql := `SELECT COUNT(*) FROM activities
	WHERE actor = $1`

	err = r.db.QueryRow(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
	).Scan(&count)
	if err != nil {
		return count, err
	}

	err = r.cache.Set(fmt.Sprintf("outbox-totalItems-%s", name), count)
	if err != nil {
		log.Println(fmt.Sprintf("error setting cache %s", fmt.Sprintf("outbox-totalItems-%s", name)))
	}

	return count, nil
}

func (r *PSQLRepository) QueryOutboxByUserName(name string, pageNum int) ([]models.Activity, error) {
	var activities []models.Activity
	_, err := r.cache.Get(fmt.Sprintf("outbox-%s-%d", name, pageNum), &activities)
	if err == nil {
		return activities, nil
	}
	log.Println(fmt.Sprintf("no cached %s", fmt.Sprintf("outbox-%s-%d", name, pageNum)))

	sql := `SELECT *
	FROM activities
	WHERE actor = $1
	ORDER BY id DESC
	OFFSET $2
	LIMIT $3`

	rows, err := r.db.Query(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
		pageNum*r.conf.PageLength,
		r.conf.PageLength+1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var activity_id int
		var object_id int
		activity := models.NewActivity()
		err = rows.Scan(
			&activity_id,
			&activity.Type,
			&activity.Actor,
			&object_id,
			&activity.Id,
		)
		if err != nil {
			return activities, err
		}
		object_iri, err := r.queryObjectIRIById(object_id)
		if err != nil {
			return activities, err
		}
		object, err := r.queryObjectByIRI(object_iri)
		if err != nil {
			activity.ChildObject = object_iri

		} else {
			activity.ChildObject = object
		}
		activity.To, err = r.queryToByActivityId(activity_id)
		if err != nil {
			return activities, err
		}
		activities = append(activities, activity)
	}
	err = rows.Err()
	if err != nil {
		return activities, err
	}

	err = r.cache.Set(fmt.Sprintf("outbox-%s-%d", name, pageNum), activities)
	if err != nil {
		log.Println(fmt.Sprintf("error setting cache %s", fmt.Sprintf("outbox-%s-%d", name, pageNum)))
	}

	return activities, nil
}

func (r *PSQLRepository) queryObjectIRIById(object_id int) (string, error) {
	sql := `SELECT iri
	FROM objects WHERE id = $1;`
	var iri string
	err := r.db.QueryRow(context.Background(), sql, object_id).Scan(
		&iri,
	)
	if err != nil {
		return iri, err
	}
	return iri, nil
}

func (r *PSQLRepository) queryObjectByIRI(iri string) (models.Object, error) {
	sql := `SELECT type, iri, content, attributed_to, in_reply_to, name
	FROM objects WHERE iri = $1;`
	object := models.NewObject()
	err := r.db.QueryRow(context.Background(), sql, iri).Scan(
		&object.Type,
		&object.Id,
		&object.Content,
		&object.AttributedTo,
		&object.InReplyTo,
		&object.Name,
	)
	if err != nil {
		return object, err
	}

	sql = `SELECT id, type, href, media_type
	FROM object_files
	WHERE object_id = (SELECT id FROM object WHERE iri = $1)`
	rows, err := r.db.Query(context.Background(), sql, iri)
	if err != nil {
		return object, err
	}
	defer rows.Close()
	var links []models.Link
	for rows.Next() {
		link := models.NewLink()
		var link_id int
		err = rows.Scan(
			&link_id,
			&link.Type,
			&link.Href,
			&link.MediaType,
		)
		if err != nil {
			return object, err
		}
		link.Id = fmt.Sprintf("%s://%s/%s/%d", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Links, link_id)
		links = append(links, link)
	}
	if links != nil {
		object.Url = links
	}
	err = rows.Err()
	if err != nil {
		return object, err
	}

	return object, nil
}

func (r *PSQLRepository) queryLinksByObjectID(id int) ([]models.Link, error) {
	var links []models.Link
	sql := `SELECT id, type, href, media_type
	FROM object_files
	WHERE object_id = $1`
	rows, err := r.db.Query(context.Background(), sql, id)
	if err != nil {
		return links, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return links, err
	}
	for rows.Next() {
		link := models.NewLink()
		var link_id int
		err = rows.Scan(
			&link_id,
			&link.Type,
			&link.Href,
			&link.MediaType,
		)
		if err != nil {
			return links, err
		}
		link.Id = fmt.Sprintf("%s://%s/%s/%d", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Links, link_id)
		links = append(links, link)
	}
	return links, nil
}

func (r *PSQLRepository) queryToByActivityId(activity_id int) ([]string, error) {
	sql := `SELECT iri
	FROM activities_to
	WHERE activity_id = $1`

	rows, err := r.db.Query(context.Background(), sql, activity_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tos []string
	for rows.Next() {
		var to string
		err = rows.Scan(
			&to,
		)
		if err != nil {
			return tos, err
		}
		tos = append(tos, to)
	}
	err = rows.Err()
	if err != nil {
		return tos, err
	}
	return tos, nil
}

func (r *PSQLRepository) QueryFollowingTotalItemsByUserName(name string) (int, error) {
	sql := `SELECT COUNT(*)
	FROM activities
	WHERE type = 'Follow'
	AND iri NOT IN (
		SELECT obj.iri FROM activities AS act
		JOIN objects AS obj ON obj.id = act.object_id
		WHERE act.type = 'Undo'
	)
	AND actor = $1`

	var count int
	err := r.db.QueryRow(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
	).Scan(
		&count,
	)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *PSQLRepository) QueryFollowingByUserName(name string, pageNum int) ([]string, error) {
	sql := `SELECT obj.iri
	FROM activities AS act
	JOIN objects AS obj ON obj.id = act.object_id
	WHERE act.type = 'Follow'
	AND act.iri NOT IN (
		SELECT obj.iri FROM activities AS act
		JOIN objects AS obj ON obj.id = act.object_id
		WHERE act.type = 'Undo'
	)
	AND act.actor = $1
	ORDER BY act.id DESC
	OFFSET $2
	LIMIT $3`

	rows, err := r.db.Query(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
		pageNum*r.conf.PageLength,
		r.conf.PageLength+1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var actors []string
	for rows.Next() {
		var actor string
		err = rows.Scan(
			&actor,
		)
		if err != nil {
			return actors, err
		}
		actors = append(actors, actor)
	}
	err = rows.Err()
	if err != nil {
		return actors, err
	}
	return actors, nil
}

func (r *PSQLRepository) QueryFollowersTotalItemsByUserName(name string) (int, error) {
	sql := `SELECT COUNT(*)
	FROM activities AS act
	JOIN objects AS obj ON obj.id = act.object_id
	WHERE act.type = 'Follow'
	AND act.iri NOT IN (
		SELECT obj.iri FROM activities AS act
		JOIN objects AS obj ON obj.id = act.object_id
		WHERE act.type = 'Undo'
	)
	AND obj.iri = $1`

	var count int
	err := r.db.QueryRow(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
	).Scan(
		&count,
	)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *PSQLRepository) QueryFollowersByUserName(name string, pageNum int) ([]string, error) {
	sql := `SELECT act.actor
	FROM activities AS act
	JOIN objects AS obj ON obj.id = act.object_id
	WHERE act.type = 'Follow'
	AND act.iri NOT IN (
		SELECT obj.iri FROM activities AS act
		JOIN objects AS obj ON obj.id = act.object_id
		WHERE act.type = 'Undo'
	)
	AND obj.iri = $1
	ORDER BY act.id DESC
	OFFSET $2
	LIMIT $3`

	rows, err := r.db.Query(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
		pageNum*r.conf.PageLength,
		r.conf.PageLength+1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var actors []string
	for rows.Next() {
		var actor string
		err = rows.Scan(
			&actor,
		)
		if err != nil {
			return actors, err
		}
		actors = append(actors, actor)
	}
	err = rows.Err()
	if err != nil {
		return actors, err
	}
	return actors, nil
}

func (r *PSQLRepository) QueryLikedTotalItemsByUserName(name string) (int, error) {
	sql := `SELECT COUNT(*)
	FROM activities
	WHERE type = 'Like'
	AND iri NOT IN (
		SELECT obj.iri FROM activities AS act
		JOIN objects AS obj ON obj.id = act.object_id
		WHERE act.type = 'Undo'
	)
	AND actor = $1`

	var count int
	err := r.db.QueryRow(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
	).Scan(
		&count,
	)
	if err != nil {
		return count, err
	}
	return count, nil
}

// func (r *PSQLRepository) QueryLikedByUserName(name string, pageNum int) ([]models.Object, error) {
func (r *PSQLRepository) QueryLikedByUserName(name string, pageNum int) ([]string, error) {
	// sql := `SELECT obj.type, obj.iri, obj.content, obj.attributed_to, obj.in_reply_to
	sql := `SELECT obj.iri
	FROM objects AS obj
	JOIN activities AS act ON act.object_id = obj.id
	WHERE act.type = 'Like'
	AND act.iri NOT IN (
		SELECT obj.iri FROM activities AS act
		JOIN objects AS obj ON obj.id = act.object_id
		WHERE act.type = 'Undo'
	)
	AND act.actor = $1
	ORDER BY act.id DESC
	OFFSET $2
	LIMIT $3`

	rows, err := r.db.Query(context.Background(), sql,
		fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
		pageNum*r.conf.PageLength,
		r.conf.PageLength+1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// var objects []models.Object
	var objects []string
	for rows.Next() {
		// object := models.NewObject()
		var object string
		// err = rows.Scan(
		// 	&object.Type,
		// 	&object.Id,
		// 	&object.Content,
		// 	&object.AttributedTo,
		// 	&object.InReplyTo,
		// )
		err = rows.Scan(
			&object,
		)
		if err != nil {
			return objects, err
		}
		objects = append(objects, object)
	}
	err = rows.Err()
	if err != nil {
		return objects, err
	}
	return objects, nil
}

func (r *PSQLRepository) QueryActivity(id int) (models.Activity, error) {
	activity := models.NewActivity()
	_, err := r.cache.Get(fmt.Sprintf("activity-%d", id), &activity)
	if err == nil {
		return activity, nil
	}
	log.Println(fmt.Sprintf("no cached %s", fmt.Sprintf("activity-%d", id)))

	sql := `SELECT * FROM activities
	WHERE id = $1
	LIMIT 1`
	var activity_id int
	var object_id int
	err = r.db.QueryRow(context.Background(), sql, id).Scan(
		&activity_id,
		&activity.Type,
		&activity.Actor,
		&object_id,
		&activity.Id,
	)
	if err != nil {
		return activity, err
	}
	object_iri, err := r.queryObjectIRIById(object_id)
	if err != nil {
		return activity, err
	}
	object, err := r.queryObjectByIRI(object_iri)
	if err != nil {
		activity.ChildObject = object_iri

	} else {
		activity.ChildObject = object
	}
	activity.To, err = r.queryToByActivityId(activity_id)
	if err != nil {
		return activity, err
	}

	err = r.cache.Set(fmt.Sprintf("activity-%d", id), activity)
	if err != nil {
		log.Println(fmt.Sprintf("error setting cache %s", fmt.Sprintf("activity-%d", id)))
	}

	return activity, nil
}

func (r *PSQLRepository) QueryObject(id int) (models.Object, error) {
	object := models.NewObject()
	_, err := r.cache.Get(fmt.Sprintf("object-%d", id), &object)
	if err == nil {
		return object, nil
	}
	log.Println(fmt.Sprintf("no cached %s", fmt.Sprintf("activity-%d", id)))

	sql := `SELECT type, iri, content, attributed_to, in_reply_to, name
	FROM objects WHERE id = $1;`
	err = r.db.QueryRow(context.Background(), sql, id).Scan(
		&object.Type,
		&object.Id,
		&object.Content,
		&object.AttributedTo,
		&object.InReplyTo,
		&object.Name,
	)
	if err != nil {
		return object, err
	}
	links, err := r.queryLinksByObjectID(id)
	if err != nil {
		return object, err
	}
	if links != nil {
		object.Url = links
	}
	err = r.cache.Set(fmt.Sprintf("object-%d", id), object)
	if err != nil {
		log.Println(fmt.Sprintf("error setting cache %s", fmt.Sprintf("object-%d", id)))
	}

	return object, nil
}

// Create a new inbox Activity with basic details
func (r *PSQLRepository) CreateInboxActivity(activityArb arb.Arb, objectArb arb.Arb, actor string, name string) (arb.Arb, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return activityArb, err
	}
	objectIRI, _ := objectArb.GetString("id")
	object_id, err := r.queryObjectID(objectIRI)
	if err != nil {
		sql := `INSERT INTO objects (iri, type, content, attributed_to, in_reply_to) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`
		err = tx.QueryRow(ctx, sql,
			objectArb["id"],
			objectArb["type"],
			objectArb["content"],
			objectArb["attributedTo"],
			objectArb["inReplyTo"],
		).Scan(&object_id)
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
	}
	activityIRI, _ := activityArb.GetString("id")
	activity_id, err := r.queryActivityID(activityIRI)
	if err != nil {
		sql := `INSERT INTO activities (type, actor, object_id, iri)
		VALUES ($1, $2, $3, $4) RETURNING id;`
		err = tx.QueryRow(ctx, sql, activityArb["type"], actor, object_id, activityArb["id"]).Scan(&activity_id)
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
	}
	iri := fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name)
	if !r.ActivityToExists(activityIRI, iri) {
		sql := `INSERT INTO activities_to (activity_id, iri) VALUES ($1, $2);`
		_, err = tx.Exec(ctx, sql, activity_id, iri)
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	err = r.cache.Del(fmt.Sprintf("inbox-%s-*", name), fmt.Sprintf("inbox-totalItems-%s", name))
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cache %s and %s", fmt.Sprintf("inbox-%s-*", name), fmt.Sprintf("inbox-totalItems-%s", name)))
	}
	// TODO: Invalidate other cache items based on activityArb["type"]
	return activityArb, nil
}

// Create a new inbox Activity with basic details
func (r *PSQLRepository) CreateInboxReferenceActivity(activityArb arb.Arb, object string, actor string, name string) (arb.Arb, error) {
	// activityIRI, _ := activityArb.GetString("id")
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return activityArb, err
	}
	object_id, err := r.queryObjectID(object)
	if err != nil {
		sql := `INSERT INTO objects (iri) 
		VALUES ($1) RETURNING id;`
		err = tx.QueryRow(ctx, sql, object).Scan(&object_id)
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
	}
	activityIRI, _ := activityArb.GetString("id")
	activity_id, err := r.queryActivityID(activityIRI)
	if err != nil {
		sql := `INSERT INTO activities (type, actor, object_id, iri)
		VALUES ($1, $2, $3, $4) RETURNING id;`
		err = tx.QueryRow(ctx, sql, activityArb["type"], actor, object_id, activityArb["id"]).Scan(&activity_id)
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
	}
	iri := fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name)
	if !r.ActivityToExists(activityIRI, iri) {
		sql := `INSERT INTO activities_to (activity_id, iri) VALUES ($1,$2);`
		_, err = tx.Exec(ctx, sql, activity_id, iri)
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
	}
	tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	err = r.cache.Del(fmt.Sprintf("inbox-%s-*", name), fmt.Sprintf("inbox-totalItems-%s", name))
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cache %s and %s", fmt.Sprintf("inbox-%s-*", name), fmt.Sprintf("inbox-totalItems-%s", name)))
	}
	// TODO: Invalidate other cache items based on activityArb["type"]
	return activityArb, nil
}

func (r *PSQLRepository) queryObjectID(iri string) (int, error) {
	sql := `SELECT id
	FROM objects WHERE iri = $1;`
	var object_id int
	err := r.db.QueryRow(context.Background(), sql, iri).Scan(&object_id)
	if err != nil {
		return object_id, err
	}
	return object_id, nil
}

func (r *PSQLRepository) queryActivityID(iri string) (int, error) {
	sql := `SELECT id
	FROM activities WHERE iri = $1;`
	var activity_id int
	err := r.db.QueryRow(context.Background(), sql, iri).Scan(&activity_id)
	if err != nil {
		return activity_id, err
	}
	return activity_id, nil
}

func (r *PSQLRepository) CreateOutboxActivity(activityArb arb.Arb, objectArb arb.Arb, name string) (arb.Arb, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return activityArb, err
	}
	// TODO: Code here to prevent duplicate objects???
	sql := `INSERT INTO objects (type, content, attributed_to, in_reply_to, name) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	var object_id int
	err = tx.QueryRow(ctx, sql,
		objectArb["type"],
		objectArb["content"],
		objectArb["attributedTo"],
		objectArb["inReplyTo"],
		objectArb["name"],
	).Scan(&object_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	objectArb["id"] = fmt.Sprintf("%s://%s/%s/%d", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Objects, object_id)
	sql = `UPDATE objects
	SET iri = $1
	WHERE id = $2;`
	_, err = tx.Exec(ctx, sql, objectArb["id"], object_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	if objectArb.Exists("url") {
		fileArbs, err := objectArb.GetArbArray("url")
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
		var urls []arb.Arb
		for _, fileArb := range fileArbs {
			sql = `INSERT INTO object_files (object_id, created, name, uuid, type, href, media_type) 
			VALUES ($1, CURRENT_TIMESTAMP, $2, $3, $4, $5, $6);`
			// TODO: return id to populate file iri
			_, err = tx.Exec(ctx, sql,
				object_id,
				fileArb["name"],
				fileArb["uuid"],
				fileArb["type"],
				fileArb["href"],
				fileArb["mediaType"],
			)
			if err != nil {
				tx.Rollback(ctx)
				return activityArb, err
			}
			urls = append(urls, fileArb)
		}
		objectArb["url"] = urls
	}
	sql = `INSERT INTO activities (type, actor, object_id)
	VALUES ($1, $2, $3) RETURNING id;`
	var activity_id int
	err = tx.QueryRow(ctx, sql, activityArb["type"], activityArb["actor"], object_id).Scan(&activity_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	activityArb["id"] = fmt.Sprintf("%s://%s/%s/%d", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Activities, activity_id)
	sql = `UPDATE activities
	SET iri = $1
	WHERE id = $2;`
	_, err = tx.Exec(ctx, sql, activityArb["id"], activity_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	// // Insert to records (need to do similar for bto, cc, bcc, and audience)
	// valueStrings, valueArgs := createRecipientsInsert(activity_id, activityArb.To)
	// sql = fmt.Sprintf("INSERT INTO activities_to (activity_id, iri) VALUES %s", strings.Join(valueStrings, ","))
	// _, err = tx.Exec(ctx, sql, valueArgs...)
	// if err != nil {
	// 	tx.Rollback(ctx)
	// 	return activityArb, err
	// }
	tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	err = r.cache.Del(fmt.Sprintf("outbox-%s-*", name), fmt.Sprintf("outbox-totalItems-%s", name))
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cache %s and %s", fmt.Sprintf("outbox-%s-*", name), fmt.Sprintf("outbox-totalItems-%s", name)))
	}
	// TODO: Invalidate other cache items based on activityArb["type"]
	return activityArb, nil
}

func (r *PSQLRepository) CreateOutboxReferenceActivity(activityArb arb.Arb, name string) (arb.Arb, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return activityArb, err
	}
	objectIRI, _ := activityArb.GetString("object")
	object_id, err := r.queryObjectID(objectIRI)
	if err != nil {
		sql := `INSERT INTO objects (iri) 
		VALUES ($1) RETURNING id;`
		err = tx.QueryRow(ctx, sql,
			activityArb["object"],
		).Scan(&object_id)
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
	}
	sql := `INSERT INTO activities (type, actor, object_id)
	VALUES ($1, $2, $3) RETURNING id;`
	var activity_id int
	err = tx.QueryRow(ctx, sql, activityArb["type"], activityArb["actor"], object_id).Scan(&activity_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	activityArb["id"] = fmt.Sprintf("%s://%s/%s/%d", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Activities, activity_id)
	sql = `UPDATE activities
	SET iri = $1
	WHERE id = $2;`
	_, err = tx.Exec(ctx, sql, activityArb["id"], activity_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	// Insert to records (need to do similar for bto, cc, bcc, and audience)
	// valueStrings, valueArgs := createRecipientsInsert(activity_id, activityArb.To)
	// sql = fmt.Sprintf("INSERT INTO activities_to (activity_id, iri) VALUES %s", strings.Join(valueStrings, ","))
	// _, err = tx.Exec(ctx, sql, valueArgs...)
	// if err != nil {
	// 	tx.Rollback(ctx)
	// 	return activityArb, err
	// }
	tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	err = r.cache.Del(fmt.Sprintf("outbox-%s-*", name), fmt.Sprintf("outbox-totalItems-%s", name))
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cache %s and %s", fmt.Sprintf("outbox-%s-*", name), fmt.Sprintf("outbox-totalItems-%s", name)))
	}
	// TODO: Invalidate other cache items based on activityArb["type"]
	return activityArb, nil
}

func (r *PSQLRepository) ActivityToExists(activityIRI string, recipientIRI string) bool {
	sql := `SELECT 1 from activities_to
	WHERE activity_id = (SELECT id from activities WHERE iri = $1 LIMIT 1)
	AND iri = $2`
	var result int
	_ = r.db.QueryRow(context.Background(), sql, activityIRI, recipientIRI).Scan(&result)
	if result != 1 {
		return false
	}
	log.Println(fmt.Sprintf("%s to %s exists", activityIRI, recipientIRI))
	return true
}

func (r *PSQLRepository) AddActivityTo(activityIRI string, recipient string) error {
	sql := `INSERT INTO activities_to (activity_id, iri) 
	VALUES (
		(SELECT id FROM activities WHERE iri = $1 LIMIT 1),
		$2
	);`
	_, err := r.db.Exec(context.Background(), sql, activityIRI, recipient)
	if err != nil {
		return err
	}
	// If local clear inbox cache of user
	if utils.IsFromHost(recipient, r.conf.ServerName) {
		name := strings.TrimPrefix(recipient, fmt.Sprintf("%s://%s/%s/", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users))
		return r.cache.Del(fmt.Sprintf("inbox-%s-*", name), fmt.Sprintf("inbox-totalItems-%s", name))
	}
	return nil
}

func (r *PSQLRepository) DeleteActivity(activityArb arb.Arb, name string) (arb.Arb, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return activityArb, err
	}
	objectIRI, _ := activityArb.GetString("object")
	object_id, err := r.queryObjectID(objectIRI)
	if err != nil {
		sql := `INSERT INTO objects (iri) 
		VALUES ($1) RETURNING id;`
		err = tx.QueryRow(ctx, sql,
			activityArb["object"],
		).Scan(&object_id)
		if err != nil {
			tx.Rollback(ctx)
			return activityArb, err
		}
	}
	sql := `INSERT INTO activities (type, actor, object_id)
	VALUES ($1, $2, $3) RETURNING id;`
	var activity_id int
	err = tx.QueryRow(ctx, sql, activityArb["type"], activityArb["actor"], object_id).Scan(&activity_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	activityArb["id"] = fmt.Sprintf("%s://%s/%s/%d", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Activities, activity_id)
	sql = `UPDATE activities
	SET iri = $1
	WHERE id = $2;`
	_, err = tx.Exec(ctx, sql, activityArb["id"], activity_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	sql = `UPDATE objects
	SET type = 'Tombstone',
	content = NULL,
	name = NULL
	WHERE id = $1;`
	_, err = tx.Exec(ctx, sql, object_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	sql = `DELETE FROM object_files
	WHERE object_id = $1;`
	_, err = tx.Exec(ctx, sql, object_id)
	if err != nil {
		tx.Rollback(ctx)
		return activityArb, err
	}
	tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	r.deleteObjectCacheInvalidation(object_id)
	err = r.cache.Del(fmt.Sprintf("outbox-%s-*", name), fmt.Sprintf("outbox-totalItems-%s", name))
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cache %s and %s", fmt.Sprintf("outbox-%s-*", name), fmt.Sprintf("outbox-totalItems-%s", name)))
	}
	// TODO: Invalidate other cache items based on activityArb["type"]
	return activityArb, nil
}

func (r *PSQLRepository) deleteObjectCacheInvalidation(object_id int) {
	// Invalidate cached Object
	err := r.cache.Del(fmt.Sprintf("object-%d", object_id))
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cache %s", fmt.Sprintf("object-%d", object_id)))
	}
	// Invalidate cached Activities referencing this Object
	activity_ids, err := r.getActivityIDsWithObjectID(object_id)
	if err != nil {
		log.Println(fmt.Sprintf("error getting activities referencing object %d: %v", object_id, err))
	}
	var activity_keys []string
	for i := range activity_ids {
		activity_keys = append(activity_keys, fmt.Sprintf("activity-%d", activity_ids[i]))
	}
	err = r.cache.Del(activity_keys...)
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cached activities referencing object %d: %v", object_id, err))
	}
	// Invalidate cached actor inboxes with Activities referencing this Object
	actors_receiving, err := r.getActorsReceivingObjectID(object_id)
	if err != nil {
		log.Println(fmt.Sprintf("error getting actors receiving object %d: %v", object_id, err))
	}
	var actors_receiving_keys []string
	for i := range actors_receiving {
		actors_receiving_keys = append(actors_receiving_keys, fmt.Sprintf("inbox-%s-*", actors_receiving[i]), fmt.Sprintf("inbox-totalItems-%s", actors_receiving[i]))
	}
	err = r.cache.Del(actors_receiving_keys...)
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cached actor inboxes referencing object %d: %v", object_id, err))
	}
	// Invalidate cached actor inboxes with Activities referencing this Object
	actors_sending, err := r.getActorsSendingObjectID(object_id)
	if err != nil {
		log.Println(fmt.Sprintf("error getting actors sending object %d: %v", object_id, err))
	}
	var actors_sending_keys []string
	for i := range actors_receiving {
		actors_sending_keys = append(actors_sending_keys, fmt.Sprintf("outbox-%s-*", actors_sending[i]), fmt.Sprintf("outbox-totalItems-%s", actors_sending[i]))
	}
	err = r.cache.Del(actors_sending_keys...)
	if err != nil {
		log.Println(fmt.Sprintf("error deleting cached actor outboxes referencing object %d: %v", object_id, err))
	}
}

func (r *PSQLRepository) getActivityIDsWithObjectID(object_id int) ([]int, error) {
	sql := `SELECT id
	FROM activities
	WHERE object_id = $1`

	rows, err := r.db.Query(context.Background(), sql, object_id)
	if err != nil {
		return nil, err
	}
	var activity_ids []int
	defer rows.Close()
	for rows.Next() {
		var activity_id int
		err = rows.Scan(&activity_id)
		if err != nil {
			return activity_ids, err
		}
		activity_ids = append(activity_ids, activity_id)
	}
	err = rows.Err()
	if err != nil {
		return activity_ids, err
	}
	return activity_ids, nil
}

func (r *PSQLRepository) getActorsReceivingObjectID(object_id int) ([]string, error) {
	sql := `SELECT DISTINCT usr.name
	FROM activities_to AS act_to
	JOIN activities AS act ON act.id = act_to.activity_id
	JOIN objects AS obj ON obj.id = act.object_id
	JOIN users AS usr ON usr.iri = act_to.iri
	WHERE obj.id = $1`

	rows, err := r.db.Query(context.Background(), sql, object_id)
	if err != nil {
		return nil, err
	}
	var actors []string
	defer rows.Close()
	for rows.Next() {
		var actor string
		err = rows.Scan(&actor)
		if err != nil {
			return actors, err
		}
		actors = append(actors, actor)
	}
	err = rows.Err()
	if err != nil {
		return actors, err
	}
	return actors, nil
}

func (r *PSQLRepository) getActorsSendingObjectID(object_id int) ([]string, error) {
	sql := `SELECT DISTINCT usr.name
	FROM activities AS act
	JOIN objects AS obj ON obj.id = act.object_id
	JOIN users AS usr ON usr.iri = act.actor
	WHERE obj.id = $1`

	rows, err := r.db.Query(context.Background(), sql, object_id)
	if err != nil {
		return nil, err
	}
	var actors []string
	defer rows.Close()
	for rows.Next() {
		var actor string
		err = rows.Scan(&actor)
		if err != nil {
			return actors, err
		}
		actors = append(actors, actor)
	}
	err = rows.Err()
	if err != nil {
		return actors, err
	}
	return actors, nil
}

func (r *PSQLRepository) GetObjectFilesByIRI(objectIRI string) ([]string, error) {
	sql := `SELECT f.href
	FROM object_files AS f
	JOIN objects AS obj ON obj.id = f.object_id
	WHERE obj.iri = $1`

	rows, err := r.db.Query(context.Background(), sql, objectIRI)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var hrefs []string
	err = rows.Err()
	if err != nil {
		return hrefs, err
	}
	for rows.Next() {
		var href string
		err = rows.Scan(&href)
		if err != nil {
			return hrefs, err
		}
		hrefs = append(hrefs, href)
	}
	return hrefs, nil
}

func (r *PSQLRepository) PurgeUnusedFiles() error {
	// TODO: Implement this
	// Need to somehow determine if a file on disk is orphaned
	return nil
}
