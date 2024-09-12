-- public.stock_keeping_unit definition

-- Drop table

-- DROP TABLE public.stock_keeping_unit;

CREATE TABLE public.stock_keeping_unit (
	snowflake_id varchar DEFAULT ''::character varying NOT NULL,
	"name" varchar DEFAULT ''::character varying NOT NULL,
	code varchar DEFAULT ''::character varying NOT NULL,
	stock_quantity int4 DEFAULT 0 NOT NULL,
	virtual_sales int4 DEFAULT 0 NOT NULL,
	price numeric DEFAULT 0 NOT NULL,
	status int4 DEFAULT 0 NOT NULL,
	sorting int4 DEFAULT 0 NOT NULL,
	object_name varchar DEFAULT ''::character varying NOT NULL,
	bucket_name varchar DEFAULT ''::character varying NOT NULL,
	actual_sales int4 DEFAULT 0 NOT NULL,
	item_id varchar DEFAULT ''::character varying NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	deleted_at timestamp NULL,
	ext varchar DEFAULT ''::character varying NOT NULL
);