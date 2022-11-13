DROP DATABASE IF EXISTS "AvitoTest";
CREATE DATABASE "AvitoTest" WITH ENCODING = 'UTF8';

\c "AvitoTest";

CREATE TABLE customers
(
  id serial, 
  balance money NOT NULL
);

CREATE TABLE services
(
  id serial, 
  cost money NOT NULL
);


CREATE TABLE orders
(
  id serial,
  customers_id smallint NOT NULL,
  services_id smallint NOT NULL, 
  price money NOT NULL
);
ALTER TABLE customers
ADD CONSTRAINT customers_pkey PRIMARY KEY (id),
ADD CONSTRAINT customers_balance CHECK (balance >= 0);

ALTER TABLE services
ADD CONSTRAINT services_pkey PRIMARY KEY (id),
ADD CONSTRAINT services_cost CHECK (cost >= 0);

ALTER TABLE orders
ADD CONSTRAINT orders_pkey PRIMARY KEY (id),
ADD CONSTRAINT orders_customers_id FOREIGN KEY (customers_id) REFERENCES customers(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
ADD CONSTRAINT orders_services_id FOREIGN KEY (services_id) REFERENCES services(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
ADD CONSTRAINT orders_price CHECK (price >= 0)
