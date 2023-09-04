import { createClient } from "@supabase/supabase-js";
import { getMostRecentSignalTrigger } from "./vaidya-service";

//----------------------------------------------
// supabase client
//----------------------------------------------

const supabaseUrl = import.meta.env.PUBLIC_SUPABASE_URL;
const supabaseAnonKey = import.meta.env.PUBLIC_SUPABASE_ANON_KEY;

const supabase = createClient<Database>(supabaseUrl, supabaseAnonKey, {
  auth: {
    persistSession: false,
  },
});

//----------------------------------------------
// supabase generated types
// command: npx supabase gen types typescript --project-id "<PROJECT_ID>" --schema public
//----------------------------------------------

type Json = string | number | boolean | null | { [key: string]: Json | undefined } | Json[];

interface Database {
  public: {
    Tables: {
      emails: {
        Row: {
          email: string;
        };
        Insert: {
          email: string;
        };
        Update: {
          email?: string;
        };
        Relationships: [];
      };
      signals: {
        Row: {
          id: number;
          low_1_date: string;
          low_2_date: string;
          ticker: string;
          trigger_date: string;
        };
        Insert: {
          id?: number;
          low_1_date: string;
          low_2_date: string;
          ticker: string;
          trigger_date: string;
        };
        Update: {
          id?: number;
          low_1_date?: string;
          low_2_date?: string;
          ticker?: string;
          trigger_date?: string;
        };
        Relationships: [];
      };
      watchlist: {
        Row: {
          last_trigger: string | null;
          ticker: string;
        };
        Insert: {
          last_trigger?: string | null;
          ticker: string;
        };
        Update: {
          last_trigger?: string | null;
          ticker?: string;
        };
        Relationships: [];
      };
    };
    Views: {
      [_ in never]: never;
    };
    Functions: {
      [_ in never]: never;
    };
    Enums: {
      [_ in never]: never;
    };
    CompositeTypes: {
      [_ in never]: never;
    };
  };
}

//-----------------------------------------------
// watchlist helpers
//-----------------------------------------------

export const getWatchlist = async () => {
  const { data, error } = await supabase.from("watchlist").select("*");
  if (error) {
    throw error;
  }
  return data;
};

export const addToWatchlist = async (tickers: string) => {
  const tickersArray = tickers.split(" ");

  const insertions: { ticker: string; last_trigger?: string }[] = [];
  for (const ticker of tickersArray) {
    const lastSignalTrigger = await getMostRecentSignalTrigger(ticker);
    insertions.push({ ticker: ticker, last_trigger: lastSignalTrigger });
  }

  const { data, error } = await supabase
    .from("watchlist")
    .upsert(insertions.map((ticker) => ({ ticker: ticker.ticker, last_trigger: ticker.last_trigger })));
  if (error) {
    throw error;
  }
  return data;
};

export const removeFromWatchlist = async (ticker: string) => {
  const { data, error } = await supabase.from("watchlist").delete().match({ ticker: ticker });
  if (error) {
    throw error;
  }
  return data;
};

export const existsInWatchlist = async (ticker: string) => {
  const list = await getWatchlist();
  return list.some((item) => item.ticker === ticker);
};

//-----------------------------------------------
// emails helpers
//-----------------------------------------------

export const addEmail = async (email: string) => {
  const { data, error } = await supabase.from("emails").insert([{ email: email }]);
  if (error) {
    throw error;
  }
  return data;
};

export const removeEmail = async (email: string) => {
  const { data, error } = await supabase.from("emails").delete().match({ email: email });
  if (error) {
    throw error;
  }
  return data;
};

//-----------------------------------------------
// signals helpers
//-----------------------------------------------

export const getSignals = async (ticker: string) => {
  const { data, error } = await supabase.from("signals").select("*").match({ ticker: ticker });
  if (error) {
    throw error;
  }
  return data;
};
