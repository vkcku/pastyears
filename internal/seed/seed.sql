insert into
  questions.exams (name, short_name)
values
  ('Civil Service Exam', 'CSE'),
  ('Combined Defense Service', 'CDS');

insert into
  questions.papers (
    name,
    short_name,
    is_optional,
    paper_type,
    exam_id
  )
values
  (
    'General Studies I',
    'GS I',
    false,
    'prelims',
    (
      select
        id
      from
        questions.exams
      where
        short_name = 'CSE'
    )
  ),
  (
    'General Knowledge',
    '',
    false,
    'prelims',
    (
      select
        id
      from
        questions.exams
      where
        short_name = 'CDS'
    )
  );

with
  editions as (
    select
      1 as edition
    union all
    select
      2 as edition
  )
insert into
  questions.question_papers (paper_id, year, exam_edition)
select
  (
    select
      id
    from
      questions.papers
    where
      name = 'General Studies I'
  ),
  year,
  0
from
  generate_series(2014, 2025) as year
union all
select
  (
    select
      id
    from
      questions.papers
    where
      name = 'General Knowledge'
  ),
  year,
  edition
from
  generate_series(2014, 2025) as year
  cross join editions;

insert into
  questions.subjects (name)
values
  ('Polity'),
  ('History'),
  ('Economics'),
  ('Geography'),
  ('Environment'),
  ('Science & Technology');

insert into
  questions.topics (name)
values
  ('Biodiversity'),
  ('Constitution'),
  ('Inclusive growth');
