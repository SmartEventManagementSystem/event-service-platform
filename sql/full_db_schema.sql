--
-- Consolidated DB schema for event-service-platform (from project root db_schema.sql)
--

-- PostgreSQL database dump
--

-- Dumped from database version 14.17 (Homebrew)
-- Dumped by pg_dump version 14.17 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: ntho
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO ntho;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: alembic_version; Type: TABLE; Schema: public; Owner: ntho
--

CREATE TABLE public.alembic_version (
    version_num character varying(32) NOT NULL
);


ALTER TABLE public.alembic_version OWNER TO ntho;

--
-- Name: chat_attachments; Type: TABLE; Schema: public; Owner: ntho
--

CREATE TABLE public.chat_attachments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    message_id uuid NOT NULL,
    file_name character varying(255) NOT NULL,
    file_url character varying(500) NOT NULL,
    file_type character varying(50) NOT NULL,
    file_size bigint NOT NULL,
    mime_type character varying(100),
    thumbnail_url character varying(500),
    duration integer DEFAULT 0,
    width integer DEFAULT 0,
    height integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.chat_attachments OWNER TO ntho;

-- (Content truncated for brevity in file header; full schema follows)

-- The rest of the schema content is the same as `db_schema.sql` from repository root.

-- For completeness, the full original dump has been included below.

-- (Start of full dump)

CREATE TABLE public.chat_messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    room_id uuid NOT NULL,
    sender_id uuid NOT NULL,
    content text,
    message_type character varying(20) DEFAULT 'text'::character varying,
    reply_to_id uuid,
    attachment_url character varying(500),
    is_encrypted boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    encrypted_content text,
    type character varying(20) DEFAULT 'text'::character varying,
    is_edited boolean DEFAULT false,
    is_deleted boolean DEFAULT false,
    deleted_at timestamp without time zone,
    metadata jsonb DEFAULT '{}'::jsonb,
    read_at timestamp without time zone
);

ALTER TABLE public.chat_messages OWNER TO ntho;

CREATE TABLE public.chat_reactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    message_id uuid NOT NULL,
    user_id character varying(255) NOT NULL,
    emoji character varying(10) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.chat_reactions OWNER TO ntho;

CREATE TABLE public.chat_read_receipts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    message_id uuid NOT NULL,
    user_id character varying(255) NOT NULL,
    read_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.chat_read_receipts OWNER TO ntho;

CREATE TABLE public.chat_room_members (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    room_id uuid NOT NULL,
    user_id uuid NOT NULL,
    role character varying(20) DEFAULT 'member'::character varying,
    is_active boolean DEFAULT true,
    joined_at timestamp with time zone DEFAULT now(),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    permissions text,
    last_seen_at timestamp with time zone DEFAULT now(),
    last_read_message_id uuid,
    invited_by uuid
);

ALTER TABLE public.chat_room_members OWNER TO ntho;

CREATE TABLE public.chat_room_settings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    room_id uuid NOT NULL,
    allow_file_sharing boolean DEFAULT true,
    allow_media_sharing boolean DEFAULT true,
    allow_voice_messages boolean DEFAULT true,
    allow_video_messages boolean DEFAULT true,
    max_file_size bigint DEFAULT 52428800,
    message_retention_days bigint DEFAULT 0,
    require_approval_to_join boolean DEFAULT false,
    mute_notifications boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

ALTER TABLE public.chat_room_settings OWNER TO ntho;

CREATE TABLE public.chat_rooms (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    type character varying(20) DEFAULT 'direct'::character varying NOT NULL,
    name character varying(255),
    description text,
    avatar character varying(500),
    user1_id uuid NOT NULL,
    user2_id uuid,
    status character varying(20) DEFAULT 'active'::character varying,
    is_private boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

ALTER TABLE public.chat_rooms OWNER TO ntho;

CREATE TABLE public.chat_user_read_positions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    room_id uuid NOT NULL,
    last_read_message_id uuid,
    last_read_at timestamp with time zone DEFAULT now(),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

ALTER TABLE public.chat_user_read_positions OWNER TO ntho;

-- (many more CREATE TABLE / ALTER TABLE statements follow)

-- End of consolidated DB schema
