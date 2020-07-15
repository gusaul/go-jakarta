CREATE DATABASE product;
CREATE USER gojakarta WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE "product" to gojakarta;
\connect product;

CREATE TABLE IF NOT EXISTS public.gj_product (
    product_id BIGSERIAL PRIMARY KEY,
    shop_id BIGINT NOT NULL,
    product_name VARCHAR(50),
    product_desc TEXT,
    product_price decimal
);
GRANT ALL ON public.gj_product TO gojakarta;

INSERT INTO public.gj_product (product_id, shop_id, product_name, product_desc, product_price)
VALUES  (123, 321, 'Product A', 'Desc of Product A', 10000),
        (234, 432, 'Product B', 'Desc of Product B', 20000),
        (345, 543, 'Product C', 'Desc of Product C', 30000),
        (456, 654, 'Product D', 'Desc of Product D', 40000),
        (567, 765, 'Product E', 'Desc of Product E', 50000);

CREATE TABLE IF NOT EXISTS public.gj_stats (
    product_id BIGINT PRIMARY KEY,
    view INTEGER,
    transactions INTEGER,
    review INTEGER,
    talk INTEGER
);
GRANT ALL ON public.gj_stats TO gojakarta;

INSERT INTO public.gj_stats (product_id, view, transactions, review, talk)
VALUES  (123, 123000, 12300, 1230, 12300),
        (234, 234000, 23400, 2340, 23400),
        (345, 345000, 34500, 3450, 34500),
        (456, 456000, 45600, 4560, 45600),
        (567, 567000, 56700, 5670, 56700);

CREATE TABLE IF NOT EXISTS public.gj_picture (
    picture_id BIGSERIAL PRIMARY KEY,
    product_id BIGINT,
    file_path VARCHAR(50),
    file_name VARCHAR(50)
);
CREATE INDEX idx_gp_product_id ON public.gj_picture (product_id);
GRANT ALL ON public.gj_picture TO gojakarta;

INSERT INTO public.gj_picture (picture_id, product_id, file_path, file_name)
VALUES  (1321, 123, 'image_product', 'product_123A.jpg'),
        (1322, 123, 'image_product', 'product_123B.jpg'),
        (1323, 123, 'image_product', 'product_123C.jpg'),
        (1432, 234, 'image_product', 'product_234A.jpg'),
        (1433, 234, 'image_product', 'product_234B.jpg'),
        (1434, 234, 'image_product', 'product_234C.jpg'),
        (1543, 345, 'image_product', 'product_345A.jpg'),
        (1544, 345, 'image_product', 'product_345B.jpg'),
        (1654, 456, 'image_product', 'product_456A.jpg'),
        (1655, 456, 'image_product', 'product_456B.jpg'),
        (1656, 456, 'image_product', 'product_456C.jpg'),
        (1657, 456, 'image_product', 'product_456D.jpg'),
        (1765, 567, 'image_product', 'product_567A.jpg'),
        (1766, 567, 'image_product', 'product_567B.jpg');
