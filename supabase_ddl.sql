create
or replace function create_days_ticker_table () returns trigger as $$
begin
  execute 'CREATE TABLE IF NOT EXISTS ' || 'days_' || replace(NEW.ticker, '.', '_') || ' (
    date date PRIMARY KEY NOT NULL,
    open real NULL,
    high real NULL,
    low real NULL,
    close real NULL,
    volume integer NULL,
    macd real NULL,
    rsi real NULL
  ) ';
  execute 'CREATE TABLE IF NOT EXISTS ' || 'vaidya_' || replace(NEW.ticker, '.', '_') || ' (
    trigger_date date PRIMARY KEY NOT NULL,
    low_2_date date NOT NULL,
    low_1_date date NOT NULL
  ) ';
  return NEW;
end;
$$ language plpgsql;

create table
public.tickers (
  ticker text not null,
  first_date date not null,
  last_date date not null,
  constraint tickers_pkey primary key (ticker)
) tablespace pg_default;

create trigger insert_ticker
after insert on tickers for each row
execute function create_days_ticker_table ();