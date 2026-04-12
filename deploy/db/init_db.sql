-- Initialization Script
-- Combined schema for go-pub (ActivityPub + auth)

DO $$
BEGIN

  -- -----------------------------------------------------------------------
  -- users
  -- Single table for both auth (username/password/uuid) and
  -- ActivityPub profile (discoverable/iri) data.
  -- -----------------------------------------------------------------------
  CREATE TABLE IF NOT EXISTS public.users (
    id           serial      NOT NULL,
    username     text        NOT NULL,
    password     text        NOT NULL DEFAULT '',
    created      timestamptz NOT NULL DEFAULT now(),
    uuid         text        NOT NULL DEFAULT '',
    discoverable bool        NOT NULL DEFAULT true,
    iri          text        NOT NULL DEFAULT '',
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_unique UNIQUE (username)
  );

  -- -----------------------------------------------------------------------
  -- user_refresh  (auth refresh tokens)
  -- -----------------------------------------------------------------------
  CREATE TABLE IF NOT EXISTS public.user_refresh (
    id      serial      NOT NULL,
    user_id int4        NOT NULL,
    jti     text        NOT NULL,
    expires timestamptz NOT NULL,
    CONSTRAINT user_refresh_pkey PRIMARY KEY (id)
  );

  ALTER TABLE public.user_refresh
    DROP CONSTRAINT IF EXISTS fk_user_refresh_user_id;
  ALTER TABLE public.user_refresh
    ADD CONSTRAINT fk_user_refresh_user_id
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

  -- -----------------------------------------------------------------------
  -- objects  (ActivityPub objects)
  -- -----------------------------------------------------------------------
  CREATE TABLE IF NOT EXISTS public.objects (
    id            serial NOT NULL,
    type          text   NULL,
    iri           text   NULL,
    content       text   NULL,
    attributed_to text   NULL,
    in_reply_to   text   NULL,
    name          text   NULL,
    CONSTRAINT objects_pkey PRIMARY KEY (id)
  );

  -- -----------------------------------------------------------------------
  -- object_files  (media attached to objects)
  -- -----------------------------------------------------------------------
  CREATE TABLE IF NOT EXISTS public.object_files (
    id         serial      NOT NULL,
    object_id  int4        NOT NULL,
    created    timestamptz NOT NULL,
    name       text        NOT NULL,
    uuid       text        NOT NULL,
    type       text        NOT NULL,
    href       text        NOT NULL,
    media_type text        NOT NULL,
    CONSTRAINT object_files_pkey PRIMARY KEY (id)
  );

  ALTER TABLE public.object_files
    DROP CONSTRAINT IF EXISTS object_files_object_id_fk;
  ALTER TABLE public.object_files
    ADD CONSTRAINT object_files_object_id_fk
    FOREIGN KEY (object_id) REFERENCES public.objects(id);

  -- -----------------------------------------------------------------------
  -- activities  (ActivityPub activities)
  -- -----------------------------------------------------------------------
  CREATE TABLE IF NOT EXISTS public.activities (
    id        serial NOT NULL,
    type      text   NOT NULL,
    actor     text   NOT NULL,
    object_id int4   NOT NULL,
    iri       text   NULL,
    CONSTRAINT activities_pkey PRIMARY KEY (id)
  );

  ALTER TABLE public.activities
    DROP CONSTRAINT IF EXISTS activities_object_id_fk;
  ALTER TABLE public.activities
    ADD CONSTRAINT activities_object_id_fk
    FOREIGN KEY (object_id) REFERENCES public.objects(id);

  -- -----------------------------------------------------------------------
  -- activities_to  (recipients of activities)
  -- -----------------------------------------------------------------------
  CREATE TABLE IF NOT EXISTS public.activities_to (
    id          serial NOT NULL,
    activity_id int4   NOT NULL,
    iri         text   NOT NULL,
    CONSTRAINT activities_to_pkey PRIMARY KEY (id)
  );

  ALTER TABLE public.activities_to
    DROP CONSTRAINT IF EXISTS activities_to_activity_id_fk;
  ALTER TABLE public.activities_to
    ADD CONSTRAINT activities_to_activity_id_fk
    FOREIGN KEY (activity_id) REFERENCES public.activities(id);

END
$$;
