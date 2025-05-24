CREATE TABLE public.staging_sales (
	order_id int4 NULL,
	product_id varchar(20) NULL,
	customer_id varchar(20) NULL,
	product_name varchar(255) NULL,
	category varchar(100) NULL,
	region varchar(100) NULL,
	date_of_sale date NULL,
	quantity_sold int4 NULL,
	unit_price numeric(10, 2) NULL,
	discount numeric(5, 2) NULL,
	shipping_cost numeric(10, 2) NULL,
	payment_method varchar(50) NULL,
	customer_name varchar(255) NULL,
	customer_email varchar(255) NULL,
	customer_address text NULL
);

-- Create categories table
CREATE TABLE public.categories (
	category_id serial4 NOT NULL,
	"name" varchar(100) NOT NULL,
	CONSTRAINT categories_name_key UNIQUE (name),
	CONSTRAINT categories_pkey PRIMARY KEY (category_id)
);

-- Create regions table
CREATE TABLE public.regions (
	region_id serial4 NOT NULL,
	"name" varchar(100) NOT NULL,
	CONSTRAINT regions_name_key UNIQUE (name),
	CONSTRAINT regions_pkey PRIMARY KEY (region_id)
);

-- Create products table
CREATE TABLE public.products (
	product_id varchar(20) NOT NULL,
	"name" varchar(255) NOT NULL,
	category_id int4 NOT NULL,
	unit_price numeric(10, 2) NOT NULL,
	CONSTRAINT products_pkey PRIMARY KEY (product_id),
	CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(category_id)
);

-- Create customers table
CREATE TABLE public.customers (
	customer_id varchar(20) NOT NULL,
	"name" varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	address text NOT NULL,
	CONSTRAINT customers_email_key UNIQUE (email),
	CONSTRAINT customers_pkey PRIMARY KEY (customer_id)
);

-- Create orders table
CREATE TABLE public.orders (
	order_id int4 NOT NULL,
	product_id varchar(20) NOT NULL,
	customer_id varchar(20) NOT NULL,
	region_id int4 NOT NULL,
	date_of_sale date NOT NULL,
	quantity_sold int4 NOT NULL,
	discount numeric(5, 2) DEFAULT 0.00 NULL,
	shipping_cost numeric(10, 2) DEFAULT 0.00 NULL,
	payment_method varchar(50) NOT NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (order_id),
	CONSTRAINT orders_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(customer_id),
	CONSTRAINT orders_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(product_id),
	CONSTRAINT orders_region_id_fkey FOREIGN KEY (region_id) REFERENCES public.regions(region_id)
);

-- public.api_log_capture;
CREATE TABLE public.api_log_capture (
	request_id uuid NOT NULL,
	respbody text NULL,
	response_status int8 DEFAULT 0 NULL,
	reqdatetime timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	realip varchar(45) DEFAULT NULL::character varying NULL,
	forwardedip varchar(45) DEFAULT NULL::character varying NULL,
	"method" varchar(10) DEFAULT NULL::character varying NULL,
	"path" text NULL,
	host text NULL,
	remoteaddr varchar(45) DEFAULT NULL::character varying NULL,
	"header" text NULL,
	endpoint text NULL,
	respdatetime timestamptz NULL,
	reqbody text NULL,
	requesttime int8 DEFAULT 0 NULL,
	responsetime int8 DEFAULT 0 NULL,
	CONSTRAINT api_log_capture_pkey PRIMARY KEY (request_id)
);
CREATE INDEX idx_endpoint ON public.api_log_capture USING btree ("left"(endpoint, 100));
CREATE INDEX idx_forwardedip ON public.api_log_capture USING btree (forwardedip);
CREATE INDEX idx_host ON public.api_log_capture USING btree ("left"(host, 100));
CREATE INDEX idx_reqdatetime ON public.api_log_capture USING btree (reqdatetime);

CREATE TABLE public.refresh_log_activity (
	request_id uuid NOT NULL,
	status varchar(20) NOT NULL,
	total_records_affected int8 NULL,
	error_message text NULL,
	created_time timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_time timestamptz NULL,
	created_by varchar(200) NULL,
	updated_by varchar(200) NULL,
	refresh_type varchar(50) NULL,
	duration_seconds int4 NULL,
	CONSTRAINT refresh_log_activity_pkey PRIMARY KEY (request_id)
);