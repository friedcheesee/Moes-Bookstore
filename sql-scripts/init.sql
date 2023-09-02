
 \c moe;

 CREATE TABLE users (
     uid SERIAL PRIMARY KEY,
     username VARCHAR(50) NOT NULL,
     email VARCHAR(100) NOT NULL,
     password VARCHAR(100) NOT NULL,
     active BOOLEAN NOT NULL DEFAULT TRUE,
     expiry DATE
 );

 CREATE TABLE books (
     bookid SERIAL PRIMARY KEY,
     book_name VARCHAR(200) NOT NULL,
     author VARCHAR(100) NOT NULL,
     genre VARCHAR(50),
     cost DECIMAL(10, 2) NOT NULL,
     download_url TEXT,
     UNIQUE (book_name, author)
 );

 CREATE TABLE bought_books (
     uid INT REFERENCES users(uid),
     bookid INT REFERENCES books(bookid),
     book_name VARCHAR(200),
     genre VARCHAR(50),
     download_url TEXT,
     review_id INT,
     PRIMARY KEY (uid, bookid)
 );

CREATE TABLE cart(
    uid INT REFERENCES users(uid),
    bookid INT REFERENCES books(bookid),
    book_name VARCHAR(200),
    author VARCHAR(100),
    cost DECIMAL(10, 2),
    PRIMARY KEY (uid, bookid)
);

ALTER TABLE users ADD CONSTRAINT unique_email_constrait UNIQUE (email);

CREATE TABLE reviews (
    reviewid SERIAL PRIMARY KEY,
    uid INT REFERENCES users(uid),
    bookid INT REFERENCES books(bookid),
    review TEXT
);

ALTER TABLE reviews ADD CONSTRAINT unique_user_book_review UNIQUE (uid, bookid);

ALTER TABLE users ADD COLUMN admin BOOLEAN DEFAULT false;


--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4
-- Dumped by pg_dump version 15.4

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
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.books (bookid, book_name, author, genre, cost, download_url) FROM stdin;
1	Sample Book	Jane Author	Fiction	29.99	https://example.com/sample-book.pdf
2	Sample Book2	Juthor	Fantasy	26.99	https://example.com/sample-book.pdf
3	learingey0	Lionello	(no genres listed)	60.00	http://dummyimage.com/129x100.png/dddddd/000000
4	tklemt1	Tandi	Action|Comedy	28.00	http://dummyimage.com/140x100.png/cc0000/ffffff
5	gbranscombe2	Georgie	Comedy	34.00	http://dummyimage.com/103x100.png/ff4444/ffffff
6	mcandlin3	Madelyn	Documentary|Musical	94.00	http://dummyimage.com/166x100.png/5fa2dd/ffffff
7	unobes4	Ursola	Mystery|Thriller	80.00	http://dummyimage.com/116x100.png/dddddd/000000
8	fgrabert5	Felisha	Drama	100.00	http://dummyimage.com/161x100.png/5fa2dd/ffffff
9	cbladder6	Christen	Film-Noir	84.00	http://dummyimage.com/243x100.png/cc0000/ffffff
10	sarlett7	Stacy	Comedy|Drama	96.00	http://dummyimage.com/145x100.png/cc0000/ffffff
11	scardenoza8	Siusan	Thriller	40.00	http://dummyimage.com/196x100.png/5fa2dd/ffffff
12	jbalbeck9	Jerad	Drama|Horror|Mystery	94.00	http://dummyimage.com/145x100.png/cc0000/ffffff
13	sgallandersa	Stephine	Crime|Drama	6.00	http://dummyimage.com/223x100.png/dddddd/000000
14	vmusgroveb	Vaclav	Drama	3.00	http://dummyimage.com/221x100.png/ff4444/ffffff
15	alambertsc	Agatha	Drama	22.00	http://dummyimage.com/164x100.png/5fa2dd/ffffff
16	rsavatierd	Rene	Documentary	18.00	http://dummyimage.com/194x100.png/ff4444/ffffff
17	amargache	Anett	Action|Adventure|Comedy|Thriller	24.00	http://dummyimage.com/164x100.png/cc0000/ffffff
18	abimf	Ahmad	Drama|Fantasy|Musical|Mystery|Sci-Fi	90.00	http://dummyimage.com/226x100.png/ff4444/ffffff
19	jkidwellg	Janeczka	Drama|Thriller	21.00	http://dummyimage.com/215x100.png/5fa2dd/ffffff
20	dommundsenh	Dianne	Comedy	60.00	http://dummyimage.com/212x100.png/5fa2dd/ffffff
21	ghansei	Garland	Drama	36.00	http://dummyimage.com/187x100.png/cc0000/ffffff
22	pfulleylovej	Prentiss	Documentary	55.00	http://dummyimage.com/203x100.png/dddddd/000000
23	ycomettoik	Yoshiko	Animation|Children|Fantasy	93.00	http://dummyimage.com/247x100.png/cc0000/ffffff
24	lgannyl	Lewiss	Adventure|Fantasy|Romance	75.00	http://dummyimage.com/200x100.png/cc0000/ffffff
25	zbalfourm	Zed	Drama|Horror|War	41.00	http://dummyimage.com/138x100.png/ff4444/ffffff
26	pbourgesn	Pet	Horror	16.00	http://dummyimage.com/172x100.png/cc0000/ffffff
27	ballderidgeo	Babs	Action|Adventure|Children|Comedy|Fantasy	53.00	http://dummyimage.com/188x100.png/5fa2dd/ffffff
28	sheinonenp	Shel	War	97.00	http://dummyimage.com/122x100.png/ff4444/ffffff
29	jcarpmileq	Jennilee	Comedy|Romance	63.00	http://dummyimage.com/120x100.png/cc0000/ffffff
30	lgerdesr	Luke	Documentary	56.00	http://dummyimage.com/248x100.png/5fa2dd/ffffff
31	zchokes	Zacharie	Comedy|Fantasy|Musical	54.00	http://dummyimage.com/117x100.png/5fa2dd/ffffff
32	lsigget	Lutero	Comedy	14.00	http://dummyimage.com/224x100.png/5fa2dd/ffffff
33	ahugeu	Adan	Comedy|Western	50.00	http://dummyimage.com/140x100.png/cc0000/ffffff
34	acosserv	Alexandr	Comedy	32.00	http://dummyimage.com/155x100.png/ff4444/ffffff
35	hbutew	Hillel	Drama	26.00	http://dummyimage.com/229x100.png/cc0000/ffffff
36	bbirchnerx	Burr	Comedy	69.00	http://dummyimage.com/223x100.png/ff4444/ffffff
37	tbartosinskiy	Tobiah	Horror|Sci-Fi|Thriller	31.00	http://dummyimage.com/234x100.png/5fa2dd/ffffff
38	vhandscombz	Veronika	Action|Drama	94.00	http://dummyimage.com/100x100.png/ff4444/ffffff
39	fgofford10	Freeland	Drama	64.00	http://dummyimage.com/120x100.png/ff4444/ffffff
40	schaffyn11	Scottie	Drama	18.00	http://dummyimage.com/199x100.png/cc0000/ffffff
41	fgiacobillo12	Flinn	Children|Comedy|Drama	26.00	http://dummyimage.com/138x100.png/dddddd/000000
42	pedsell13	Pepillo	Documentary	51.00	http://dummyimage.com/103x100.png/cc0000/ffffff
43	msaywood14	Malachi	Drama	5.00	http://dummyimage.com/187x100.png/cc0000/ffffff
44	bhiley15	Bellina	Crime|Mystery|Thriller	80.00	http://dummyimage.com/189x100.png/ff4444/ffffff
45	mhuish16	Mariya	Animation|Children|Comedy	35.00	http://dummyimage.com/102x100.png/dddddd/000000
46	qfusco17	Quinn	Comedy|Drama|Romance	49.00	http://dummyimage.com/231x100.png/dddddd/000000
47	pmcgroarty18	Pansy	Adventure|Fantasy|Sci-Fi	50.00	http://dummyimage.com/174x100.png/5fa2dd/ffffff
48	bevett19	Brena	Comedy|Drama|War	63.00	http://dummyimage.com/162x100.png/ff4444/ffffff
49	sscarratt1a	Sheelah	Action|Adventure|Drama|Fantasy|Mystery|IMAX	79.00	http://dummyimage.com/157x100.png/ff4444/ffffff
50	fmccarlie1b	Frankie	Comedy|Romance	38.00	http://dummyimage.com/227x100.png/cc0000/ffffff
51	rplaskitt1c	Reid	Drama	94.00	http://dummyimage.com/248x100.png/ff4444/ffffff
52	ifuente1d	Irina	(no genres listed)	98.00	http://dummyimage.com/106x100.png/5fa2dd/ffffff
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (uid, username, email, password, active, expiry, admin) FROM stdin;
2	apple	apple@mail.com	$2a$10$38raH.z8lb.gYj/c8EVWuORhWdc7js5WOrghcLacsXFC66ML0XPvG	t	\N	f
4	banana	banana@mail.com	$2a$10$oaUqlIEQMcOdn8UFezlNGuDaEQsX3TUBFULhLlGeoMiOJA.mXsNCe	t	\N	f
6	myself	appp@mail.com	$2a$10$y.Gj0mOdHkTSUJQpCOTt6OpPRK3GX1LpNsOb/FUEgSc1.Yc3gW5qq	t	\N	f
7	myselfs	apppd@mail.com	$2a$10$mRGwxMCGU/txIVBhrIJ73.zZrXG.Sperngvnaubm0EXzUVrmVcJLO	t	\N	f
8	myselfsa	appdpd@mail.com	$2a$10$IC.dEOnyZpl.dAZ.XRvYEOp1q3DyZDjdR/ZAnDmwhqEly9AK7FwKS	t	\N	f
9	fresh	new@mail.com	$2a$10$5wT9fSKonGXVxx2UEjeAW.Z5d047DedBwH0.9nkqKXn.uVhF5o6fC	t	\N	f
10	fresh	newdty@mail.com	$2a$10$Iei.qnjzMfrT9S4Rpxarauvy1kZImZXWVwukPm9lKtDiuiJyjQhg2	t	\N	f
1	friedcheese	friedd@mail.com	$2a$10$xMzjkTM4wWW4CRiEyV8oe.NTIiKuvqedLOdA6zmMWUN2G33Z3rFFa	f	\N	f
5	banana	fried@mail.com	$2a$10$A06wv0fqu/oCsS9mUgfuI.3O7HbPW70YVchMQJsNCXJsAQnBRuk8a	t	\N	t
11	fresh	newdtyd@mail.com	$2a$10$3.vAhiCbMH1YRhT34BrMYOXh1bs2xlHVJDiz3Mtx2KV.8t8QGHm8.	t	\N	f
\.


--
-- Data for Name: bought_books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.bought_books (uid, bookid, book_name, genre, download_url, review_id) FROM stdin;
1	1	Sample Book	Fiction	abc.com	\N
1	2	Sample Book2	Fantasy	https://example.com/sample-book.pdf	\N
2	1	Sample Book	Fiction	https://example.com/sample-book.pdf	\N
1	8	fgrabert5	Drama	http://dummyimage.com/161x100.png/5fa2dd/ffffff	\N
1	9	cbladder6	Film-Noir	http://dummyimage.com/243x100.png/cc0000/ffffff	\N
1	10	sarlett7	Comedy|Drama	http://dummyimage.com/145x100.png/cc0000/ffffff	\N
5	2	Sample Book2	Fantasy	https://example.com/sample-book.pdf	\N
5	3	learingey0	(no genres listed)	http://dummyimage.com/129x100.png/dddddd/000000	\N
5	4	tklemt1	Action|Comedy	http://dummyimage.com/140x100.png/cc0000/ffffff	\N
5	6	mcandlin3	Documentary|Musical	http://dummyimage.com/166x100.png/5fa2dd/ffffff	\N
5	8	fgrabert5	Drama	http://dummyimage.com/161x100.png/5fa2dd/ffffff	\N
5	10	sarlett7	Comedy|Drama	http://dummyimage.com/145x100.png/cc0000/ffffff	\N
5	11	scardenoza8	Thriller	http://dummyimage.com/196x100.png/5fa2dd/ffffff	\N
9	11	scardenoza8	Thriller	http://dummyimage.com/196x100.png/5fa2dd/ffffff	\N
9	12	jbalbeck9	Drama|Horror|Mystery	http://dummyimage.com/145x100.png/cc0000/ffffff	\N
9	8	fgrabert5	Drama	http://dummyimage.com/161x100.png/5fa2dd/ffffff	\N
9	5	gbranscombe2	Comedy	http://dummyimage.com/103x100.png/ff4444/ffffff	\N
9	15	alambertsc	Drama	http://dummyimage.com/164x100.png/5fa2dd/ffffff	\N
\.


--
-- Data for Name: cart; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cart (uid, bookid, book_name, author, cost) FROM stdin;
2	1	Sample Book	Jane Author	29.99
1	10	sarlett7	Stacy	96.00
\.


--
-- Data for Name: reviews; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reviews (reviewid, uid, bookid, review) FROM stdin;
3	1	1	this is a review
5	5	1	this is a review2
6	1	2	apple banana
7	9	15	very nice
\.


--
-- Name: books_bookid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.books_bookid_seq', 54, true);


--
-- Name: reviews_reviewid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reviews_reviewid_seq', 7, true);


--
-- Name: users_uid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_uid_seq', 11, true);


--
-- PostgreSQL database dump complete
--

