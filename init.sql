CREATE SCHEMA IF NOT EXISTS url_shortener;

-- URL table sequence
CREATE SEQUENCE IF NOT EXISTS url_shortener.urls_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- URL Analytics table sequence
CREATE SEQUENCE IF NOT EXISTS url_shortener.url_analytics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS url_shortener.url
(
    id BIGINT NOT NULL DEFAULT nextval('url_shortener.urls_id_seq'::regclass),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    short_code VARCHAR(255) NOT NULL,
    long_url TEXT NOT NULL,
    last_accessed_at TIMESTAMPTZ,
    click_count INTEGER DEFAULT 0,
    is_custom_url BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    CONSTRAINT urls_pkey PRIMARY KEY (id),
    CONSTRAINT urls_long_url_key UNIQUE (long_url),
    CONSTRAINT urls_short_code_key UNIQUE (short_code)
);

ALTER TABLE IF EXISTS url_shortener.url
    OWNER TO postgres;

CREATE INDEX IF NOT EXISTS idx_url_short_code
    ON url_shortener.url (short_code);

CREATE INDEX IF NOT EXISTS idx_url_is_active
    ON url_shortener.url (is_active);

CREATE TABLE IF NOT EXISTS url_shortener.url_analytics
(
    id BIGINT NOT NULL DEFAULT nextval('url_shortener.url_analytics_id_seq'::regclass),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    date VARCHAR(255) NOT NULL,
    total_clicks INTEGER DEFAULT 0,
    CONSTRAINT url_analytics_pkey PRIMARY KEY (id),
    CONSTRAINT url_analytics_date_key UNIQUE (date)
);

ALTER TABLE IF EXISTS url_shortener.url_analytics
    OWNER TO postgres;

CREATE INDEX IF NOT EXISTS idx_url_analytics_date
    ON url_shortener.url_analytics (date);
