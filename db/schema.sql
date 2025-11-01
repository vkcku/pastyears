
-- Dumped from database version 18.0
-- Dumped by pg_dump version 18.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: questions; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA questions;


--
-- Name: paper_type; Type: TYPE; Schema: questions; Owner: -
--

CREATE TYPE questions.paper_type AS ENUM (
    'prelims',
    'mains'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: exams; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.exams (
    id integer NOT NULL,
    name text NOT NULL,
    short_name text NOT NULL,
    created_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'::text) NOT NULL,
    updated_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'::text) NOT NULL,
    CONSTRAINT ck__exams__name__not_empty CHECK ((name <> ''::text))
);


--
-- Name: exams_id_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

ALTER TABLE questions.exams ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME questions.exams_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: papers; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.papers (
    id integer NOT NULL,
    name text NOT NULL,
    short_name text NOT NULL,
    is_optional boolean NOT NULL,
    paper_type questions.paper_type NOT NULL,
    exam_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'::text) NOT NULL,
    updated_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'::text) NOT NULL,
    CONSTRAINT ck__papers__is_optional__only_for_mains CHECK (((is_optional = false) OR ((is_optional = true) AND (paper_type = 'mains'::questions.paper_type)))),
    CONSTRAINT ck__papers__name__not_empty CHECK ((name <> ''::text)),
    CONSTRAINT ck__papers__paper_type__enum CHECK ((paper_type = ANY (ARRAY['prelims'::questions.paper_type, 'mains'::questions.paper_type])))
);


--
-- Name: papers_id_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

ALTER TABLE questions.papers ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME questions.papers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: prelims_questions; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.prelims_questions (
    id_pk integer NOT NULL,
    id text NOT NULL,
    rt_id integer NOT NULL,
    answer_a integer NOT NULL,
    answer_b integer NOT NULL,
    answer_c integer NOT NULL,
    answer_d integer NOT NULL,
    correct_answer text NOT NULL,
    qp_id integer NOT NULL,
    subject_id integer NOT NULL,
    CONSTRAINT ck__prelims_questions__correct_answer__valid CHECK ((correct_answer = ANY (ARRAY['a'::text, 'b'::text, 'c'::text, 'd'::text]))),
    CONSTRAINT ck__prelims_questions__id__not_empty CHECK ((id <> ''::text))
);


--
-- Name: prelims_questions_id_pk_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

ALTER TABLE questions.prelims_questions ALTER COLUMN id_pk ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME questions.prelims_questions_id_pk_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: prelims_questions_topics; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.prelims_questions_topics (
    id integer NOT NULL,
    question_id integer NOT NULL,
    topic_id integer NOT NULL
);


--
-- Name: prelims_questions_topics_id_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

ALTER TABLE questions.prelims_questions_topics ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME questions.prelims_questions_topics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: question_papers; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.question_papers (
    id integer NOT NULL,
    paper_id integer NOT NULL,
    year smallint NOT NULL,
    exam_edition smallint NOT NULL,
    CONSTRAINT ck__question_papers__exam_edition__valid CHECK ((exam_edition = ANY (ARRAY[0, 1, 2]))),
    CONSTRAINT ck__question_papers__year__in_range CHECK (((year >= 1990) AND (year <= 2025)))
);


--
-- Name: question_papers_id_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

ALTER TABLE questions.question_papers ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME questions.question_papers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: rich_text; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.rich_text (
    id integer NOT NULL,
    raw_text text NOT NULL,
    html text NOT NULL,
    CONSTRAINT ck__rich_text__not_empty CHECK (((raw_text <> ''::text) AND (html <> ''::text)))
);


--
-- Name: rich_text_id_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

ALTER TABLE questions.rich_text ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME questions.rich_text_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: subjects; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.subjects (
    id integer NOT NULL,
    name text NOT NULL,
    CONSTRAINT ck__subjects__name__not_empty CHECK ((name <> ''::text))
);


--
-- Name: subjects_id_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

ALTER TABLE questions.subjects ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME questions.subjects_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: topics; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.topics (
    id integer NOT NULL,
    name text NOT NULL,
    CONSTRAINT ck__topics__name__not_empty CHECK ((name <> ''::text))
);


--
-- Name: topics_id_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

ALTER TABLE questions.topics ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME questions.topics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: exams exams_pkey; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.exams
    ADD CONSTRAINT exams_pkey PRIMARY KEY (id);


--
-- Name: papers papers_pkey; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.papers
    ADD CONSTRAINT papers_pkey PRIMARY KEY (id);


--
-- Name: prelims_questions prelims_questions_pkey; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT prelims_questions_pkey PRIMARY KEY (id_pk);


--
-- Name: prelims_questions_topics prelims_questions_topics_pkey; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions_topics
    ADD CONSTRAINT prelims_questions_topics_pkey PRIMARY KEY (id);


--
-- Name: question_papers question_papers_pkey; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.question_papers
    ADD CONSTRAINT question_papers_pkey PRIMARY KEY (id);


--
-- Name: rich_text rich_text_pkey; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.rich_text
    ADD CONSTRAINT rich_text_pkey PRIMARY KEY (id);


--
-- Name: subjects subjects_pkey; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.subjects
    ADD CONSTRAINT subjects_pkey PRIMARY KEY (id);


--
-- Name: topics topics_pkey; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.topics
    ADD CONSTRAINT topics_pkey PRIMARY KEY (id);


--
-- Name: exams uq__exams__name; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.exams
    ADD CONSTRAINT uq__exams__name UNIQUE (name);


--
-- Name: papers uq__papers__name__paper_type__exam_id__is_optional; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.papers
    ADD CONSTRAINT uq__papers__name__paper_type__exam_id__is_optional UNIQUE (name, paper_type, exam_id, is_optional);


--
-- Name: prelims_questions uq__prelims_questions__id; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT uq__prelims_questions__id UNIQUE (id);


--
-- Name: prelims_questions_topics uq__prelims_questions_topics__question_id__topic_id; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions_topics
    ADD CONSTRAINT uq__prelims_questions_topics__question_id__topic_id UNIQUE (question_id, topic_id);


--
-- Name: question_papers uq__question_papers__year__paper_id__exam_edition; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.question_papers
    ADD CONSTRAINT uq__question_papers__year__paper_id__exam_edition UNIQUE (year, paper_id, exam_edition);


--
-- Name: subjects uq__subjects__name; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.subjects
    ADD CONSTRAINT uq__subjects__name UNIQUE (name);


--
-- Name: topics uq__topics__name; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.topics
    ADD CONSTRAINT uq__topics__name UNIQUE (name);


--
-- Name: papers papers_exam_id_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.papers
    ADD CONSTRAINT papers_exam_id_fkey FOREIGN KEY (exam_id) REFERENCES questions.exams(id);


--
-- Name: prelims_questions prelims_questions_answer_a_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT prelims_questions_answer_a_fkey FOREIGN KEY (answer_a) REFERENCES questions.rich_text(id);


--
-- Name: prelims_questions prelims_questions_answer_b_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT prelims_questions_answer_b_fkey FOREIGN KEY (answer_b) REFERENCES questions.rich_text(id);


--
-- Name: prelims_questions prelims_questions_answer_c_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT prelims_questions_answer_c_fkey FOREIGN KEY (answer_c) REFERENCES questions.rich_text(id);


--
-- Name: prelims_questions prelims_questions_answer_d_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT prelims_questions_answer_d_fkey FOREIGN KEY (answer_d) REFERENCES questions.rich_text(id);


--
-- Name: prelims_questions prelims_questions_qp_id_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT prelims_questions_qp_id_fkey FOREIGN KEY (qp_id) REFERENCES questions.question_papers(id);


--
-- Name: prelims_questions prelims_questions_rt_id_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT prelims_questions_rt_id_fkey FOREIGN KEY (rt_id) REFERENCES questions.rich_text(id);


--
-- Name: prelims_questions prelims_questions_subject_id_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions
    ADD CONSTRAINT prelims_questions_subject_id_fkey FOREIGN KEY (subject_id) REFERENCES questions.subjects(id);


--
-- Name: prelims_questions_topics prelims_questions_topics_question_id_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions_topics
    ADD CONSTRAINT prelims_questions_topics_question_id_fkey FOREIGN KEY (question_id) REFERENCES questions.prelims_questions(id_pk);


--
-- Name: prelims_questions_topics prelims_questions_topics_topic_id_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.prelims_questions_topics
    ADD CONSTRAINT prelims_questions_topics_topic_id_fkey FOREIGN KEY (topic_id) REFERENCES questions.topics(id);


--
-- Name: question_papers question_papers_paper_id_fkey; Type: FK CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.question_papers
    ADD CONSTRAINT question_papers_paper_id_fkey FOREIGN KEY (paper_id) REFERENCES questions.papers(id);


--
-- PostgreSQL database dump complete
--



--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20251029001051');
