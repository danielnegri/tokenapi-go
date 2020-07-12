create table if not exists secret_tokens
(
    data text not null
        constraint secret_tokens_data_key
            unique
        constraint no_dash_character
            check (strpos(data, '-'::text) = 0)
);

-- alter table secret_tokens owner to interviewed;
