insert into
  exams (name, short_name)
values
  ('Civil Service Exam', 'CSE'),
  ('Combined Defense Service', 'CDS');

insert into
  papers (
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
        exams
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
        exams
      where
        short_name = 'CDS'
    )
  );

with recursive
  generate_years (year) as (
    select
      2014 as year
    union all
    select
      year + 1
    from
      generate_years
    where
      year + 1 <= 2025
  ),
  generate_editions as (
    select
      1 as edition
    union all
    select
      2 as edition
  )
insert into
  question_papers (paper_id, year, exam_edition)
select
  (
    select
      id
    from
      papers
    where
      name = 'General Studies I'
  ),
  year,
  0
from
  generate_years
union all
select
  (
    select
      id
    from
      papers
    where
      name = 'General Knowledge'
  ),
  year,
  edition
from
  generate_years
  cross join generate_editions;

insert into
  subjects (name)
values
  ('Polity'),
  ('History'),
  ('Economics'),
  ('Geography'),
  ('Environment'),
  ('Science & Technology');

insert into
  topics (name)
values
  ('Biodiversity'),
  ('Constitution'),
  ('Inclusive growth');
