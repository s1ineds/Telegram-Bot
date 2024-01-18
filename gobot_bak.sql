--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1
-- Dumped by pg_dump version 16.1

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: expenses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.expenses (
    exp_id integer NOT NULL,
    category character varying(128),
    item_name character varying(128),
    price integer,
    buy_date date,
    user_id integer
);


ALTER TABLE public.expenses OWNER TO postgres;

--
-- Name: expenses_exp_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.expenses_exp_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.expenses_exp_id_seq OWNER TO postgres;

--
-- Name: expenses_exp_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.expenses_exp_id_seq OWNED BY public.expenses.exp_id;


--
-- Name: income; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income (
    inc_id integer NOT NULL,
    category character varying(128),
    item_name character varying(128),
    price integer,
    buy_date date,
    user_id integer
);


ALTER TABLE public.income OWNER TO postgres;

--
-- Name: income_inc_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.income_inc_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.income_inc_id_seq OWNER TO postgres;

--
-- Name: income_inc_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.income_inc_id_seq OWNED BY public.income.inc_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    user_id integer NOT NULL,
    chat_id bigint,
    first_name character varying(128)
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_user_id_seq OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.users.user_id;


--
-- Name: expenses exp_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expenses ALTER COLUMN exp_id SET DEFAULT nextval('public.expenses_exp_id_seq'::regclass);


--
-- Name: income inc_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income ALTER COLUMN inc_id SET DEFAULT nextval('public.income_inc_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_user_id_seq'::regclass);


--
-- Data for Name: expenses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.expenses (exp_id, category, item_name, price, buy_date, user_id) FROM stdin;
1	Транспорт	Бензин	6000	2023-10-15	1
2	Транспорт	Автобус	160	2023-10-16	1
3	Продукты питания	Хлеб	320	2023-10-20	1
4	Продукты питания	Рис	250	2023-10-25	1
5	Продукты питания	Макароны	1250	2023-10-25	1
6	Продукты питания	Майонез	800	2023-10-25	1
7	Продукты питания	Конфеты Merci	2500	2023-10-25	1
8	Продукты питания	Мука	1500	2023-10-25	1
9	Транспорт	Автобус	160	2023-10-25	1
10	Образование	Курсы рисования	250000	2023-11-01	1
11	Подарки	Цветы	19000	2023-11-15	1
12	Развлечения	Кино	5000	2023-11-29	1
13	Продукты питания	Яблочипсы	250	2023-12-05	1
14	Продукты питания	КокаКола	670	2023-12-05	1
15	Медикаменты	Аскорбинка	320	2023-12-15	1
16	Медикаменты	Презервативы	1150	2023-12-15	1
17	Ком. Услуги	Свет	10000	2023-12-20	1
18	Ком. Услуги	Вода	3000	2023-12-20	1
19	Ком. Услуги	Мусор	550	2023-12-20	1
20	Ком. Услуги	Интернет	9000	2023-12-20	1
\.


--
-- Data for Name: income; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.income (inc_id, category, item_name, price, buy_date, user_id) FROM stdin;
1	Доп. доходы	Подкинули в переходе	15000	2023-11-01	1
2	Доп. доходы	Майнинг	16000	2023-11-05	1
3	Заработная плата	Зарплата	16000	2023-12-10	1
4	Прочее	Продажа оперативки	4000	2023-12-10	1
5	Прочее	Продажа материнки	20000	2023-12-15	1
6	Прочее	Продажа инструмента	32000	2023-12-20	1
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (user_id, chat_id, first_name) FROM stdin;
1	5495193096	Denis
\.


--
-- Name: expenses_exp_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.expenses_exp_id_seq', 20, true);


--
-- Name: income_inc_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.income_inc_id_seq', 6, true);


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_user_id_seq', 1, true);


--
-- Name: expenses expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_pkey PRIMARY KEY (exp_id);


--
-- Name: income income_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_pkey PRIMARY KEY (inc_id);


--
-- Name: users users_chat_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_chat_id_key UNIQUE (chat_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: expenses expenses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: income income_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- PostgreSQL database dump complete
--

