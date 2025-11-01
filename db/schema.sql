CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE exams (
  id integer primary key,
  name text not null,
  short_name text not null,
  created_at text not null default (datetime('now')),
  updated_at text not null default (datetime('now')),
  constraint uq__exams__name unique (name),
  constraint ck__exams__name__not_empty check (name != '')
) strict;
CREATE TABLE papers (
  id integer primary key,
  name text not null,
  short_name text not null,
  is_optional integer not null,
  paper_type text not null,
  exam_id integer not null references exams (id),
  created_at text not null default (datetime('now')),
  updated_at text not null default (datetime('now')),
  constraint uq__papers__name__paper_type__exam_id__is_optional unique (name, paper_type, exam_id, is_optional),
  constraint ck__papers__is_optional__boolean check (is_optional in (0, 1)),
  constraint ck__papers__paper_type__enum check (paper_type in ('prelims', 'mains')),
  constraint ck__papers__name__not_empty check (name != ''),
  constraint ck__papers__is_optional__only_for_mains check (
    is_optional = 0
    or (
      is_optional = 1
      and paper_type = 'mains'
    )
  )
) strict;
CREATE TABLE question_papers (
  id integer primary key,
  paper_id integer not null references papers (id),
  year integer not null,
  exam_edition integer not null,
  constraint ck__question_papers__exam_edition__valid check (exam_edition in (0, 1, 2)),
  constraint ck__question_papers__year__in_range check (
    year >= 1990
    and year <= 2025
  ),
  constraint uq__question_papers__year__paper_id__exam_edition unique (year, paper_id, exam_edition)
) strict;
CREATE TABLE subjects (
  id integer primary key,
  name text not null,
  constraint uq__subjects__name unique (name),
  constraint ck__subjects__name__not_empty check (name != '')
) strict;
CREATE TABLE topics (
  id integer primary key,
  name text not null,
  constraint uq__topics__name unique (name),
  constraint ck__topics__name__not_empty check (name != '')
) strict;
CREATE TABLE rich_text (
  id integer primary key,
  raw_text text not null,
  html text not null,
  constraint ck__rich_text__not_empty check (
    raw_text != ''
    and html != ''
  )
) strict;
CREATE TABLE prelims_questions (
  id_pk integer primary key,
  id text not null,
  rt_id integer not null references rich_text (id),
  answer_a integer not null references rich_text (id),
  answer_b integer not null references rich_text (id),
  answer_c integer not null references rich_text (id),
  answer_d integer not null references rich_text (id),
  correct_answer text not null,
  qp_id integer not null references question_papers (id),
  subject_id integer not null references subjects (id),
  constraint uq__prelims_questions__id unique (id),
  constraint ck__prelims_questions__id__not_empty check (id != ''),
  constraint ck__prelims_questions__correct_answer__valid check (correct_answer in ('a', 'b', 'c', 'd'))
) strict;
CREATE TABLE prelims_questions_topics (
  id integer primary key,
  question_id integer not null references prelims_questions (id_pk),
  topic_id integer not null references topics (id),
  constraint uq__prelims_questions_topics__question_id__topic_id unique (question_id, topic_id)
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20251029001051');
