-- migrate:up
create schema if not exists questions;

create table questions.exams (
  id int generated always as identity primary key,
  name text not null,
  short_name text not null,
  created_at timestamptz not null default (now() at time zone 'UTC'),
  updated_at timestamptz not null default (now() at time zone 'UTC'),
  constraint uq__exams__name unique (name),
  constraint ck__exams__name__not_empty check (name != '')
);

create type questions.paper_type as enum('prelims', 'mains');

create table questions.papers (
  id int generated always as identity primary key,
  name text not null,
  short_name text not null,
  is_optional boolean not null,
  paper_type questions.paper_type not null,
  exam_id int not null references questions.exams (id),
  created_at timestamptz not null default (now() at time zone 'UTC'),
  updated_at timestamptz not null default (now() at time zone 'UTC'),
  constraint uq__papers__name__paper_type__exam_id__is_optional unique (name, paper_type, exam_id, is_optional),
  constraint ck__papers__paper_type__enum check (paper_type in ('prelims', 'mains')),
  constraint ck__papers__name__not_empty check (name != ''),
  constraint ck__papers__is_optional__only_for_mains check (
    is_optional = false
    or (
      is_optional = true
      and paper_type = 'mains'
    )
  )
);

create table questions.question_papers (
  id int generated always as identity primary key,
  paper_id int not null references questions.papers (id),
  year smallint not null,
  exam_edition smallint not null,
  constraint ck__question_papers__exam_edition__valid check (exam_edition in (0, 1, 2)),
  constraint ck__question_papers__year__in_range check (
    year >= 1990
    and year <= 2025
  ),
  constraint uq__question_papers__year__paper_id__exam_edition unique (year, paper_id, exam_edition)
);

create table questions.subjects (
  id int generated always as identity primary key,
  name text not null,
  constraint uq__subjects__name unique (name),
  constraint ck__subjects__name__not_empty check (name != '')
);

create table questions.topics (
  id int generated always as identity primary key,
  name text not null,
  constraint uq__topics__name unique (name),
  constraint ck__topics__name__not_empty check (name != '')
);

create table questions.rich_text (
  id int generated always as identity primary key,
  raw_text text not null,
  html text not null,
  constraint ck__rich_text__not_empty check (
    raw_text != ''
    and html != ''
  )
);

create table questions.prelims_questions (
  id_pk int generated always as identity primary key,
  id text not null,
  rt_id int not null references questions.rich_text (id),
  answer_a int not null references questions.rich_text (id),
  answer_b int not null references questions.rich_text (id),
  answer_c int not null references questions.rich_text (id),
  answer_d int not null references questions.rich_text (id),
  correct_answer text not null,
  qp_id int not null references questions.question_papers (id),
  subject_id int not null references questions.subjects (id),
  constraint uq__prelims_questions__id unique (id),
  constraint ck__prelims_questions__id__not_empty check (id != ''),
  constraint ck__prelims_questions__correct_answer__valid check (correct_answer in ('a', 'b', 'c', 'd'))
);

create table questions.prelims_questions_topics (
  id int generated always as identity primary key,
  question_id int not null references questions.prelims_questions (id_pk),
  topic_id int not null references questions.topics (id),
  constraint uq__prelims_questions_topics__question_id__topic_id unique (question_id, topic_id)
);

-- migrate:down
drop table if exists questions.prelims_questions_topics;

drop table if exists questions.prelims_questions;

drop table if exists questions.rich_text;

drop table if exists questions.topics;

drop table if exists questions.subjects;

drop table if exists questions.question_papers;

drop table if exists questions.papers;

drop type if exists questions.paper_type;

drop table if exists questions.exams;
