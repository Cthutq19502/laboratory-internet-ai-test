CREATE TABLE IF not exists public.contacts (
     id serial4 NOT NULL,
     name varchar NOT NULL,
     phone varchar NOT NULL,
     email varchar NOT NULL,
     "comment" varchar NOT NULL,
     date_create timestamptz DEFAULT now() NOT NULL,
     tonal varchar NULL,
     CONSTRAINT contacts_pk PRIMARY KEY (id)
);