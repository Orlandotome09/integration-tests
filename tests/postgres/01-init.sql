CREATE TABLE economical_activities(
    code_id text primary key not null,
    description text not null,
    risk_value bool,
    created_at timestamptz,
    updated_at timestamptz
);

insert into economical_activities values ('9900-8/00','Organismos internacionais e outras instituições extraterritoriais', true);