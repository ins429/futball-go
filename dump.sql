--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: cards; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE cards (
    id integer NOT NULL,
    name text,
    created_at date,
    updated_at date
);


ALTER TABLE public.cards OWNER TO root;

--
-- Name: cards_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE cards_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cards_id_seq OWNER TO root;

--
-- Name: cards_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE cards_id_seq OWNED BY cards.id;


--
-- Name: players; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE players (
    id integer NOT NULL,
    namealias text,
    name text,
    club text,
    "position" text,
    appearances text,
    goals numeric,
    shots numeric,
    penalties numeric,
    assists numeric,
    crosses numeric,
    offsides numeric,
    savesmade numeric,
    owngoals numeric,
    cleansheets numeric,
    blocks numeric,
    clearances numeric,
    fouls numeric,
    cards numeric,
    dob text,
    height text,
    age numeric,
    weight text,
    "national" text,
    image text
);


ALTER TABLE public.players OWNER TO root;

--
-- Name: players_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE players_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.players_id_seq OWNER TO root;

--
-- Name: players_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE players_id_seq OWNED BY players.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE users (
    id integer NOT NULL,
    username text,
    fb_id text,
    firstname text,
    lastname text,
    encrypted_password text,
    created_at date,
    updated_at date,
    players json[]
);


ALTER TABLE public.users OWNER TO root;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO root;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: wc_players; Type: TABLE; Schema: public; Owner: root; Tablespace: 
--

CREATE TABLE wc_players (
    id integer NOT NULL,
    uid text,
    name text,
    age numeric,
    foot text,
    goals numeric,
    birthcountry text,
    birthcity text,
    penaltygoals numeric,
    birthdate text,
    image text,
    weight numeric,
    assists numeric,
    "national" text,
    "position" text,
    height numeric,
    owngoals numeric,
    club text
);


ALTER TABLE public.wc_players OWNER TO root;

--
-- Name: wc_players_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE wc_players_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.wc_players_id_seq OWNER TO root;

--
-- Name: wc_players_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE wc_players_id_seq OWNED BY wc_players.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY cards ALTER COLUMN id SET DEFAULT nextval('cards_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY players ALTER COLUMN id SET DEFAULT nextval('players_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY wc_players ALTER COLUMN id SET DEFAULT nextval('wc_players_id_seq'::regclass);


--
-- Data for Name: cards; Type: TABLE DATA; Schema: public; Owner: root
--

COPY cards (id, name, created_at, updated_at) FROM stdin;
\.

--
-- Name: players_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('players_id_seq', 4, true);

--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('users_id_seq', 1, true);

--
-- Name: public; Type: ACL; Schema: -; Owner: root
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM root;
GRANT ALL ON SCHEMA public TO root;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

