create table
  public.watchlist (
    ticker text not null,
    last_trigger date null,
    constraint watchlist_pkey primary key (ticker)
  ) tablespace pg_default;

create table
public.emails (
  email text not null,
  constraint emails_pkey primary key (email)
) tablespace pg_default;

-- -- the following DDLs are deprecated for this project

-- create
-- or replace function create_ticker_tables () returns trigger as $$
-- begin
--   execute 'CREATE TABLE IF NOT EXISTS ' || 'days_' || replace(NEW.ticker, '.', '_') || ' (
--     date date PRIMARY KEY NOT NULL,
--     open real NULL,
--     high real NULL,
--     low real NULL,
--     close real NULL,
--     volume integer NULL,
--     macd real NULL,
--     rsi real NULL,
--     ema_12 real NULL,
--     ema_26 real NULL,
--     avg_gain real NULL,
--     avg_loss real NULL
--   ) ';
--   execute 'CREATE TABLE IF NOT EXISTS ' || 'vaidya_' || replace(NEW.ticker, '.', '_') || ' (
--     trigger_date date PRIMARY KEY NOT NULL,
--     low_2_date date NOT NULL,
--     low_1_date date NOT NULL
--   ) ';
--   return NEW;
-- end;
-- $$ language plpgsql;

-- create table
-- public.tickers (
--   ticker text not null,
--   first_date date not null,
--   last_date date not null,
--   last_vaiday_signal_date date null,
--   constraint tickers_pkey primary key (ticker)
-- ) tablespace pg_default;

-- create trigger insert_ticker
-- after insert on tickers for each row
-- execute function create_ticker_tables ();