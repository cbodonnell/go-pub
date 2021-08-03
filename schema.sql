-- Drop tables

-- DROP TABLE public.activities_to;

-- DROP TABLE public.activities;

-- DROP TABLE public.objects;

-- Create tables

-- public.objects definition

CREATE TABLE public.objects (
	id serial NOT NULL,
	"type" text NOT NULL,
	iri text NULL,
	"content" text NULL,
	attributed_to text NULL,
	CONSTRAINT objects_pkey PRIMARY KEY (id)
);

-- public.activities definition

CREATE TABLE public.activities (
	id serial NOT NULL,
	"type" text NOT NULL,
	actor text NOT NULL,
	object_id int4 NOT NULL,
	iri text NULL,
	CONSTRAINT outbox_activities_pkey PRIMARY KEY (id)
);

ALTER TABLE public.activities ADD CONSTRAINT activities_object_id_fk FOREIGN KEY (object_id) REFERENCES public.objects(id);

-- public.activities_to definition

CREATE TABLE public.activities_to (
	id int4 NOT NULL,
	activity_id int4 NOT NULL,
	iri text NOT NULL,
	CONSTRAINT outbox_to_pkey PRIMARY KEY (id)
);

ALTER TABLE public.activities_to ADD CONSTRAINT activities_to_activity_id_fk FOREIGN KEY (activity_id) REFERENCES public.activities(id);