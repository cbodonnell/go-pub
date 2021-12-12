-- Initialization Script

-- Create tables

-- public.users definition

DO $$
BEGIN

	CREATE TABLE IF NOT EXISTS public.users (
		id serial NOT NULL,
		"name" text NOT NULL,
		discoverable bool NOT NULL,
		iri text NOT NULL,
		CONSTRAINT users_pkey PRIMARY KEY (id)
	);

	-- public.objects definition

	CREATE TABLE IF NOT EXISTS public.objects (
		id serial NOT NULL,
		"type" text NULL,
		iri text NULL,
		"content" text NULL,
		attributed_to text NULL,
		in_reply_to text NULL,
		"name" text NULL,
		CONSTRAINT objects_pkey PRIMARY KEY (id)
	);

	-- public.object_files definition

	CREATE TABLE IF NOT EXISTS public.object_files (
		id serial NOT NULL,
		object_id int4 NOT NULL,
		created timestamptz NOT NULL,
		"name" text NOT NULL,
		uuid text NOT NULL,
		"type" text NOT NULL,
		href text NOT NULL,
		media_type text NOT NULL,
		CONSTRAINT object_files_pkey PRIMARY KEY (id)
	);

	ALTER TABLE public.object_files DROP CONSTRAINT IF EXISTS object_files_object_id_fk;
	ALTER TABLE public.object_files ADD CONSTRAINT object_files_object_id_fk FOREIGN KEY (object_id) REFERENCES public.objects(id);

	-- public.activities definition

	CREATE TABLE IF NOT EXISTS public.activities (
		id serial NOT NULL,
		"type" text NOT NULL,
		actor text NOT NULL,
		object_id int4 NOT NULL,
		iri text NULL,
		CONSTRAINT outbox_activities_pkey PRIMARY KEY (id)
	);

	ALTER TABLE public.activities DROP CONSTRAINT IF EXISTS activities_object_id_fk;
	ALTER TABLE public.activities ADD CONSTRAINT activities_object_id_fk FOREIGN KEY (object_id) REFERENCES public.objects(id);

	-- public.activities_to definition

	CREATE TABLE IF NOT EXISTS public.activities_to (
		id serial NOT NULL,
		activity_id int4 NOT NULL,
		iri text NOT NULL,
		CONSTRAINT outbox_to_pkey PRIMARY KEY (id)
	);

	ALTER TABLE public.activities_to DROP CONSTRAINT IF EXISTS activities_to_activity_id_fk;
	ALTER TABLE public.activities_to ADD CONSTRAINT activities_to_activity_id_fk FOREIGN KEY (activity_id) REFERENCES public.activities(id);

END
$$

